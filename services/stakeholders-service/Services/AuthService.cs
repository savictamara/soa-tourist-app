using Microsoft.AspNetCore.Identity;
using Microsoft.Extensions.Options;
using Microsoft.EntityFrameworkCore;
using Microsoft.IdentityModel.Tokens;
using Npgsql;
using StakeholdersService.DTOs;
using StakeholdersService.Entities;
using StakeholdersService.Exceptions;
using StakeholdersService.Models;
using StakeholdersService.Repositories;
using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using System.Text;
using System.Text.Json;

namespace StakeholdersService.Services;

public class AuthService : IAuthService
{
    private static readonly string[] AllowedRoles = ["Guide", "Tourist"];

    private readonly IUserRepository _userRepository;
    private readonly IPasswordHasher<User> _passwordHasher;
    private readonly JwtSettings _jwtSettings;
    private readonly IHttpClientFactory _httpClientFactory;
    private readonly ILogger<AuthService> _logger;

    public AuthService(
        IUserRepository userRepository,
        IOptions<JwtSettings> jwtOptions,
        IHttpClientFactory httpClientFactory,
        ILogger<AuthService> logger)
    {
        _userRepository = userRepository;
        _passwordHasher = new PasswordHasher<User>();
        _jwtSettings = jwtOptions.Value;
        _httpClientFactory = httpClientFactory;
        _logger = logger;
    }

    public async Task<AuthResponseDto> RegisterAsync(RegisterUserRequestDto request, CancellationToken cancellationToken = default)
    {
        var username = request.Username.Trim();
        var email = request.Email.Trim();
        var role = request.Role.Trim();

        if (!AllowedRoles.Contains(role))
        {
            throw new ArgumentException("Role must be Guide or Tourist.");
        }

        if (await _userRepository.UsernameExistsAsync(username, cancellationToken))
        {
            throw new ConflictException("Username already exists.");
        }

        if (await _userRepository.EmailExistsAsync(email, cancellationToken))
        {
            throw new ConflictException("Email already exists.");
        }

        var user = new User
        {
            Username = username,
            Email = email,
            Role = role,
            IsBlocked = false,
            FirstName = string.Empty,
            LastName = string.Empty,
            ProfileImage = null,
            Biography = null,
            Motto = null
        };

        user.PasswordHash = _passwordHasher.HashPassword(user, request.Password);

        User createdUser;

        try
        {
            createdUser = await _userRepository.AddAsync(user, cancellationToken);
        }
        catch (DbUpdateException exception) when (exception.InnerException is PostgresException postgresException &&
                                                 postgresException.SqlState == PostgresErrorCodes.UniqueViolation)
        {
            throw new ConflictException("Username or email already exists.");
        }

        // SAGA T2: sync user node to follower service
        // Compensation C1: delete user from PostgreSQL if follower sync fails
        try
        {
            await SyncUserToFollowerServiceAsync(createdUser, cancellationToken);
            _logger.LogInformation("Registration SAGA completed userId={UserId} username={Username}", createdUser.Id, createdUser.Username);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Registration SAGA T2 failed userId={UserId} username={Username} — executing compensation C1: deleting user", createdUser.Id, createdUser.Username);
            await _userRepository.DeleteAsync(createdUser.Id, cancellationToken);
            throw new InvalidOperationException("Registration failed: could not sync user profile. Please try again.");
        }

        return CreateAuthResponse(createdUser);
    }

    private async Task SyncUserToFollowerServiceAsync(User user, CancellationToken cancellationToken)
    {
        var client = _httpClientFactory.CreateClient("follower");
        var payload = new { userId = user.Id.ToString(), username = user.Username, role = user.Role };
        var content = new StringContent(JsonSerializer.Serialize(payload), System.Text.Encoding.UTF8, "application/json");
        var response = await client.PostAsync("api/followers/users", content, cancellationToken);
        response.EnsureSuccessStatusCode();
        _logger.LogInformation("Registration SAGA T2 succeeded userId={UserId} username={Username} synced to follower service", user.Id, user.Username);
    }

    public async Task<AuthResponseDto> LoginAsync(LoginRequestDto request, CancellationToken cancellationToken = default)
    {
        var username = request.Username.Trim();
        var user = await _userRepository.GetByUsernameAsync(username, cancellationToken);

        if (user is null)
        {
            throw new UnauthorizedAccessException("Invalid username or password.");
        }

        if (user.IsBlocked)
        {
            throw new UnauthorizedAccessException("This account is blocked.");
        }

        var verificationResult = _passwordHasher.VerifyHashedPassword(user, user.PasswordHash, request.Password);

        if (verificationResult == PasswordVerificationResult.Failed)
        {
            throw new UnauthorizedAccessException("Invalid username or password.");
        }

        return CreateAuthResponse(user);
    }

    private static UserResponseDto MapToResponse(User user)
    {
        return new UserResponseDto
        {
            Id = user.Id,
            Username = user.Username,
            Email = user.Email,
            Role = user.Role,
            IsBlocked = user.IsBlocked,
            FirstName = user.FirstName,
            LastName = user.LastName,
            ProfileImage = user.ProfileImage,
            Biography = user.Biography,
            Motto = user.Motto
        };
    }

    private AuthResponseDto CreateAuthResponse(User user)
    {
        return new AuthResponseDto
        {
            Token = GenerateToken(user),
            User = MapToResponse(user)
        };
    }

    private string GenerateToken(User user)
    {
        var claims = new List<Claim>
        {
            new(JwtRegisteredClaimNames.Sub, user.Id.ToString()),
            new(JwtRegisteredClaimNames.UniqueName, user.Username),
            new(ClaimTypes.NameIdentifier, user.Id.ToString()),
            new(ClaimTypes.Name, user.Username),
            new(ClaimTypes.Email, user.Email),
            new(ClaimTypes.Role, user.Role)
        };

        var key = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(_jwtSettings.Key));
        var credentials = new SigningCredentials(key, SecurityAlgorithms.HmacSha256);
        var expires = DateTime.UtcNow.AddMinutes(_jwtSettings.ExpiresInMinutes);

        var token = new JwtSecurityToken(
            issuer: _jwtSettings.Issuer,
            audience: _jwtSettings.Audience,
            claims: claims,
            expires: expires,
            signingCredentials: credentials);

        return new JwtSecurityTokenHandler().WriteToken(token);
    }
}

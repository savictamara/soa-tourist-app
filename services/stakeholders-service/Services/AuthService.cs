using Microsoft.AspNetCore.Identity;
using StakeholdersService.DTOs;
using StakeholdersService.Entities;
using StakeholdersService.Exceptions;
using StakeholdersService.Repositories;

namespace StakeholdersService.Services;

public class AuthService : IAuthService
{
    private static readonly string[] AllowedRoles = ["Guide", "Tourist"];

    private readonly IUserRepository _userRepository;
    private readonly IPasswordHasher<User> _passwordHasher;

    public AuthService(IUserRepository userRepository)
    {
        _userRepository = userRepository;
        _passwordHasher = new PasswordHasher<User>();
    }

    public async Task<UserResponseDto> RegisterAsync(RegisterUserRequestDto request, CancellationToken cancellationToken = default)
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
            IsBlocked = false
        };

        user.PasswordHash = _passwordHasher.HashPassword(user, request.Password);

        var createdUser = await _userRepository.AddAsync(user, cancellationToken);

        return MapToResponse(createdUser);
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
}

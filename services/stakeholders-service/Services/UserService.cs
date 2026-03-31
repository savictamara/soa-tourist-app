using StakeholdersService.DTOs;
using StakeholdersService.Repositories;

namespace StakeholdersService.Services;

public class UserService : IUserService
{
    private const string AdministratorRole = "Administrator";
    private readonly IUserRepository _userRepository;

    public UserService(IUserRepository userRepository)
    {
        _userRepository = userRepository;
    }

    public async Task<List<UserResponseDto>> GetAllAsync(string adminUsername, CancellationToken cancellationToken = default)
    {
        var adminUser = await _userRepository.GetByUsernameAsync(adminUsername, cancellationToken);

        if (adminUser is null || adminUser.Role != AdministratorRole)
        {
            throw new UnauthorizedAccessException("Only an administrator can view all user accounts.");
        }

        var users = await _userRepository.GetAllAsync(cancellationToken);

        return users.Select(user => new UserResponseDto
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
        }).ToList();
    }

    public async Task BlockAsync(long id, string adminUsername, CancellationToken cancellationToken = default)
    {
        var adminUser = await _userRepository.GetByUsernameAsync(adminUsername, cancellationToken);

        if (adminUser is null || adminUser.Role != AdministratorRole)
        {
            throw new UnauthorizedAccessException("Only an administrator can block user accounts.");
        }

        var user = await _userRepository.GetByIdAsync(id, cancellationToken);

        if (user is null)
        {
            throw new KeyNotFoundException($"User with id {id} was not found.");
        }

        if (!user.IsBlocked)
        {
            user.IsBlocked = true;
            await _userRepository.SaveChangesAsync(cancellationToken);
        }
    }

    public async Task<UserProfileResponseDto> GetProfileAsync(long id, CancellationToken cancellationToken = default)
    {
        var user = await _userRepository.GetByIdAsync(id, cancellationToken);

        if (user is null)
        {
            throw new KeyNotFoundException($"User with id {id} was not found.");
        }

        return new UserProfileResponseDto
        {
            FirstName = user.FirstName,
            LastName = user.LastName,
            ProfileImage = user.ProfileImage,
            Biography = user.Biography,
            Motto = user.Motto
        };
    }

    public async Task<UserProfileResponseDto> GetProfileByUsernameAsync(string username, CancellationToken cancellationToken = default)
    {
        var user = await _userRepository.GetByUsernameAsync(username, cancellationToken);

        if (user is null)
        {
            throw new KeyNotFoundException($"User '{username}' was not found.");
        }

        return new UserProfileResponseDto
        {
            FirstName = user.FirstName,
            LastName = user.LastName,
            ProfileImage = user.ProfileImage,
            Biography = user.Biography,
            Motto = user.Motto
        };
    }

    public async Task<UserProfileResponseDto> UpdateProfileAsync(string username, UpdateUserProfileRequestDto request, CancellationToken cancellationToken = default)
    {
        var user = await _userRepository.GetByUsernameAsync(username, cancellationToken);

        if (user is null)
        {
            throw new KeyNotFoundException($"User '{username}' was not found.");
        }

        user.FirstName = request.FirstName?.Trim();
        user.LastName = request.LastName?.Trim();
        user.ProfileImage = request.ProfileImage?.Trim();
        user.Biography = request.Biography?.Trim();
        user.Motto = request.Motto?.Trim();

        await _userRepository.SaveChangesAsync(cancellationToken);

        return new UserProfileResponseDto
        {
            FirstName = user.FirstName,
            LastName = user.LastName,
            ProfileImage = user.ProfileImage,
            Biography = user.Biography,
            Motto = user.Motto
        };
    }
}

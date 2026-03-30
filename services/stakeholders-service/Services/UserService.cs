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
}

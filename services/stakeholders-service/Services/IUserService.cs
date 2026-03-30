using StakeholdersService.DTOs;

namespace StakeholdersService.Services;

public interface IUserService
{
    Task<List<UserResponseDto>> GetAllAsync(string adminUsername, CancellationToken cancellationToken = default);
}

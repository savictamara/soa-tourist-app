using StakeholdersService.DTOs;

namespace StakeholdersService.Services;

public interface IUserService
{
    Task<List<UserResponseDto>> GetAllAsync(string adminUsername, CancellationToken cancellationToken = default);
    Task BlockAsync(long id, string adminUsername, CancellationToken cancellationToken = default);
    Task<UserProfileResponseDto> GetProfileAsync(long id, CancellationToken cancellationToken = default);
    Task<UserProfileResponseDto> GetProfileByUsernameAsync(string username, CancellationToken cancellationToken = default);
    Task<UserProfileResponseDto> UpdateProfileAsync(string username, UpdateUserProfileRequestDto request, CancellationToken cancellationToken = default);
}

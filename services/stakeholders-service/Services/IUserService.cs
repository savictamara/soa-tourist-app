using StakeholdersService.DTOs;

namespace StakeholdersService.Services;

public interface IUserService
{
    Task<List<PublicUserDto>> GetPublicUsersAsync(CancellationToken cancellationToken = default);
    Task<List<UserResponseDto>> GetAllAsync(string adminUsername, CancellationToken cancellationToken = default);
    Task BlockAsync(long id, string adminUsername, CancellationToken cancellationToken = default);
    Task<UserProfileResponseDto> GetProfileAsync(long id, CancellationToken cancellationToken = default);
    Task<UserProfileResponseDto> GetProfileByUsernameAsync(string username, CancellationToken cancellationToken = default);
    Task<UserProfileResponseDto> UpdateProfileAsync(string username, UpdateUserProfileRequestDto request, CancellationToken cancellationToken = default);
    Task<TouristPositionResponseDto> GetTouristPositionAsync(string username, CancellationToken cancellationToken = default);
    Task<TouristPositionResponseDto> UpdateTouristPositionAsync(string username, UpdateTouristPositionRequestDto request, CancellationToken cancellationToken = default);
}

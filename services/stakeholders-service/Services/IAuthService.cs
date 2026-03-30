using StakeholdersService.DTOs;

namespace StakeholdersService.Services;

public interface IAuthService
{
    Task<UserResponseDto> RegisterAsync(RegisterUserRequestDto request, CancellationToken cancellationToken = default);
}

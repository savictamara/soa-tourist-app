using StakeholdersService.DTOs;

namespace StakeholdersService.Services;

public interface IHealthService
{
    HealthResponseDto GetStatus();
}

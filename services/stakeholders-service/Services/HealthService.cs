using StakeholdersService.DTOs;
using StakeholdersService.Repositories;

namespace StakeholdersService.Services;

public class HealthService : IHealthService
{
    private readonly IHealthRepository _healthRepository;

    public HealthService(IHealthRepository healthRepository)
    {
        _healthRepository = healthRepository;
    }

    public HealthResponseDto GetStatus()
    {
        var healthStatus = _healthRepository.GetStatus();

        return new HealthResponseDto
        {
            Service = healthStatus.ServiceName,
            Status = healthStatus.Status
        };
    }
}

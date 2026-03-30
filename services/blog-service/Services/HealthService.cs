using BlogService.DTOs;
using BlogService.Repositories;

namespace BlogService.Services;

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

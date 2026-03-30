using BlogService.DTOs;

namespace BlogService.Services;

public interface IHealthService
{
    HealthResponseDto GetStatus();
}

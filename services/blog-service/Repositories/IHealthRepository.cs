using BlogService.Models;

namespace BlogService.Repositories;

public interface IHealthRepository
{
    HealthStatus GetStatus();
}

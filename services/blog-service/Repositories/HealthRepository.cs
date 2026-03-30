using BlogService.Models;

namespace BlogService.Repositories;

public class HealthRepository : IHealthRepository
{
    public HealthStatus GetStatus()
    {
        return new HealthStatus
        {
            ServiceName = "blog-service",
            Status = "ok"
        };
    }
}

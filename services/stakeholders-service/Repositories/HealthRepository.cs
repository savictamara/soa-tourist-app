using StakeholdersService.Models;

namespace StakeholdersService.Repositories;

public class HealthRepository : IHealthRepository
{
    public HealthStatus GetStatus()
    {
        return new HealthStatus
        {
            ServiceName = "stakeholders-service",
            Status = "ok"
        };
    }
}

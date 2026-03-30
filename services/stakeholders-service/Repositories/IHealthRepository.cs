using StakeholdersService.Models;

namespace StakeholdersService.Repositories;

public interface IHealthRepository
{
    HealthStatus GetStatus();
}

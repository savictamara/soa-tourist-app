using Microsoft.EntityFrameworkCore;
using StakeholdersService.Entities;

namespace StakeholdersService.Data;

public class StakeholdersDbContext : DbContext
{
    public StakeholdersDbContext(DbContextOptions<StakeholdersDbContext> options) : base(options)
    {
    }

    public DbSet<User> Users => Set<User>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);
        modelBuilder.ApplyConfigurationsFromAssembly(typeof(StakeholdersDbContext).Assembly);
    }
}

using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Design;

namespace StakeholdersService.Data;

public class StakeholdersDbContextFactory : IDesignTimeDbContextFactory<StakeholdersDbContext>
{
    public StakeholdersDbContext CreateDbContext(string[] args)
    {
        var environmentName = Environment.GetEnvironmentVariable("ASPNETCORE_ENVIRONMENT");

        var configuration = new ConfigurationBuilder()
            .SetBasePath(Directory.GetCurrentDirectory())
            .AddJsonFile("appsettings.json", optional: false)
            .AddJsonFile($"appsettings.{environmentName}.json", optional: true)
            .AddEnvironmentVariables()
            .Build();

        var connectionString = configuration.GetConnectionString("DefaultConnection")
            ?? throw new InvalidOperationException("Connection string 'DefaultConnection' was not found.");

        var optionsBuilder = new DbContextOptionsBuilder<StakeholdersDbContext>();
        optionsBuilder.UseNpgsql(connectionString);

        return new StakeholdersDbContext(optionsBuilder.Options);
    }
}

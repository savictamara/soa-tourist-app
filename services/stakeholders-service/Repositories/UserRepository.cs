using Microsoft.EntityFrameworkCore;
using StakeholdersService.Data;
using StakeholdersService.Entities;

namespace StakeholdersService.Repositories;

public class UserRepository : IUserRepository
{
    private readonly StakeholdersDbContext _dbContext;

    public UserRepository(StakeholdersDbContext dbContext)
    {
        _dbContext = dbContext;
    }

    public Task<bool> UsernameExistsAsync(string username, CancellationToken cancellationToken = default)
    {
        var normalizedUsername = username.Trim().ToLower();

        return _dbContext.Users.AnyAsync(
            user => user.Username.ToLower() == normalizedUsername,
            cancellationToken);
    }

    public Task<bool> EmailExistsAsync(string email, CancellationToken cancellationToken = default)
    {
        var normalizedEmail = email.Trim().ToLower();

        return _dbContext.Users.AnyAsync(
            user => user.Email.ToLower() == normalizedEmail,
            cancellationToken);
    }

    public async Task<User> AddAsync(User user, CancellationToken cancellationToken = default)
    {
        _dbContext.Users.Add(user);
        await _dbContext.SaveChangesAsync(cancellationToken);
        return user;
    }

    public Task<User?> GetByUsernameAsync(string username, CancellationToken cancellationToken = default)
    {
        var normalizedUsername = username.Trim().ToLower();

        return _dbContext.Users.FirstOrDefaultAsync(
            user => user.Username.ToLower() == normalizedUsername,
            cancellationToken);
    }

    public Task<User?> GetByIdAsync(long id, CancellationToken cancellationToken = default)
    {
        return _dbContext.Users.FirstOrDefaultAsync(user => user.Id == id, cancellationToken);
    }

    public Task<List<User>> GetAllAsync(CancellationToken cancellationToken = default)
    {
        return _dbContext.Users
            .OrderBy(user => user.Id)
            .ToListAsync(cancellationToken);
    }

    public Task SaveChangesAsync(CancellationToken cancellationToken = default)
    {
        return _dbContext.SaveChangesAsync(cancellationToken);
    }
}

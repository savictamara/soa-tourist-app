using BlogService.Data;
using BlogService.Entities;
using Microsoft.EntityFrameworkCore;

namespace BlogService.Repositories;

public class BlogRepository : IBlogRepository
{
    private readonly BlogDbContext _dbContext;

    public BlogRepository(BlogDbContext dbContext)
    {
        _dbContext = dbContext;
    }

    public async Task<BlogPost> AddAsync(BlogPost blogPost, CancellationToken cancellationToken = default)
    {
        _dbContext.BlogPosts.Add(blogPost);
        await _dbContext.SaveChangesAsync(cancellationToken);

        return await _dbContext.BlogPosts
            .Include(post => post.Images)
            .SingleAsync(post => post.Id == blogPost.Id, cancellationToken);
    }

    public Task<List<BlogPost>> GetAllAsync(CancellationToken cancellationToken = default)
    {
        return _dbContext.BlogPosts
            .Include(post => post.Images)
            .OrderByDescending(post => post.CreatedAtUtc)
            .ThenByDescending(post => post.Id)
            .ToListAsync(cancellationToken);
    }
}

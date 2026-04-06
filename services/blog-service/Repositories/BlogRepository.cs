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
            .Include(post => post.Comments.OrderByDescending(comment => comment.CreatedAtUtc).ThenByDescending(comment => comment.Id))
            .SingleAsync(post => post.Id == blogPost.Id, cancellationToken);
    }

    public Task<List<BlogPost>> GetAllAsync(CancellationToken cancellationToken = default)
    {
        return _dbContext.BlogPosts
            .Include(post => post.Images)
            .Include(post => post.Comments.OrderByDescending(comment => comment.CreatedAtUtc).ThenByDescending(comment => comment.Id))
            .OrderByDescending(post => post.CreatedAtUtc)
            .ThenByDescending(post => post.Id)
            .ToListAsync(cancellationToken);
    }

    public Task<BlogPost?> GetByIdAsync(long blogPostId, CancellationToken cancellationToken = default)
    {
        return _dbContext.BlogPosts
            .Include(post => post.Images)
            .Include(post => post.Comments.OrderByDescending(comment => comment.CreatedAtUtc).ThenByDescending(comment => comment.Id))
            .SingleOrDefaultAsync(post => post.Id == blogPostId, cancellationToken);
    }

    public async Task<BlogComment> AddCommentAsync(BlogComment comment, CancellationToken cancellationToken = default)
    {
        _dbContext.BlogComments.Add(comment);
        await _dbContext.SaveChangesAsync(cancellationToken);

        return await _dbContext.BlogComments
            .SingleAsync(existingComment => existingComment.Id == comment.Id, cancellationToken);
    }

    public Task<BlogComment?> GetCommentByIdAsync(long blogPostId, long commentId, CancellationToken cancellationToken = default)
    {
        return _dbContext.BlogComments
            .SingleOrDefaultAsync(
                comment => comment.BlogPostId == blogPostId && comment.Id == commentId,
                cancellationToken);
    }

    public Task SaveChangesAsync(CancellationToken cancellationToken = default)
    {
        return _dbContext.SaveChangesAsync(cancellationToken);
    }
}

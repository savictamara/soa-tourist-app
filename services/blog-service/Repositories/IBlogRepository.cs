using BlogService.Entities;

namespace BlogService.Repositories;

public interface IBlogRepository
{
    Task<BlogPost> AddAsync(BlogPost blogPost, CancellationToken cancellationToken = default);
    Task<List<BlogPost>> GetAllAsync(CancellationToken cancellationToken = default);
    Task<BlogPost?> GetByIdAsync(long blogPostId, CancellationToken cancellationToken = default);
    Task<BlogComment> AddCommentAsync(BlogComment comment, CancellationToken cancellationToken = default);
    Task<BlogComment?> GetCommentByIdAsync(long blogPostId, long commentId, CancellationToken cancellationToken = default);
    Task<BlogLike?> GetLikeAsync(long blogPostId, string username, CancellationToken cancellationToken = default);
    Task<BlogLike> AddLikeAsync(BlogLike like, CancellationToken cancellationToken = default);
    Task RemoveLikeAsync(BlogLike like, CancellationToken cancellationToken = default);
    Task SaveChangesAsync(CancellationToken cancellationToken = default);
}

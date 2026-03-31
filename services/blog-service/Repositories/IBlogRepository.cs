using BlogService.Entities;

namespace BlogService.Repositories;

public interface IBlogRepository
{
    Task<BlogPost> AddAsync(BlogPost blogPost, CancellationToken cancellationToken = default);
    Task<List<BlogPost>> GetAllAsync(CancellationToken cancellationToken = default);
}

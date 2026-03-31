using BlogService.DTOs;

namespace BlogService.Services;

public interface IBlogService
{
    Task<List<BlogPostResponseDto>> GetAllAsync(CancellationToken cancellationToken = default);

    Task<BlogPostResponseDto> CreateAsync(
        string username,
        string role,
        CreateBlogPostRequestDto request,
        CancellationToken cancellationToken = default);
}

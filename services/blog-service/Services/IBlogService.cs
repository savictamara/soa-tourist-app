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

    Task<BlogCommentResponseDto> AddCommentAsync(
        long blogPostId,
        string username,
        string role,
        CreateBlogCommentRequestDto request,
        CancellationToken cancellationToken = default);

    Task<BlogCommentResponseDto> UpdateCommentAsync(
        long blogPostId,
        long commentId,
        string username,
        UpdateBlogCommentRequestDto request,
        CancellationToken cancellationToken = default);
}

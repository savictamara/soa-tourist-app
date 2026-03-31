using BlogService.DTOs;
using BlogService.Entities;
using BlogService.Repositories;

namespace BlogService.Services;

public class BlogService : IBlogService
{
    private static readonly HashSet<string> AllowedRoles = ["Guide", "Tourist"];
    private readonly IBlogRepository _blogRepository;

    public BlogService(IBlogRepository blogRepository)
    {
        _blogRepository = blogRepository;
    }

    public async Task<List<BlogPostResponseDto>> GetAllAsync(CancellationToken cancellationToken = default)
    {
        var posts = await _blogRepository.GetAllAsync(cancellationToken);
        return posts.Select(MapToResponse).ToList();
    }

    public async Task<BlogPostResponseDto> CreateAsync(
        string username,
        string role,
        CreateBlogPostRequestDto request,
        CancellationToken cancellationToken = default)
    {
        if (!AllowedRoles.Contains(role))
        {
            throw new UnauthorizedAccessException("Only Guide and Tourist users can create blogs.");
        }

        var title = request.Title.Trim();
        var descriptionMarkdown = request.DescriptionMarkdown.Trim();

        if (string.IsNullOrWhiteSpace(title))
        {
            throw new ArgumentException("Title is required.");
        }

        if (string.IsNullOrWhiteSpace(descriptionMarkdown))
        {
            throw new ArgumentException("Description markdown is required.");
        }

        var imageUrls = request.ImageUrls?
            .Select(url => url.Trim())
            .Where(url => !string.IsNullOrWhiteSpace(url))
            .Distinct()
            .ToList() ?? [];

        var blogPost = new BlogPost
        {
            Title = title,
            DescriptionMarkdown = descriptionMarkdown,
            CreatedAtUtc = DateTime.UtcNow,
            AuthorUsername = username,
            Images = imageUrls
                .Select(url => new BlogPostImage { ImageUrl = url })
                .ToList()
        };

        var createdBlogPost = await _blogRepository.AddAsync(blogPost, cancellationToken);

        return MapToResponse(createdBlogPost);
    }

    private static BlogPostResponseDto MapToResponse(BlogPost post)
    {
        return new BlogPostResponseDto
        {
            Id = post.Id,
            Title = post.Title,
            DescriptionMarkdown = post.DescriptionMarkdown,
            CreatedAtUtc = post.CreatedAtUtc,
            AuthorUsername = post.AuthorUsername,
            ImageUrls = post.Images.Select(image => image.ImageUrl).ToList()
        };
    }
}

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

    public async Task<List<BlogPostResponseDto>> GetAllAsync(string currentUsername, CancellationToken cancellationToken = default)
    {
        var posts = await _blogRepository.GetAllAsync(cancellationToken);
        return posts.Select(post => MapToResponse(post, currentUsername)).ToList();
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

        return MapToResponse(createdBlogPost, username);
    }

    public async Task<BlogCommentResponseDto> AddCommentAsync(
        long blogPostId,
        string username,
        string role,
        CreateBlogCommentRequestDto request,
        CancellationToken cancellationToken = default)
    {
        EnsureAllowedRole(role, "Only Guide and Tourist users can comment on blogs.");

        var blogPost = await _blogRepository.GetByIdAsync(blogPostId, cancellationToken);

        if (blogPost is null)
        {
            throw new KeyNotFoundException("Blog post not found.");
        }

        var text = request.Text.Trim();

        if (string.IsNullOrWhiteSpace(text))
        {
            throw new ArgumentException("Comment text is required.");
        }

        var now = DateTime.UtcNow;
        var comment = new BlogComment
        {
            BlogPostId = blogPostId,
            AuthorUsername = username,
            AuthorRole = role,
            Text = text,
            CreatedAtUtc = now,
            LastModifiedAtUtc = now
        };

        var createdComment = await _blogRepository.AddCommentAsync(comment, cancellationToken);
        return MapCommentToResponse(createdComment);
    }

    public async Task<BlogCommentResponseDto> UpdateCommentAsync(
        long blogPostId,
        long commentId,
        string username,
        UpdateBlogCommentRequestDto request,
        CancellationToken cancellationToken = default)
    {
        var comment = await _blogRepository.GetCommentByIdAsync(blogPostId, commentId, cancellationToken);

        if (comment is null)
        {
            throw new KeyNotFoundException("Comment not found.");
        }

        if (!string.Equals(comment.AuthorUsername, username, StringComparison.Ordinal))
        {
            throw new UnauthorizedAccessException("You can only edit your own comments.");
        }

        var text = request.Text.Trim();

        if (string.IsNullOrWhiteSpace(text))
        {
            throw new ArgumentException("Comment text is required.");
        }

        comment.Text = text;
        comment.LastModifiedAtUtc = DateTime.UtcNow;
        await _blogRepository.SaveChangesAsync(cancellationToken);

        return MapCommentToResponse(comment);
    }

    public async Task<BlogLikeStatusResponseDto> LikeAsync(
        long blogPostId,
        string username,
        string role,
        CancellationToken cancellationToken = default)
    {
        EnsureAllowedRole(role, "Only Guide and Tourist users can like blogs.");

        var blogPost = await _blogRepository.GetByIdAsync(blogPostId, cancellationToken);

        if (blogPost is null)
        {
            throw new KeyNotFoundException("Blog post not found.");
        }

        var existingLike = await _blogRepository.GetLikeAsync(blogPostId, username, cancellationToken);

        if (existingLike is not null)
        {
            return CreateLikeStatus(blogPost, username);
        }

        await _blogRepository.AddLikeAsync(new BlogLike
        {
            BlogPostId = blogPostId,
            Username = username,
            CreatedAtUtc = DateTime.UtcNow
        }, cancellationToken);

        blogPost = await _blogRepository.GetByIdAsync(blogPostId, cancellationToken)
            ?? throw new KeyNotFoundException("Blog post not found.");

        return CreateLikeStatus(blogPost, username);
    }

    public async Task<BlogLikeStatusResponseDto> UnlikeAsync(
        long blogPostId,
        string username,
        string role,
        CancellationToken cancellationToken = default)
    {
        EnsureAllowedRole(role, "Only Guide and Tourist users can like blogs.");

        var blogPost = await _blogRepository.GetByIdAsync(blogPostId, cancellationToken);

        if (blogPost is null)
        {
            throw new KeyNotFoundException("Blog post not found.");
        }

        var existingLike = await _blogRepository.GetLikeAsync(blogPostId, username, cancellationToken);

        if (existingLike is not null)
        {
            await _blogRepository.RemoveLikeAsync(existingLike, cancellationToken);
        }

        blogPost = await _blogRepository.GetByIdAsync(blogPostId, cancellationToken)
            ?? throw new KeyNotFoundException("Blog post not found.");

        return CreateLikeStatus(blogPost, username);
    }

    private static void EnsureAllowedRole(string role, string message)
    {
        if (!AllowedRoles.Contains(role))
        {
            throw new UnauthorizedAccessException(message);
        }
    }

    private static BlogPostResponseDto MapToResponse(BlogPost post, string currentUsername)
    {
        return new BlogPostResponseDto
        {
            Id = post.Id,
            Title = post.Title,
            DescriptionMarkdown = post.DescriptionMarkdown,
            CreatedAtUtc = post.CreatedAtUtc,
            AuthorUsername = post.AuthorUsername,
            ImageUrls = post.Images.Select(image => image.ImageUrl).ToList(),
            LikesCount = post.Likes.Count,
            IsLikedByCurrentUser = post.Likes.Any(like => like.Username == currentUsername),
            Comments = post.Comments
                .OrderByDescending(comment => comment.CreatedAtUtc)
                .ThenByDescending(comment => comment.Id)
                .Select(MapCommentToResponse)
                .ToList()
        };
    }

    private static BlogLikeStatusResponseDto CreateLikeStatus(BlogPost blogPost, string currentUsername)
    {
        return new BlogLikeStatusResponseDto
        {
            BlogPostId = blogPost.Id,
            LikesCount = blogPost.Likes.Count,
            IsLikedByCurrentUser = blogPost.Likes.Any(like => like.Username == currentUsername)
        };
    }

    private static BlogCommentResponseDto MapCommentToResponse(BlogComment comment)
    {
        return new BlogCommentResponseDto
        {
            Id = comment.Id,
            AuthorUsername = comment.AuthorUsername,
            AuthorRole = comment.AuthorRole,
            Text = comment.Text,
            CreatedAtUtc = comment.CreatedAtUtc,
            LastModifiedAtUtc = comment.LastModifiedAtUtc
        };
    }
}

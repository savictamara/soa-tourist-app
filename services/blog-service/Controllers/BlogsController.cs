using System.Security.Claims;
using BlogService.DTOs;
using BlogService.Services;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;

namespace BlogService.Controllers;

[ApiController]
[Route("api/blogs")]
[Authorize]
public class BlogsController : ControllerBase
{
    private readonly IBlogService _blogService;

    public BlogsController(IBlogService blogService)
    {
        _blogService = blogService;
    }

    [HttpGet]
    [Authorize(Roles = "Guide,Tourist")]
    [ProducesResponseType(typeof(List<BlogPostResponseDto>), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status401Unauthorized)]
    public async Task<ActionResult<List<BlogPostResponseDto>>> GetAll(CancellationToken cancellationToken)
    {
        var username = User.FindFirstValue(ClaimTypes.Name);

        if (string.IsNullOrWhiteSpace(username))
        {
            return Unauthorized(new { message = "Current user is not authenticated." });
        }

        var posts = await _blogService.GetAllAsync(username, cancellationToken);
        return Ok(posts);
    }

    [HttpGet("followed/{userId:long}")]
    [Authorize(Roles = "Guide,Tourist")]
    [ProducesResponseType(typeof(List<BlogPostResponseDto>), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status401Unauthorized)]
    public async Task<ActionResult<List<BlogPostResponseDto>>> GetFollowed(long userId, CancellationToken cancellationToken)
    {
        var username = User.FindFirstValue(ClaimTypes.Name);

        if (string.IsNullOrWhiteSpace(username))
        {
            return Unauthorized(new { message = "Current user is not authenticated." });
        }

        var posts = await _blogService.GetFollowedBlogsAsync(userId, username, cancellationToken);
        return Ok(posts);
    }

    [HttpPost]
    [Authorize(Roles = "Guide,Tourist")]
    [ProducesResponseType(typeof(BlogPostResponseDto), StatusCodes.Status201Created)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    [ProducesResponseType(StatusCodes.Status401Unauthorized)]
    public async Task<ActionResult<BlogPostResponseDto>> Create(
        [FromBody] CreateBlogPostRequestDto request,
        CancellationToken cancellationToken)
    {
        var username = User.FindFirstValue(ClaimTypes.Name);
        var role = User.FindFirstValue(ClaimTypes.Role);
        var userIdStr = User.FindFirstValue(ClaimTypes.NameIdentifier);

        if (string.IsNullOrWhiteSpace(username) || string.IsNullOrWhiteSpace(role))
        {
            return Unauthorized(new { message = "Current user is not authenticated." });
        }

        if (!long.TryParse(userIdStr, out var authorId))
        {
            return Unauthorized(new { message = "Current user is not authenticated." });
        }

        try
        {
            var createdBlogPost = await _blogService.CreateAsync(username, role, authorId, request, cancellationToken);
            return Created($"/api/blogs/{createdBlogPost.Id}", createdBlogPost);
        }
        catch (ArgumentException exception)
        {
            return BadRequest(new { message = exception.Message });
        }
        catch (UnauthorizedAccessException exception)
        {
            return Unauthorized(new { message = exception.Message });
        }
    }

    [HttpPost("{blogId:long}/likes")]
    [Authorize(Roles = "Guide,Tourist")]
    [ProducesResponseType(typeof(BlogLikeStatusResponseDto), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status401Unauthorized)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<ActionResult<BlogLikeStatusResponseDto>> Like(
        long blogId,
        CancellationToken cancellationToken)
    {
        var username = User.FindFirstValue(ClaimTypes.Name);
        var role = User.FindFirstValue(ClaimTypes.Role);

        if (string.IsNullOrWhiteSpace(username) || string.IsNullOrWhiteSpace(role))
        {
            return Unauthorized(new { message = "Current user is not authenticated." });
        }

        try
        {
            var status = await _blogService.LikeAsync(blogId, username, role, cancellationToken);
            return Ok(status);
        }
        catch (KeyNotFoundException exception)
        {
            return NotFound(new { message = exception.Message });
        }
        catch (UnauthorizedAccessException exception)
        {
            return Unauthorized(new { message = exception.Message });
        }
    }

    [HttpDelete("{blogId:long}/likes")]
    [Authorize(Roles = "Guide,Tourist")]
    [ProducesResponseType(typeof(BlogLikeStatusResponseDto), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status401Unauthorized)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<ActionResult<BlogLikeStatusResponseDto>> Unlike(
        long blogId,
        CancellationToken cancellationToken)
    {
        var username = User.FindFirstValue(ClaimTypes.Name);
        var role = User.FindFirstValue(ClaimTypes.Role);

        if (string.IsNullOrWhiteSpace(username) || string.IsNullOrWhiteSpace(role))
        {
            return Unauthorized(new { message = "Current user is not authenticated." });
        }

        try
        {
            var status = await _blogService.UnlikeAsync(blogId, username, role, cancellationToken);
            return Ok(status);
        }
        catch (KeyNotFoundException exception)
        {
            return NotFound(new { message = exception.Message });
        }
        catch (UnauthorizedAccessException exception)
        {
            return Unauthorized(new { message = exception.Message });
        }
    }

    [HttpPost("{blogId:long}/comments")]
    [Authorize(Roles = "Guide,Tourist")]
    [ProducesResponseType(typeof(BlogCommentResponseDto), StatusCodes.Status201Created)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    [ProducesResponseType(StatusCodes.Status401Unauthorized)]
    [ProducesResponseType(StatusCodes.Status403Forbidden)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<ActionResult<BlogCommentResponseDto>> AddComment(
        long blogId,
        [FromBody] CreateBlogCommentRequestDto request,
        CancellationToken cancellationToken)
    {
        var username = User.FindFirstValue(ClaimTypes.Name);
        var role = User.FindFirstValue(ClaimTypes.Role);
        var userIdStr = User.FindFirstValue(ClaimTypes.NameIdentifier);

        if (string.IsNullOrWhiteSpace(username) || string.IsNullOrWhiteSpace(role))
        {
            return Unauthorized(new { message = "Current user is not authenticated." });
        }

        if (!long.TryParse(userIdStr, out var commenterId))
        {
            return Unauthorized(new { message = "Current user is not authenticated." });
        }

        try
        {
            var createdComment = await _blogService.AddCommentAsync(blogId, username, role, commenterId, request, cancellationToken);
            return Created($"/api/blogs/{blogId}/comments/{createdComment.Id}", createdComment);
        }
        catch (ArgumentException exception)
        {
            return BadRequest(new { message = exception.Message });
        }
        catch (KeyNotFoundException exception)
        {
            return NotFound(new { message = exception.Message });
        }
        catch (UnauthorizedAccessException exception)
        {
            return StatusCode(StatusCodes.Status403Forbidden, new { message = exception.Message });
        }
    }

    [HttpPut("{blogId:long}/comments/{commentId:long}")]
    [Authorize(Roles = "Guide,Tourist")]
    [ProducesResponseType(typeof(BlogCommentResponseDto), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    [ProducesResponseType(StatusCodes.Status401Unauthorized)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<ActionResult<BlogCommentResponseDto>> UpdateComment(
        long blogId,
        long commentId,
        [FromBody] UpdateBlogCommentRequestDto request,
        CancellationToken cancellationToken)
    {
        var username = User.FindFirstValue(ClaimTypes.Name);

        if (string.IsNullOrWhiteSpace(username))
        {
            return Unauthorized(new { message = "Current user is not authenticated." });
        }

        try
        {
            var updatedComment = await _blogService.UpdateCommentAsync(blogId, commentId, username, request, cancellationToken);
            return Ok(updatedComment);
        }
        catch (ArgumentException exception)
        {
            return BadRequest(new { message = exception.Message });
        }
        catch (KeyNotFoundException exception)
        {
            return NotFound(new { message = exception.Message });
        }
        catch (UnauthorizedAccessException exception)
        {
            return Unauthorized(new { message = exception.Message });
        }
    }
}

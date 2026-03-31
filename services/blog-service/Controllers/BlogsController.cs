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
        var posts = await _blogService.GetAllAsync(cancellationToken);
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

        if (string.IsNullOrWhiteSpace(username) || string.IsNullOrWhiteSpace(role))
        {
            return Unauthorized(new { message = "Current user is not authenticated." });
        }

        try
        {
            var createdBlogPost = await _blogService.CreateAsync(username, role, request, cancellationToken);
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
}

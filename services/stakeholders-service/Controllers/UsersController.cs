using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using System.Security.Claims;
using StakeholdersService.DTOs;
using StakeholdersService.Services;

namespace StakeholdersService.Controllers;

[ApiController]
[Route("api/users")]
[Authorize]
public class UsersController : ControllerBase
{
    private readonly IUserService _userService;

    public UsersController(IUserService userService)
    {
        _userService = userService;
    }

    [HttpGet("public")]
    [AllowAnonymous]
    [ProducesResponseType(typeof(List<PublicUserDto>), StatusCodes.Status200OK)]
    public async Task<ActionResult<List<PublicUserDto>>> GetPublic(CancellationToken cancellationToken)
    {
        var users = await _userService.GetPublicUsersAsync(cancellationToken);
        return Ok(users);
    }

    [HttpGet]
    [Authorize(Roles = "Administrator")]
    [ProducesResponseType(typeof(List<UserResponseDto>), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status401Unauthorized)]
    public async Task<ActionResult<List<UserResponseDto>>> GetAll(CancellationToken cancellationToken)
    {
        var adminUsername = User.FindFirstValue(ClaimTypes.Name);

        if (string.IsNullOrWhiteSpace(adminUsername))
        {
            return Unauthorized(new { message = "Current user is not authenticated." });
        }

        try
        {
            var users = await _userService.GetAllAsync(adminUsername.ToString(), cancellationToken);
            return Ok(users);
        }
        catch (UnauthorizedAccessException exception)
        {
            return Unauthorized(new { message = exception.Message });
        }
    }

    [HttpPut("{id:long}/block")]
    [Authorize(Roles = "Administrator")]
    [ProducesResponseType(StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status401Unauthorized)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<ActionResult> Block(long id, CancellationToken cancellationToken)
    {
        var adminUsername = User.FindFirstValue(ClaimTypes.Name);

        if (string.IsNullOrWhiteSpace(adminUsername))
        {
            return Unauthorized(new { message = "Current user is not authenticated." });
        }

        try
        {
            await _userService.BlockAsync(id, adminUsername.ToString(), cancellationToken);
            return Ok(new { message = "User blocked successfully." });
        }
        catch (UnauthorizedAccessException exception)
        {
            return Unauthorized(new { message = exception.Message });
        }
        catch (KeyNotFoundException exception)
        {
            return NotFound(new { message = exception.Message });
        }
    }

    [HttpGet("{id:long}/profile")]
    [Authorize(Roles = "Administrator")]
    [ProducesResponseType(typeof(UserProfileResponseDto), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<ActionResult<UserProfileResponseDto>> GetProfile(long id, CancellationToken cancellationToken)
    {
        try
        {
            var profile = await _userService.GetProfileAsync(id, cancellationToken);
            return Ok(profile);
        }
        catch (KeyNotFoundException exception)
        {
            return NotFound(new { message = exception.Message });
        }
    }

    [HttpGet("me/profile")]
    [ProducesResponseType(typeof(UserProfileResponseDto), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status401Unauthorized)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<ActionResult<UserProfileResponseDto>> GetMyProfile(CancellationToken cancellationToken)
    {
        var username = User.FindFirstValue(ClaimTypes.Name);

        if (string.IsNullOrWhiteSpace(username))
        {
            return Unauthorized(new { message = "Current user is not authenticated." });
        }

        try
        {
            var profile = await _userService.GetProfileByUsernameAsync(username, cancellationToken);
            return Ok(profile);
        }
        catch (KeyNotFoundException exception)
        {
            return NotFound(new { message = exception.Message });
        }
    }

    [HttpPut("me/profile")]
    [ProducesResponseType(typeof(UserProfileResponseDto), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    [ProducesResponseType(StatusCodes.Status401Unauthorized)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<ActionResult<UserProfileResponseDto>> UpdateMyProfile(
        [FromBody] UpdateUserProfileRequestDto request,
        CancellationToken cancellationToken)
    {
        var username = User.FindFirstValue(ClaimTypes.Name);

        if (string.IsNullOrWhiteSpace(username))
        {
            return Unauthorized(new { message = "Current user is not authenticated." });
        }

        try
        {
            var profile = await _userService.UpdateProfileAsync(username, request, cancellationToken);
            return Ok(profile);
        }
        catch (KeyNotFoundException exception)
        {
            return NotFound(new { message = exception.Message });
        }
    }

    [HttpGet("me/position")]
    [Authorize(Roles = "Tourist")]
    [ProducesResponseType(typeof(TouristPositionResponseDto), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status401Unauthorized)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<ActionResult<TouristPositionResponseDto>> GetMyPosition(CancellationToken cancellationToken)
    {
        var username = User.FindFirstValue(ClaimTypes.Name);

        if (string.IsNullOrWhiteSpace(username))
        {
            return Unauthorized(new { message = "Current user is not authenticated." });
        }

        try
        {
            var position = await _userService.GetTouristPositionAsync(username, cancellationToken);
            return Ok(position);
        }
        catch (UnauthorizedAccessException exception)
        {
            return Unauthorized(new { message = exception.Message });
        }
        catch (KeyNotFoundException exception)
        {
            return NotFound(new { message = exception.Message });
        }
    }

    [HttpPut("me/position")]
    [Authorize(Roles = "Tourist")]
    [ProducesResponseType(typeof(TouristPositionResponseDto), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    [ProducesResponseType(StatusCodes.Status401Unauthorized)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<ActionResult<TouristPositionResponseDto>> UpdateMyPosition(
        [FromBody] UpdateTouristPositionRequestDto request,
        CancellationToken cancellationToken)
    {
        var username = User.FindFirstValue(ClaimTypes.Name);

        if (string.IsNullOrWhiteSpace(username))
        {
            return Unauthorized(new { message = "Current user is not authenticated." });
        }

        try
        {
            var position = await _userService.UpdateTouristPositionAsync(username, request, cancellationToken);
            return Ok(position);
        }
        catch (ArgumentException exception)
        {
            return BadRequest(new { message = exception.Message });
        }
        catch (UnauthorizedAccessException exception)
        {
            return Unauthorized(new { message = exception.Message });
        }
        catch (KeyNotFoundException exception)
        {
            return NotFound(new { message = exception.Message });
        }
    }
}

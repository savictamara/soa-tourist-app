using Microsoft.AspNetCore.Mvc;
using StakeholdersService.DTOs;
using StakeholdersService.Services;

namespace StakeholdersService.Controllers;

[ApiController]
[Route("api/users")]
public class UsersController : ControllerBase
{
    private readonly IUserService _userService;

    public UsersController(IUserService userService)
    {
        _userService = userService;
    }

    [HttpGet]
    [ProducesResponseType(typeof(List<UserResponseDto>), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status401Unauthorized)]
    public async Task<ActionResult<List<UserResponseDto>>> GetAll(CancellationToken cancellationToken)
    {
        if (!Request.Headers.TryGetValue("X-Admin-Username", out var adminUsername) ||
            string.IsNullOrWhiteSpace(adminUsername))
        {
            return Unauthorized(new { message = "Admin username is required." });
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
}

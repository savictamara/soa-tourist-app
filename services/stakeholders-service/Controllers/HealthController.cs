using Microsoft.AspNetCore.Mvc;
using StakeholdersService.DTOs;
using StakeholdersService.Services;

namespace StakeholdersService.Controllers;

[ApiController]
[Route("api/health")]
public class HealthController : ControllerBase
{
    private readonly IHealthService _healthService;

    public HealthController(IHealthService healthService)
    {
        _healthService = healthService;
    }

    [HttpGet]
    [ProducesResponseType(typeof(HealthResponseDto), StatusCodes.Status200OK)]
    public ActionResult<HealthResponseDto> Get()
    {
        return Ok(_healthService.GetStatus());
    }
}

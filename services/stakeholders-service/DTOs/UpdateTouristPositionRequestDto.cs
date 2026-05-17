using System.ComponentModel.DataAnnotations;

namespace StakeholdersService.DTOs;

public class UpdateTouristPositionRequestDto
{
    [Range(-90, 90)]
    public double? Latitude { get; set; }

    [Range(-180, 180)]
    public double? Longitude { get; set; }
}

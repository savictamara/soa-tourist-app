using System.ComponentModel.DataAnnotations;

namespace StakeholdersService.DTOs;

public class UpdateUserProfileRequestDto
{
    [MaxLength(100)]
    public string? FirstName { get; set; }

    [MaxLength(100)]
    public string? LastName { get; set; }

    public string? ProfileImage { get; set; }

    [MaxLength(4000)]
    public string? Biography { get; set; }

    [MaxLength(250)]
    public string? Motto { get; set; }
}

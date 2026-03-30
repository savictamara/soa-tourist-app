using System.ComponentModel.DataAnnotations;

namespace StakeholdersService.DTOs;

public class RegisterUserRequestDto
{
    [Required]
    public string Username { get; set; } = string.Empty;

    [Required]
    public string Password { get; set; } = string.Empty;

    [Required]
    [EmailAddress]
    public string Email { get; set; } = string.Empty;

    [Required]
    [RegularExpression("^(Guide|Tourist)$", ErrorMessage = "Role must be Guide or Tourist.")]
    public string Role { get; set; } = string.Empty;
}

namespace StakeholdersService.Entities;

public class User
{
    public long Id { get; set; }
    public string Username { get; set; } = string.Empty;
    public string PasswordHash { get; set; } = string.Empty;
    public string Email { get; set; } = string.Empty;
    public string Role { get; set; } = string.Empty;
    public bool IsBlocked { get; set; }
    public string? FirstName { get; set; }
    public string? LastName { get; set; }
    public string? ProfileImage { get; set; }
    public string? Biography { get; set; }
    public string? Motto { get; set; }
    public double? Latitude { get; set; }
    public double? Longitude { get; set; }
}

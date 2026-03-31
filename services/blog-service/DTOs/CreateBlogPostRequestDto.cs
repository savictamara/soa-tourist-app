using System.ComponentModel.DataAnnotations;

namespace BlogService.DTOs;

public class CreateBlogPostRequestDto
{
    [Required]
    [MaxLength(200)]
    public string Title { get; set; } = string.Empty;

    [Required]
    [MaxLength(12000)]
    public string DescriptionMarkdown { get; set; } = string.Empty;

    public List<string>? ImageUrls { get; set; }
}

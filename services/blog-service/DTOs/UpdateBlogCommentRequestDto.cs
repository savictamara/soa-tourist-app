using System.ComponentModel.DataAnnotations;

namespace BlogService.DTOs;

public class UpdateBlogCommentRequestDto
{
    [Required]
    [MaxLength(2000)]
    public string Text { get; set; } = string.Empty;
}

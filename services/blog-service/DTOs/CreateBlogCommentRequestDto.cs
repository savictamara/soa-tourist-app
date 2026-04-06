using System.ComponentModel.DataAnnotations;

namespace BlogService.DTOs;

public class CreateBlogCommentRequestDto
{
    [Required]
    [MaxLength(2000)]
    public string Text { get; set; } = string.Empty;
}

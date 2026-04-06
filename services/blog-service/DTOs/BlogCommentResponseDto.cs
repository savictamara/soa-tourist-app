namespace BlogService.DTOs;

public class BlogCommentResponseDto
{
    public long Id { get; set; }
    public string AuthorUsername { get; set; } = string.Empty;
    public string AuthorRole { get; set; } = string.Empty;
    public string Text { get; set; } = string.Empty;
    public DateTime CreatedAtUtc { get; set; }
    public DateTime LastModifiedAtUtc { get; set; }
}

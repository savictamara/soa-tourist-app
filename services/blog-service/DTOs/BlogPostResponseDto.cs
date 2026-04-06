namespace BlogService.DTOs;

public class BlogPostResponseDto
{
    public long Id { get; set; }
    public string Title { get; set; } = string.Empty;
    public string DescriptionMarkdown { get; set; } = string.Empty;
    public DateTime CreatedAtUtc { get; set; }
    public string AuthorUsername { get; set; } = string.Empty;
    public List<string> ImageUrls { get; set; } = [];
    public int LikesCount { get; set; }
    public bool IsLikedByCurrentUser { get; set; }
    public List<BlogCommentResponseDto> Comments { get; set; } = [];
}

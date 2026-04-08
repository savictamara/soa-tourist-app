namespace BlogService.Entities;

public class BlogComment
{
    public long Id { get; set; }
    public long BlogPostId { get; set; }
    public BlogPost BlogPost { get; set; } = null!;
    public string AuthorUsername { get; set; } = string.Empty;
    public string AuthorRole { get; set; } = string.Empty;
    public string Text { get; set; } = string.Empty;
    public DateTime CreatedAtUtc { get; set; }
    public DateTime LastModifiedAtUtc { get; set; }
}

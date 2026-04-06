namespace BlogService.Entities;

public class BlogLike
{
    public long Id { get; set; }
    public long BlogPostId { get; set; }
    public BlogPost BlogPost { get; set; } = null!;
    public string Username { get; set; } = string.Empty;
    public DateTime CreatedAtUtc { get; set; }
}

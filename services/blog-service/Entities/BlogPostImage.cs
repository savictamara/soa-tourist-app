namespace BlogService.Entities;

public class BlogPostImage
{
    public long Id { get; set; }
    public long BlogPostId { get; set; }
    public string ImageUrl { get; set; } = string.Empty;
    public BlogPost BlogPost { get; set; } = null!;
}

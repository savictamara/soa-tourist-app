namespace BlogService.Entities;

public class BlogPost
{
    public long Id { get; set; }
    public string Title { get; set; } = string.Empty;
    public string DescriptionMarkdown { get; set; } = string.Empty;
    public DateTime CreatedAtUtc { get; set; }
    public string AuthorUsername { get; set; } = string.Empty;
    public long AuthorId { get; set; }
    public ICollection<BlogPostImage> Images { get; set; } = new List<BlogPostImage>();
    public ICollection<BlogComment> Comments { get; set; } = new List<BlogComment>();
    public ICollection<BlogLike> Likes { get; set; } = new List<BlogLike>();
}

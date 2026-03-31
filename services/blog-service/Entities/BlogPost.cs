namespace BlogService.Entities;

public class BlogPost
{
    public long Id { get; set; }
    public string Title { get; set; } = string.Empty;
    public string DescriptionMarkdown { get; set; } = string.Empty;
    public DateTime CreatedAtUtc { get; set; }
    public string AuthorUsername { get; set; } = string.Empty;
    public ICollection<BlogPostImage> Images { get; set; } = new List<BlogPostImage>();
}

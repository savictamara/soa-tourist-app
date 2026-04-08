namespace BlogService.DTOs;

public class BlogLikeStatusResponseDto
{
    public long BlogPostId { get; set; }
    public int LikesCount { get; set; }
    public bool IsLikedByCurrentUser { get; set; }
}

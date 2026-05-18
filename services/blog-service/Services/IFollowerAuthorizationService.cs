namespace BlogService.Services;

public interface IFollowerAuthorizationService
{
    Task<List<long>> GetFollowedAuthorIds(long userId);
    Task<List<string>> GetFollowedAuthorUsernames(long userId);
    Task<bool> IsFollowing(long followerId, long targetId);
}

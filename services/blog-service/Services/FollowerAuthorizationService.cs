using System.Net.Http.Json;
using System.Text.Json.Serialization;

namespace BlogService.Services;

file record FollowerUserDto(
    [property: JsonPropertyName("userId")] string UserId,
    [property: JsonPropertyName("username")] string Username,
    [property: JsonPropertyName("role")] string Role);

public class FollowerAuthorizationService : IFollowerAuthorizationService
{
    private readonly HttpClient _httpClient;
    private readonly string _baseUrl;

    public FollowerAuthorizationService(HttpClient httpClient, IConfiguration configuration)
    {
        _httpClient = httpClient;
        _baseUrl = (configuration["FollowerService:BaseUrl"] ?? "http://localhost:8086").TrimEnd('/');
    }

    public async Task<List<long>> GetFollowedAuthorIds(long userId)
    {
        try
        {
            var response = await _httpClient.GetAsync($"{_baseUrl}/api/followers/{userId}/followed-author-ids");
            if (!response.IsSuccessStatusCode)
                return [];

            var ids = await response.Content.ReadFromJsonAsync<List<string>>() ?? [];
            return ids
                .Where(id => long.TryParse(id, out _))
                .Select(long.Parse)
                .ToList();
        }
        catch
        {
            return [];
        }
    }

    public async Task<List<string>> GetFollowedAuthorUsernames(long userId)
    {
        try
        {
            var response = await _httpClient.GetAsync($"{_baseUrl}/api/followers/{userId}/following");
            if (!response.IsSuccessStatusCode)
                return [];

            var users = await response.Content.ReadFromJsonAsync<List<FollowerUserDto>>() ?? [];
            return users
                .Where(u => !string.IsNullOrWhiteSpace(u.Username))
                .Select(u => u.Username)
                .ToList();
        }
        catch
        {
            return [];
        }
    }

    public async Task<bool> IsFollowing(long followerId, long targetId)
    {
        try
        {
            var response = await _httpClient.GetAsync($"{_baseUrl}/api/followers/{followerId}/follows/{targetId}");
            if (!response.IsSuccessStatusCode)
                return false;

            return await response.Content.ReadFromJsonAsync<bool>();
        }
        catch
        {
            return false;
        }
    }
}

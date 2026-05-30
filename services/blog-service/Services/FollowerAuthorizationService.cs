using System.Net.Http.Json;
using System.Text.Json;
using System.Text.Json.Serialization;
using Grpc.Core;
using Grpc.Net.Client;

namespace BlogService.Services;

internal record FollowerUserDto(
    [property: JsonPropertyName("userId")] string UserId,
    [property: JsonPropertyName("username")] string Username,
    [property: JsonPropertyName("role")] string Role);

internal record CanCommentRequest(long CommenterId, long AuthorId);
internal record CanCommentResponse(bool Allowed);
internal record GetFollowedAuthorsRequest(long UserId);
internal record GetFollowedAuthorsResponse(List<long> AuthorIds);

public class FollowerAuthorizationService : IFollowerAuthorizationService
{
    private static readonly JsonSerializerOptions JsonOptions = new(JsonSerializerDefaults.Web);
    private static readonly Marshaller<CanCommentRequest> CanCommentRequestMarshaller = JsonMarshaller<CanCommentRequest>();
    private static readonly Marshaller<CanCommentResponse> CanCommentResponseMarshaller = JsonMarshaller<CanCommentResponse>();
    private static readonly Marshaller<GetFollowedAuthorsRequest> FollowedAuthorsRequestMarshaller = JsonMarshaller<GetFollowedAuthorsRequest>();
    private static readonly Marshaller<GetFollowedAuthorsResponse> FollowedAuthorsResponseMarshaller = JsonMarshaller<GetFollowedAuthorsResponse>();
    private static readonly Method<CanCommentRequest, CanCommentResponse> CanCommentMethod = new(
        MethodType.Unary,
        "followerrpc.FollowerRpc",
        "CanComment",
        CanCommentRequestMarshaller,
        CanCommentResponseMarshaller);
    private static readonly Method<GetFollowedAuthorsRequest, GetFollowedAuthorsResponse> GetFollowedAuthorsMethod = new(
        MethodType.Unary,
        "followerrpc.FollowerRpc",
        "GetFollowedAuthors",
        FollowedAuthorsRequestMarshaller,
        FollowedAuthorsResponseMarshaller);

    private readonly HttpClient _httpClient;
    private readonly string _baseUrl;
    private readonly CallInvoker _callInvoker;
    private readonly ILogger<FollowerAuthorizationService> _logger;

    public FollowerAuthorizationService(HttpClient httpClient, IConfiguration configuration, ILogger<FollowerAuthorizationService> logger)
    {
        _httpClient = httpClient;
        _baseUrl = (configuration["FollowerService:BaseUrl"] ?? "http://localhost:8086").TrimEnd('/');
        var grpcUrl = configuration["FollowerService:GrpcUrl"] ?? "http://localhost:9092";
        _callInvoker = GrpcChannel.ForAddress(grpcUrl).CreateCallInvoker();
        _logger = logger;
    }

    public async Task<List<long>> GetFollowedAuthorIds(long userId)
    {
        try
        {
            var response = await _callInvoker.AsyncUnaryCall(
                GetFollowedAuthorsMethod,
                null,
                new CallOptions(deadline: DateTime.UtcNow.AddSeconds(5)),
                new GetFollowedAuthorsRequest(userId));
            _logger.LogInformation("RPC GetFollowedAuthors userId={UserId} count={Count}", userId, response.AuthorIds.Count);
            return response.AuthorIds;
        }
        catch (Exception ex)
        {
            _logger.LogWarning(ex, "RPC GetFollowedAuthors failed for userId={UserId}", userId);
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
            var response = await _callInvoker.AsyncUnaryCall(
                CanCommentMethod,
                null,
                new CallOptions(deadline: DateTime.UtcNow.AddSeconds(5)),
                new CanCommentRequest(followerId, targetId));
            _logger.LogInformation("RPC CanComment commenterId={CommenterId} authorId={AuthorId} allowed={Allowed}", followerId, targetId, response.Allowed);
            return response.Allowed;
        }
        catch (Exception ex)
        {
            _logger.LogWarning(ex, "RPC CanComment failed for commenterId={CommenterId} authorId={AuthorId}", followerId, targetId);
            return false;
        }
    }

    private static Marshaller<T> JsonMarshaller<T>() =>
        new(
            value => JsonSerializer.SerializeToUtf8Bytes(value, JsonOptions),
            bytes => JsonSerializer.Deserialize<T>(bytes, JsonOptions)
                     ?? throw new InvalidOperationException($"Could not deserialize {typeof(T).Name}"));
}

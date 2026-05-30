package com.soa.followers.rpc;

import java.util.List;

public final class FollowerRpcMessages {
    private FollowerRpcMessages() {
    }

    public record CanCommentRequest(long commenterId, long authorId) {
    }

    public record CanCommentResponse(boolean allowed) {
    }

    public record GetFollowedAuthorsRequest(long userId) {
    }

    public record GetFollowedAuthorsResponse(List<Long> authorIds) {
    }

    public record FollowUserRequest(
            String followerId,
            String followingId,
            String followerUsername,
            String followingUsername,
            String followerRole,
            String followingRole) {
    }

    public record FollowResponse(boolean success, String message) {
    }

    public record UnfollowUserRequest(String followerId, String followingId) {
    }

    public record UnfollowResponse(boolean success, String message) {
    }

    public record GetRecommendationsRequest(String userId) {
    }

    public record UserMessage(String userId, String username, String role) {
    }

    public record RecommendationsResponse(List<UserMessage> recommendations) {
    }
}

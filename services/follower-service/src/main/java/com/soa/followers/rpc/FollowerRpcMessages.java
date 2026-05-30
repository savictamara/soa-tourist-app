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
}

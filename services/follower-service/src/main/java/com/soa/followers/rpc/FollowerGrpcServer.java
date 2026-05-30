package com.soa.followers.rpc;

import com.soa.followers.dto.FollowRequest;
import com.soa.followers.model.UserNode;
import com.soa.followers.service.FollowerService;
import io.grpc.BindableService;
import io.grpc.MethodDescriptor;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.ServerServiceDefinition;
import io.grpc.stub.ServerCalls;
import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.springframework.web.server.ResponseStatusException;

import java.io.IOException;
import java.util.List;

@Component
public class FollowerGrpcServer {
    private static final Logger LOGGER = LoggerFactory.getLogger(FollowerGrpcServer.class);
    private final FollowerService followerService;
    private final int port;
    private Server server;

    public FollowerGrpcServer(FollowerService followerService, @Value("${grpc.server.port:9092}") int port) {
        this.followerService = followerService;
        this.port = port;
    }

    @PostConstruct
    public void start() throws IOException {
        server = ServerBuilder.forPort(port)
                .addService(new FollowerRpcBindableService(followerService))
                .build()
                .start();
        LOGGER.info("follower-service gRPC listening on {}", port);
    }

    @PreDestroy
    public void stop() {
        if (server != null) {
            server.shutdown();
        }
    }

    private static class FollowerRpcBindableService implements BindableService {
        private final FollowerService followerService;

        private FollowerRpcBindableService(FollowerService followerService) {
            this.followerService = followerService;
        }

        @Override
        public ServerServiceDefinition bindService() {
            MethodDescriptor<FollowerRpcMessages.CanCommentRequest, FollowerRpcMessages.CanCommentResponse> canComment =
                    MethodDescriptor.<FollowerRpcMessages.CanCommentRequest, FollowerRpcMessages.CanCommentResponse>newBuilder()
                            .setType(MethodDescriptor.MethodType.UNARY)
                            .setFullMethodName(MethodDescriptor.generateFullMethodName("followerrpc.FollowerRpc", "CanComment"))
                            .setRequestMarshaller(new JsonMarshaller<>(FollowerRpcMessages.CanCommentRequest.class))
                            .setResponseMarshaller(new JsonMarshaller<>(FollowerRpcMessages.CanCommentResponse.class))
                            .build();

            MethodDescriptor<FollowerRpcMessages.GetFollowedAuthorsRequest, FollowerRpcMessages.GetFollowedAuthorsResponse> followedAuthors =
                    MethodDescriptor.<FollowerRpcMessages.GetFollowedAuthorsRequest, FollowerRpcMessages.GetFollowedAuthorsResponse>newBuilder()
                            .setType(MethodDescriptor.MethodType.UNARY)
                            .setFullMethodName(MethodDescriptor.generateFullMethodName("followerrpc.FollowerRpc", "GetFollowedAuthors"))
                            .setRequestMarshaller(new JsonMarshaller<>(FollowerRpcMessages.GetFollowedAuthorsRequest.class))
                            .setResponseMarshaller(new JsonMarshaller<>(FollowerRpcMessages.GetFollowedAuthorsResponse.class))
                            .build();

            MethodDescriptor<FollowerRpcMessages.FollowUserRequest, FollowerRpcMessages.FollowResponse> followUser =
                    MethodDescriptor.<FollowerRpcMessages.FollowUserRequest, FollowerRpcMessages.FollowResponse>newBuilder()
                            .setType(MethodDescriptor.MethodType.UNARY)
                            .setFullMethodName(MethodDescriptor.generateFullMethodName("followerrpc.FollowerRpc", "FollowUser"))
                            .setRequestMarshaller(new JsonMarshaller<>(FollowerRpcMessages.FollowUserRequest.class))
                            .setResponseMarshaller(new JsonMarshaller<>(FollowerRpcMessages.FollowResponse.class))
                            .build();

            MethodDescriptor<FollowerRpcMessages.UnfollowUserRequest, FollowerRpcMessages.UnfollowResponse> unfollowUser =
                    MethodDescriptor.<FollowerRpcMessages.UnfollowUserRequest, FollowerRpcMessages.UnfollowResponse>newBuilder()
                            .setType(MethodDescriptor.MethodType.UNARY)
                            .setFullMethodName(MethodDescriptor.generateFullMethodName("followerrpc.FollowerRpc", "UnfollowUser"))
                            .setRequestMarshaller(new JsonMarshaller<>(FollowerRpcMessages.UnfollowUserRequest.class))
                            .setResponseMarshaller(new JsonMarshaller<>(FollowerRpcMessages.UnfollowResponse.class))
                            .build();

            MethodDescriptor<FollowerRpcMessages.GetRecommendationsRequest, FollowerRpcMessages.RecommendationsResponse> recommendations =
                    MethodDescriptor.<FollowerRpcMessages.GetRecommendationsRequest, FollowerRpcMessages.RecommendationsResponse>newBuilder()
                            .setType(MethodDescriptor.MethodType.UNARY)
                            .setFullMethodName(MethodDescriptor.generateFullMethodName("followerrpc.FollowerRpc", "GetRecommendations"))
                            .setRequestMarshaller(new JsonMarshaller<>(FollowerRpcMessages.GetRecommendationsRequest.class))
                            .setResponseMarshaller(new JsonMarshaller<>(FollowerRpcMessages.RecommendationsResponse.class))
                            .build();

            return ServerServiceDefinition.builder("followerrpc.FollowerRpc")
                    .addMethod(canComment, ServerCalls.asyncUnaryCall((request, observer) -> {
                        boolean allowed = request.commenterId() == request.authorId()
                                || followerService.isFollowing(String.valueOf(request.commenterId()), String.valueOf(request.authorId()));
                        LOGGER.info("RPC CanComment commenterId={} authorId={} allowed={}", request.commenterId(), request.authorId(), allowed);
                        observer.onNext(new FollowerRpcMessages.CanCommentResponse(allowed));
                        observer.onCompleted();
                    }))
                    .addMethod(followedAuthors, ServerCalls.asyncUnaryCall((request, observer) -> {
                        List<Long> ids = followerService.getFollowedAuthorIds(String.valueOf(request.userId())).stream()
                                .filter(value -> value != null && value.matches("\\d+"))
                                .map(Long::parseLong)
                                .toList();
                        LOGGER.info("RPC GetFollowedAuthors userId={} count={}", request.userId(), ids.size());
                        observer.onNext(new FollowerRpcMessages.GetFollowedAuthorsResponse(ids));
                        observer.onCompleted();
                    }))
                    .addMethod(followUser, ServerCalls.asyncUnaryCall((request, observer) -> {
                        try {
                            followerService.follow(new FollowRequest(request.followerId(), request.followingId()));
                            LOGGER.info("RPC FollowUser followerId={} followingId={} success=true", request.followerId(), request.followingId());
                            observer.onNext(new FollowerRpcMessages.FollowResponse(true, "followed successfully"));
                        } catch (ResponseStatusException ex) {
                            String message = ex.getReason() != null ? ex.getReason() : "follow failed";
                            LOGGER.info("RPC FollowUser followerId={} followingId={} success=false message={}", request.followerId(), request.followingId(), message);
                            observer.onNext(new FollowerRpcMessages.FollowResponse(false, message));
                        }
                        observer.onCompleted();
                    }))
                    .addMethod(unfollowUser, ServerCalls.asyncUnaryCall((request, observer) -> {
                        try {
                            followerService.unfollow(new FollowRequest(request.followerId(), request.followingId()));
                            LOGGER.info("RPC UnfollowUser followerId={} followingId={} success=true relationshipRemoved=true", request.followerId(), request.followingId());
                            observer.onNext(new FollowerRpcMessages.UnfollowResponse(true, "unfollowed successfully"));
                        } catch (ResponseStatusException ex) {
                            String message = ex.getReason() != null ? ex.getReason() : "unfollow failed";
                            LOGGER.info("RPC UnfollowUser followerId={} followingId={} success=false message={}", request.followerId(), request.followingId(), message);
                            observer.onNext(new FollowerRpcMessages.UnfollowResponse(false, message));
                        }
                        observer.onCompleted();
                    }))
                    .addMethod(recommendations, ServerCalls.asyncUnaryCall((request, observer) -> {
                        List<FollowerRpcMessages.UserMessage> users = followerService.getRecommendations(request.userId()).stream()
                                .map(FollowerGrpcServer.FollowerRpcBindableService::toUserMessage)
                                .toList();
                        LOGGER.info("RPC GetRecommendations userId={} count={}", request.userId(), users.size());
                        observer.onNext(new FollowerRpcMessages.RecommendationsResponse(users));
                        observer.onCompleted();
                    }))
                    .build();
        }

        private static FollowerRpcMessages.UserMessage toUserMessage(UserNode user) {
            return new FollowerRpcMessages.UserMessage(user.getUserId(), user.getUsername(), user.getRole());
        }
    }
}

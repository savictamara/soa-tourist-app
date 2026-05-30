package rpc

import (
	"context"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FollowUserRequest struct {
	FollowerID        string `json:"followerId"`
	FollowingID       string `json:"followingId"`
	FollowerUsername  string `json:"followerUsername"`
	FollowingUsername string `json:"followingUsername"`
	FollowerRole      string `json:"followerRole"`
	FollowingRole     string `json:"followingRole"`
}

type FollowResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type UnfollowUserRequest struct {
	FollowerID  string `json:"followerId"`
	FollowingID string `json:"followingId"`
}

type UnfollowResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type GetRecommendationsRequest struct {
	UserID string `json:"userId"`
}

type UserMessage struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type RecommendationsResponse struct {
	Recommendations []UserMessage `json:"recommendations"`
}

type CanCommentRequest struct {
	CommenterID int64 `json:"commenterId"`
	AuthorID    int64 `json:"authorId"`
}

type CanCommentResponse struct {
	Allowed bool `json:"allowed"`
}

type GetFollowedAuthorsRequest struct {
	UserID int64 `json:"userId"`
}

type GetFollowedAuthorsResponse struct {
	AuthorIDs []int64 `json:"authorIds"`
}

type FollowerClient struct {
	conn *grpc.ClientConn
}

func NewFollowerClient() (*FollowerClient, error) {
	address := strings.TrimSpace(os.Getenv("FOLLOWER_GRPC_ADDR"))
	if address == "" {
		address = "localhost:9092"
	}
	RegisterCodec()
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(JSONCodec{})),
	)
	if err != nil {
		return nil, err
	}
	log.Printf("gateway-rpc-adapter follower gRPC client configured for %s", address)
	return &FollowerClient{conn: conn}, nil
}

func (c *FollowerClient) Close() error {
	return c.conn.Close()
}

func (c *FollowerClient) FollowUser(ctx context.Context, request *FollowUserRequest) (*FollowResponse, error) {
	response := new(FollowResponse)
	err := c.conn.Invoke(ctx, "/followerrpc.FollowerRpc/FollowUser", request, response)
	return response, err
}

func (c *FollowerClient) UnfollowUser(ctx context.Context, request *UnfollowUserRequest) (*UnfollowResponse, error) {
	response := new(UnfollowResponse)
	err := c.conn.Invoke(ctx, "/followerrpc.FollowerRpc/UnfollowUser", request, response)
	return response, err
}

func (c *FollowerClient) GetRecommendations(ctx context.Context, request *GetRecommendationsRequest) (*RecommendationsResponse, error) {
	response := new(RecommendationsResponse)
	err := c.conn.Invoke(ctx, "/followerrpc.FollowerRpc/GetRecommendations", request, response)
	return response, err
}

func (c *FollowerClient) CanComment(ctx context.Context, request *CanCommentRequest) (*CanCommentResponse, error) {
	response := new(CanCommentResponse)
	err := c.conn.Invoke(ctx, "/followerrpc.FollowerRpc/CanComment", request, response)
	return response, err
}

func (c *FollowerClient) GetFollowedAuthors(ctx context.Context, request *GetFollowedAuthorsRequest) (*GetFollowedAuthorsResponse, error) {
	response := new(GetFollowedAuthorsResponse)
	err := c.conn.Invoke(ctx, "/followerrpc.FollowerRpc/GetFollowedAuthors", request, response)
	return response, err
}

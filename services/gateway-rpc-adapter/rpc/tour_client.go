package rpc

import (
	"context"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CompletedKeyPointMessage struct {
	KeyPointID   string `json:"keyPointId"`
	KeyPointName string `json:"keyPointName"`
	ReachedAt    string `json:"reachedAt"`
}

type TourExecutionMessage struct {
	ID                 string                     `json:"id"`
	TourID             string                     `json:"tourId"`
	TourName           string                     `json:"tourName"`
	TouristID          string                     `json:"touristId"`
	Status             string                     `json:"status"`
	StartedAt          string                     `json:"startedAt"`
	CompletedAt        string                     `json:"completedAt,omitempty"`
	AbandonedAt        string                     `json:"abandonedAt,omitempty"`
	LastActivityAt     string                     `json:"lastActivityAt"`
	StartLatitude      float64                    `json:"startLatitude"`
	StartLongitude     float64                    `json:"startLongitude"`
	CurrentLatitude    float64                    `json:"currentLatitude"`
	CurrentLongitude   float64                    `json:"currentLongitude"`
	CompletedKeyPoints []CompletedKeyPointMessage `json:"completedKeyPoints"`
}

type StartTourExecutionRequest struct {
	TourID    string  `json:"tourId"`
	TouristID string  `json:"touristId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type StartTourExecutionResponse struct {
	Success   bool                 `json:"success"`
	Message   string               `json:"message"`
	Execution TourExecutionMessage `json:"execution"`
}

type CheckTourExecutionLocationRequest struct {
	ExecutionID string  `json:"executionId"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type CheckTourExecutionLocationResponse struct {
	Success            bool                      `json:"success"`
	Message            string                    `json:"message"`
	Execution          TourExecutionMessage      `json:"execution"`
	KeyPointReached    bool                      `json:"keyPointReached"`
	ReachedKeyPoint    *CompletedKeyPointMessage `json:"reachedKeyPoint,omitempty"`
	DistanceMeters     float64                   `json:"distanceMeters"`
	LastActivityAt     string                    `json:"lastActivityAt"`
	CompletedCount     int                       `json:"completedCount"`
	TotalKeyPointCount int                       `json:"totalKeyPointCount"`
}

type TourClient struct {
	conn *grpc.ClientConn
}

func NewTourClient() (*TourClient, error) {
	address := strings.TrimSpace(os.Getenv("TOUR_GRPC_ADDR"))
	if address == "" {
		address = "localhost:9094"
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
	log.Printf("gateway-rpc-adapter tour gRPC client configured for %s", address)
	return &TourClient{conn: conn}, nil
}

func (c *TourClient) Close() error {
	return c.conn.Close()
}

func (c *TourClient) StartTourExecution(ctx context.Context, request *StartTourExecutionRequest) (*StartTourExecutionResponse, error) {
	response := new(StartTourExecutionResponse)
	err := c.conn.Invoke(ctx, "/tourgateway.TourGatewayRpc/StartTourExecution", request, response)
	return response, err
}

func (c *TourClient) CheckTourExecutionLocation(ctx context.Context, request *CheckTourExecutionLocationRequest) (*CheckTourExecutionLocationResponse, error) {
	response := new(CheckTourExecutionLocationResponse)
	err := c.conn.Invoke(ctx, "/tourgateway.TourGatewayRpc/CheckTourExecutionLocation", request, response)
	return response, err
}

package gatewayrpc

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"tour-service/models"
	"tour-service/rpc/tourgateway"
	"tour-service/service"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type TourGatewayServer struct {
	service *service.TourExecutionService
}

func StartServer(address string, executionService *service.TourExecutionService) {
	tourgateway.RegisterCodec()
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen on tour gateway gRPC address %s: %v", address, err)
	}
	server := grpc.NewServer(grpc.ForceServerCodec(tourgateway.JSONCodec{}))
	tourgateway.RegisterTourGatewayRpcServer(server, &TourGatewayServer{service: executionService})
	log.Printf("tour-service gateway gRPC listening on %s", address)
	if err = server.Serve(listener); err != nil {
		log.Fatalf("failed to serve tour gateway gRPC: %v", err)
	}
}

func (s *TourGatewayServer) StartTourExecution(ctx context.Context, request *tourgateway.StartTourExecutionRequest) (*tourgateway.StartTourExecutionResponse, error) {
	tourID, err := primitive.ObjectIDFromHex(request.TourID)
	if err != nil {
		return &tourgateway.StartTourExecutionResponse{Success: false, Message: "invalid tourId"}, nil
	}
	execution, err := s.service.StartTour(ctx, tourID, models.StartTourExecutionRequest{
		TouristID: request.TouristID,
		Latitude:  request.Latitude,
		Longitude: request.Longitude,
	})
	if err != nil {
		log.Printf("RPC StartTourExecution tourId=%s touristId=%s success=false message=%s", request.TourID, request.TouristID, err.Error())
		return &tourgateway.StartTourExecutionResponse{
			Success: false,
			Message: errorMessage(err),
		}, nil
	}
	log.Printf("RPC StartTourExecution tourId=%s touristId=%s success=true executionId=%s", request.TourID, request.TouristID, execution.ID.Hex())
	return &tourgateway.StartTourExecutionResponse{
		Success:   true,
		Message:   "tour execution started",
		Execution: executionToMessage(execution),
	}, nil
}

func (s *TourGatewayServer) CheckTourExecutionLocation(ctx context.Context, request *tourgateway.CheckTourExecutionLocationRequest) (*tourgateway.CheckTourExecutionLocationResponse, error) {
	executionID, err := primitive.ObjectIDFromHex(request.ExecutionID)
	if err != nil {
		return &tourgateway.CheckTourExecutionLocationResponse{Success: false, Message: "invalid executionId"}, nil
	}
	response, err := s.service.CheckLocation(ctx, executionID, models.CheckTourExecutionLocationRequest{
		Latitude:  request.Latitude,
		Longitude: request.Longitude,
	})
	if err != nil {
		log.Printf("RPC CheckTourExecutionLocation executionId=%s success=false message=%s", request.ExecutionID, err.Error())
		return &tourgateway.CheckTourExecutionLocationResponse{
			Success: false,
			Message: errorMessage(err),
		}, nil
	}
	log.Printf("RPC CheckTourExecutionLocation executionId=%s success=true reached=%t", request.ExecutionID, response.KeyPointReached)
	return &tourgateway.CheckTourExecutionLocationResponse{
		Success:            true,
		Message:            "location checked",
		Execution:          executionToMessage(response.Execution),
		KeyPointReached:    response.KeyPointReached,
		ReachedKeyPoint:    completedKeyPointPtrToMessage(response.ReachedKeyPoint),
		DistanceMeters:     response.DistanceMeters,
		LastActivityAt:     response.LastActivityAt.UTC().Format(time.RFC3339Nano),
		CompletedCount:     response.CompletedCount,
		TotalKeyPointCount: response.TotalKeyPointCnt,
	}, nil
}

func errorMessage(err error) string {
	if errors.Is(err, service.ErrPurchaseRequired) {
		return "tour must be purchased before it can be started"
	}
	if errors.Is(err, service.ErrExecutionConflict) {
		return "tour execution is already active"
	}
	if errors.Is(err, service.ErrExecutionNotActive) {
		return "tour execution is not active"
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "tour execution or tour not found"
	}
	return err.Error()
}

func executionToMessage(execution models.TourExecution) tourgateway.TourExecutionMessage {
	completed := make([]tourgateway.CompletedKeyPointMessage, 0, len(execution.CompletedKeyPoints))
	for _, keyPoint := range execution.CompletedKeyPoints {
		completed = append(completed, completedKeyPointToMessage(keyPoint))
	}

	message := tourgateway.TourExecutionMessage{
		ID:                 execution.ID.Hex(),
		TourID:             execution.TourID.Hex(),
		TourName:           execution.TourName,
		TouristID:          execution.TouristID,
		Status:             execution.Status,
		StartedAt:          execution.StartedAt.UTC().Format(time.RFC3339Nano),
		LastActivityAt:     execution.LastActivityAt.UTC().Format(time.RFC3339Nano),
		StartLatitude:      execution.StartLatitude,
		StartLongitude:     execution.StartLongitude,
		CurrentLatitude:    execution.CurrentLatitude,
		CurrentLongitude:   execution.CurrentLongitude,
		CompletedKeyPoints: completed,
	}
	if execution.CompletedAt != nil {
		message.CompletedAt = execution.CompletedAt.UTC().Format(time.RFC3339Nano)
	}
	if execution.AbandonedAt != nil {
		message.AbandonedAt = execution.AbandonedAt.UTC().Format(time.RFC3339Nano)
	}
	return message
}

func completedKeyPointPtrToMessage(keyPoint *models.CompletedKeyPoint) *tourgateway.CompletedKeyPointMessage {
	if keyPoint == nil {
		return nil
	}
	message := completedKeyPointToMessage(*keyPoint)
	return &message
}

func completedKeyPointToMessage(keyPoint models.CompletedKeyPoint) tourgateway.CompletedKeyPointMessage {
	return tourgateway.CompletedKeyPointMessage{
		KeyPointID:   keyPoint.KeyPointID.Hex(),
		KeyPointName: keyPoint.KeyPointName,
		ReachedAt:    keyPoint.ReachedAt.UTC().Format(time.RFC3339Nano),
	}
}

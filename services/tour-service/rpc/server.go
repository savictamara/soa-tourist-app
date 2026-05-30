package rpc

import (
	"context"
	"log"
	"net"

	"tour-service/rpc/tourpurchase"
	"tour-service/service"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

type TourPurchaseServer struct {
	tourService *service.TourService
}

func NewTourPurchaseServer(tourService *service.TourService) *TourPurchaseServer {
	return &TourPurchaseServer{tourService: tourService}
}

func (s *TourPurchaseServer) GetTourForPurchase(ctx context.Context, req *tourpurchase.GetTourForPurchaseRequest) (*tourpurchase.GetTourForPurchaseResponse, error) {
	tourID, err := primitive.ObjectIDFromHex(req.TourID)
	if err != nil {
		return &tourpurchase.GetTourForPurchaseResponse{TourID: req.TourID, Status: "invalid"}, nil
	}
	tour, err := s.tourService.GetTourByID(ctx, tourID)
	if err != nil {
		return &tourpurchase.GetTourForPurchaseResponse{TourID: req.TourID, Status: "missing"}, nil
	}
	log.Printf("RPC GetTourForPurchase tourId=%s status=%s price=%.2f archived=%t", req.TourID, tour.Status, tour.Price, tour.ArchivedAt != nil)
	return &tourpurchase.GetTourForPurchaseResponse{
		TourID:   tour.ID.Hex(),
		Name:     tour.Name,
		Price:    tour.Price,
		Status:   tour.Status,
		Archived: tour.ArchivedAt != nil,
	}, nil
}

func (s *TourPurchaseServer) ValidateTourPurchase(ctx context.Context, req *tourpurchase.ValidateTourPurchaseRequest) (*tourpurchase.ValidateTourPurchaseResponse, error) {
	tour, err := s.GetTourForPurchase(ctx, &tourpurchase.GetTourForPurchaseRequest{TourID: req.TourID})
	if err != nil {
		return &tourpurchase.ValidateTourPurchaseResponse{Valid: false, Message: "tour lookup failed"}, nil
	}
	if tour.Status != "published" || tour.Archived {
		log.Printf("RPC ValidateTourPurchase tourId=%s valid=false message=not purchasable", req.TourID)
		return &tourpurchase.ValidateTourPurchaseResponse{Valid: false, Message: "tour is not published or is archived"}, nil
	}
	if tour.Price <= 0 {
		log.Printf("RPC ValidateTourPurchase tourId=%s valid=false message=price missing", req.TourID)
		return &tourpurchase.ValidateTourPurchaseResponse{Valid: false, Message: "tour price must be greater than 0"}, nil
	}
	log.Printf("RPC ValidateTourPurchase tourId=%s valid=true", req.TourID)
	return &tourpurchase.ValidateTourPurchaseResponse{Valid: true, Message: "ok"}, nil
}

func StartGRPCServer(address string, tourService *service.TourService) {
	tourpurchase.RegisterCodec()
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen on gRPC address %s: %v", address, err)
	}
	server := grpc.NewServer(grpc.ForceServerCodec(tourpurchase.JSONCodec{}))
	tourpurchase.RegisterTourPurchaseRPCServer(server, NewTourPurchaseServer(tourService))
	log.Printf("tour-service gRPC listening on %s", address)
	if err = server.Serve(listener); err != nil {
		log.Fatalf("tour-service gRPC failed: %v", err)
	}
}

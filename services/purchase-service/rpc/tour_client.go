package rpc

import (
	"context"
	"log"
	"os"
	"strings"

	"purchase-service/models"
	"purchase-service/rpc/tourpurchase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TourRPCClient struct {
	conn   *grpc.ClientConn
	client tourpurchase.TourPurchaseRPCClient
}

func NewTourRPCClient() (*TourRPCClient, error) {
	address := strings.TrimSpace(os.Getenv("TOUR_SERVICE_GRPC_ADDR"))
	if address == "" {
		address = "localhost:9091"
	}
	tourpurchase.RegisterCodec()
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(tourpurchase.JSONCodec{})),
	)
	if err != nil {
		return nil, err
	}
	log.Printf("purchase-service gRPC client configured for tour-service at %s", address)
	return &TourRPCClient{
		conn:   conn,
		client: tourpurchase.NewTourPurchaseRPCClient(conn),
	}, nil
}

func (c *TourRPCClient) Close() error {
	return c.conn.Close()
}

func (c *TourRPCClient) GetTourForPurchase(ctx context.Context, tourID string) (models.TourSnapshot, error) {
	res, err := c.client.GetTourForPurchase(ctx, &tourpurchase.GetTourForPurchaseRequest{TourID: tourID})
	if err != nil {
		return models.TourSnapshot{}, err
	}
	log.Printf("RPC GetTourForPurchase response tourId=%s status=%s price=%.2f archived=%t", res.TourID, res.Status, res.Price, res.Archived)
	return models.TourSnapshot{
		ID:       res.TourID,
		Name:     res.Name,
		Status:   res.Status,
		Archived: res.Archived,
		Price:    res.Price,
	}, nil
}

func (c *TourRPCClient) ValidateTourPurchase(ctx context.Context, tourID string) (bool, string, error) {
	res, err := c.client.ValidateTourPurchase(ctx, &tourpurchase.ValidateTourPurchaseRequest{TourID: tourID})
	if err != nil {
		return false, "", err
	}
	log.Printf("RPC ValidateTourPurchase response tourId=%s valid=%t message=%s", tourID, res.Valid, res.Message)
	return res.Valid, res.Message, nil
}

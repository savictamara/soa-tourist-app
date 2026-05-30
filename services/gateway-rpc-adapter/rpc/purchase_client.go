package rpc

import (
	"context"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AddToCartRequest struct {
	TouristID string  `json:"touristId"`
	TourID    string  `json:"tourId"`
	TourName  string  `json:"tourName"`
	Price     float64 `json:"price"`
}

type CheckoutRequest struct {
	TouristID string `json:"touristId"`
}

type OrderItemMessage struct {
	TourID   string  `json:"tourId"`
	TourName string  `json:"tourName"`
	Price    float64 `json:"price"`
}

type ShoppingCartMessage struct {
	TouristID  string             `json:"touristId"`
	Items      []OrderItemMessage `json:"items"`
	TotalPrice float64            `json:"totalPrice"`
}

type PurchaseTokenMessage struct {
	TouristID   string `json:"touristId"`
	TourID      string `json:"tourId"`
	PurchasedAt string `json:"purchasedAt"`
	Token       string `json:"token"`
}

type ShoppingCartResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Cart    ShoppingCartMessage `json:"cart"`
}

type CheckoutResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Tokens  []PurchaseTokenMessage `json:"tokens"`
}

type PurchaseClient struct {
	conn *grpc.ClientConn
}

func NewPurchaseClient() (*PurchaseClient, error) {
	address := strings.TrimSpace(os.Getenv("PURCHASE_GRPC_ADDR"))
	if address == "" {
		address = "localhost:9093"
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
	log.Printf("gateway-rpc-adapter purchase gRPC client configured for %s", address)
	return &PurchaseClient{conn: conn}, nil
}

func (c *PurchaseClient) Close() error {
	return c.conn.Close()
}

func (c *PurchaseClient) AddToCart(ctx context.Context, request *AddToCartRequest) (*ShoppingCartResponse, error) {
	response := new(ShoppingCartResponse)
	err := c.conn.Invoke(ctx, "/purchasegateway.PurchaseGatewayRpc/AddToCart", request, response)
	return response, err
}

func (c *PurchaseClient) Checkout(ctx context.Context, request *CheckoutRequest) (*CheckoutResponse, error) {
	response := new(CheckoutResponse)
	err := c.conn.Invoke(ctx, "/purchasegateway.PurchaseGatewayRpc/Checkout", request, response)
	return response, err
}

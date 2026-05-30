package gatewayrpc

import (
	"context"
	"log"
	"net"
	"time"

	"purchase-service/models"
	"purchase-service/rpc/purchasegateway"
	"purchase-service/service"

	"google.golang.org/grpc"
)

type PurchaseGatewayServer struct {
	service *service.PurchaseService
}

func StartServer(address string, purchaseService *service.PurchaseService) {
	purchasegateway.RegisterCodec()
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen on purchase gateway gRPC address %s: %v", address, err)
	}
	server := grpc.NewServer(grpc.ForceServerCodec(purchasegateway.JSONCodec{}))
	purchasegateway.RegisterPurchaseGatewayRpcServer(server, &PurchaseGatewayServer{service: purchaseService})
	log.Printf("purchase-service gateway gRPC listening on %s", address)
	if err = server.Serve(listener); err != nil {
		log.Fatalf("failed to serve purchase gateway gRPC: %v", err)
	}
}

func (s *PurchaseGatewayServer) AddToCart(ctx context.Context, request *purchasegateway.AddToCartRequest) (*purchasegateway.ShoppingCartResponse, error) {
	cart, err := s.service.AddItem(ctx, request.TouristID, models.OrderItem{
		TourID:   request.TourID,
		TourName: request.TourName,
		Price:    request.Price,
	})
	if err != nil {
		log.Printf("RPC AddToCart touristId=%s tourId=%s success=false message=%s", request.TouristID, request.TourID, err.Error())
		return &purchasegateway.ShoppingCartResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	log.Printf("RPC AddToCart touristId=%s tourId=%s success=true totalPrice=%.2f", request.TouristID, request.TourID, cart.TotalPrice)
	return &purchasegateway.ShoppingCartResponse{
		Success: true,
		Message: "item added to cart",
		Cart:    cartToMessage(cart),
	}, nil
}

func (s *PurchaseGatewayServer) Checkout(ctx context.Context, request *purchasegateway.CheckoutRequest) (*purchasegateway.CheckoutResponse, error) {
	tokens, err := s.service.Checkout(ctx, request.TouristID)
	if err != nil {
		log.Printf("RPC Checkout touristId=%s success=false message=%s", request.TouristID, err.Error())
		return &purchasegateway.CheckoutResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	log.Printf("RPC Checkout touristId=%s success=true tokens=%d", request.TouristID, len(tokens))
	return &purchasegateway.CheckoutResponse{
		Success: true,
		Message: "checkout successful",
		Tokens:  tokensToMessages(tokens),
	}, nil
}

func cartToMessage(cart models.ShoppingCart) purchasegateway.ShoppingCartMessage {
	items := make([]purchasegateway.OrderItemMessage, 0, len(cart.Items))
	for _, item := range cart.Items {
		items = append(items, purchasegateway.OrderItemMessage{
			TourID:   item.TourID,
			TourName: item.TourName,
			Price:    item.Price,
		})
	}
	return purchasegateway.ShoppingCartMessage{
		TouristID:  cart.TouristID,
		Items:      items,
		TotalPrice: cart.TotalPrice,
	}
}

func tokensToMessages(tokens []models.TourPurchaseToken) []purchasegateway.PurchaseTokenMessage {
	result := make([]purchasegateway.PurchaseTokenMessage, 0, len(tokens))
	for _, token := range tokens {
		purchasedAt := token.PurchasedAt.UTC().Format(time.RFC3339Nano)
		result = append(result, purchasegateway.PurchaseTokenMessage{
			TouristID:   token.TouristID,
			TourID:      token.TourID,
			PurchasedAt: purchasedAt,
			Token:       token.Token,
		})
	}
	return result
}

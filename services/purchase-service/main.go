package main

import (
	"context"
	"log"
	"time"

	"purchase-service/config"
	"purchase-service/gatewayrpc"
	"purchase-service/handlers"
	"purchase-service/repository"
	"purchase-service/routes"
	"purchase-service/rpc"
	"purchase-service/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	mongoCfg, err := config.ConnectMongo()
	if err != nil {
		log.Fatalf("failed to connect mongo: %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if disconnectErr := mongoCfg.Client.Disconnect(ctx); disconnectErr != nil {
			log.Printf("failed to disconnect mongo: %v", disconnectErr)
		}
	}()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	repo := repository.NewPurchaseRepository(mongoCfg.Carts, mongoCfg.PurchaseTokens)
	tourRPCClient, err := rpc.NewTourRPCClient()
	if err != nil {
		log.Fatalf("failed to configure tour-service gRPC client: %v", err)
	}
	defer tourRPCClient.Close()

	purchaseService := service.NewPurchaseService(repo, tourRPCClient)
	go gatewayrpc.StartServer(":9093", purchaseService)

	purchaseHandler := handlers.NewPurchaseHandler(purchaseService)
	routes.RegisterRoutes(router, purchaseHandler)

	if err = router.Run(":8087"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

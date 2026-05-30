package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"gateway-rpc-adapter/rpc"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	purchaseClient, err := rpc.NewPurchaseClient()
	if err != nil {
		log.Fatalf("failed to configure purchase gRPC client: %v", err)
	}
	defer purchaseClient.Close()

	followerClient, err := rpc.NewFollowerClient()
	if err != nil {
		log.Fatalf("failed to configure follower gRPC client: %v", err)
	}
	defer followerClient.Close()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/api/purchases/cart/:touristId/items", func(c *gin.Context) {
		var request rpc.AddToCartRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		request.TouristID = c.Param("touristId")
		log.Printf("Gateway RPC adapter received AddToCart HTTP request touristId=%s tourId=%s", request.TouristID, request.TourID)
		response, err := purchaseClient.AddToCart(c.Request.Context(), &request)
		if err != nil {
			log.Printf("Gateway RPC adapter AddToCart gRPC failed touristId=%s tourId=%s error=%v", request.TouristID, request.TourID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add item"})
			return
		}
		log.Printf("Gateway RPC adapter called Purchase gRPC AddToCart touristId=%s tourId=%s success=%t", request.TouristID, request.TourID, response.Success)
		if !response.Success {
			c.JSON(errorStatus(response.Message), gin.H{"error": response.Message})
			return
		}
		c.JSON(http.StatusCreated, response.Cart)
	})

	router.POST("/api/purchases/cart/:touristId/checkout", func(c *gin.Context) {
		request := rpc.CheckoutRequest{TouristID: c.Param("touristId")}
		log.Printf("Gateway RPC adapter received Checkout HTTP request touristId=%s", request.TouristID)
		response, err := purchaseClient.Checkout(c.Request.Context(), &request)
		if err != nil {
			log.Printf("Gateway RPC adapter Checkout gRPC failed touristId=%s error=%v", request.TouristID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to checkout"})
			return
		}
		log.Printf("Gateway RPC adapter called Purchase gRPC Checkout touristId=%s success=%t tokens=%d", request.TouristID, response.Success, len(response.Tokens))
		if !response.Success {
			c.JSON(errorStatus(response.Message), gin.H{"error": response.Message})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"tokens": response.Tokens})
	})

	router.POST("/api/followers/follow", func(c *gin.Context) {
		var request rpc.FollowUserRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		log.Printf("Gateway RPC adapter received Follow HTTP request followerId=%s followingId=%s", request.FollowerID, request.FollowingID)
		response, err := followerClient.FollowUser(c.Request.Context(), &request)
		if err != nil {
			log.Printf("Gateway RPC adapter FollowUser gRPC failed followerId=%s followingId=%s error=%v", request.FollowerID, request.FollowingID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to follow user"})
			return
		}
		log.Printf("Gateway RPC adapter called Follower gRPC FollowUser followerId=%s followingId=%s success=%t", request.FollowerID, request.FollowingID, response.Success)
		if !response.Success {
			c.JSON(errorStatus(response.Message), gin.H{"error": response.Message})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": response.Message})
	})

	router.DELETE("/api/followers/follow", func(c *gin.Context) {
		var request rpc.UnfollowUserRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		log.Printf("Gateway RPC adapter received DELETE /api/followers/follow followerId=%s followingId=%s", request.FollowerID, request.FollowingID)
		response, err := followerClient.UnfollowUser(c.Request.Context(), &request)
		if err != nil {
			log.Printf("Gateway RPC adapter UnfollowUser gRPC failed followerId=%s followingId=%s error=%v", request.FollowerID, request.FollowingID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unfollow user"})
			return
		}
		log.Printf("Gateway RPC adapter called Follower gRPC UnfollowUser followerId=%s followingId=%s success=%t", request.FollowerID, request.FollowingID, response.Success)
		if !response.Success {
			c.JSON(errorStatus(response.Message), gin.H{"error": response.Message})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": response.Message})
	})

	router.GET("/api/followers/:userId/recommendations", func(c *gin.Context) {
		request := rpc.GetRecommendationsRequest{UserID: c.Param("userId")}
		log.Printf("Gateway RPC adapter received Recommendations HTTP request userId=%s", request.UserID)
		response, err := followerClient.GetRecommendations(c.Request.Context(), &request)
		if err != nil {
			log.Printf("Gateway RPC adapter GetRecommendations gRPC failed userId=%s error=%v", request.UserID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch recommendations"})
			return
		}
		log.Printf("Gateway RPC adapter called Follower gRPC GetRecommendations userId=%s count=%d", request.UserID, len(response.Recommendations))
		c.JSON(http.StatusOK, response.Recommendations)
	})

	if err = router.Run(":8090"); err != nil {
		log.Fatalf("failed to start gateway-rpc-adapter: %v", err)
	}
}

func errorStatus(message string) int {
	normalized := strings.ToLower(message)
	if strings.Contains(normalized, "already") {
		return http.StatusConflict
	}
	return http.StatusBadRequest
}

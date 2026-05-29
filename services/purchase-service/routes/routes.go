package routes

import (
	"purchase-service/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, handler *handlers.PurchaseHandler) {
	api := router.Group("/api/purchases")
	{
		api.GET("/health", handler.Health)
		api.GET("/cart/:touristId", handler.GetCart)
		api.POST("/cart/:touristId/items", handler.AddItem)
		api.DELETE("/cart/:touristId/items/:tourId", handler.RemoveItem)
		api.POST("/cart/:touristId/checkout", handler.Checkout)
		api.GET("/tokens/:touristId", handler.GetTokens)
		api.GET("/tokens/:touristId/tour/:tourId", handler.IsPurchased)
	}
}

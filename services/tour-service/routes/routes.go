package routes

import (
	"tour-service/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, handler *handlers.TourHandler) {
	api := router.Group("/api/tours")
	{
		api.GET("/health", handler.Health)
		api.POST("", handler.CreateTour)
		api.GET("", handler.GetTours)
		api.GET("/author/:authorId", handler.GetToursByAuthor)
		api.GET("/:tourId", handler.GetTour)
		api.POST("/:tourId/key-points", handler.AddKeyPoint)
		api.GET("/:tourId/key-points", handler.GetKeyPoints)
		api.POST("/:tourId/reviews", handler.AddReview)
		api.GET("/:tourId/reviews", handler.GetReviews)
	}
}

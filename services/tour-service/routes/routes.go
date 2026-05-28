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
		api.GET("/published", handler.GetPublishedTours)
		api.GET("/author/:authorId", handler.GetToursByAuthor)
		api.GET("/:tourId", handler.GetTour)
		api.PUT("/:tourId/durations", handler.UpdateDurations)
		api.POST("/:tourId/publish", handler.PublishTour)
		api.POST("/:tourId/archive", handler.ArchiveTour)
		api.POST("/:tourId/reactivate", handler.ReactivateTour)
		api.POST("/:tourId/key-points", handler.AddKeyPoint)
		api.GET("/:tourId/key-points", handler.GetKeyPoints)
		api.PUT("/:tourId/key-points/:keyPointId", handler.UpdateKeyPoint)
		api.DELETE("/:tourId/key-points/:keyPointId", handler.DeleteKeyPoint)
		api.POST("/:tourId/reviews", handler.AddReview)
		api.GET("/:tourId/reviews", handler.GetReviews)
	}
}

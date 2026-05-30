package routes

import (
	"tour-service/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, handler *handlers.TourHandler, executionHandler *handlers.TourExecutionHandler) {
	api := router.Group("/api/tours")
	{
		api.GET("/health", handler.Health)
		api.POST("", handler.CreateTour)
		api.GET("", handler.GetTours)
		api.GET("/published", handler.GetPublishedTours)
		api.GET("/available", handler.GetAvailableTours)
		api.GET("/author/:authorId", handler.GetToursByAuthor)
		api.GET("/executions/active/:touristId", executionHandler.GetActiveExecution)
		api.POST("/executions/:executionId/check-location", executionHandler.CheckLocation)
		api.POST("/executions/:executionId/complete", executionHandler.Complete)
		api.POST("/executions/:executionId/abandon", executionHandler.Abandon)
		api.GET("/:tourId", handler.GetTour)
		api.GET("/:tourId/executions/latest/:touristId", executionHandler.GetLatestExecution)
		api.POST("/:tourId/executions/start", executionHandler.StartTour)
		api.PUT("/:tourId/durations", handler.UpdateDurations)
		api.PUT("/:tourId/price", handler.UpdatePrice)
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

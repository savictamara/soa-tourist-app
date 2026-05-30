package handlers

import (
	"errors"
	"net/http"

	"tour-service/models"
	"tour-service/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TourExecutionHandler struct {
	executionService *service.TourExecutionService
}

func NewTourExecutionHandler(executionService *service.TourExecutionService) *TourExecutionHandler {
	return &TourExecutionHandler{executionService: executionService}
}

func (h *TourExecutionHandler) StartTour(c *gin.Context) {
	tourID, ok := parseObjectID(c, "tourId")
	if !ok {
		return
	}
	var req models.StartTourExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	execution, err := h.executionService.StartTour(c.Request.Context(), tourID, req)
	if err != nil {
		writeExecutionError(c, err, "failed to start tour")
		return
	}
	c.JSON(http.StatusCreated, execution)
}

func (h *TourExecutionHandler) GetActiveExecution(c *gin.Context) {
	execution, err := h.executionService.GetActiveByTouristID(c.Request.Context(), c.Param("touristId"))
	if err != nil {
		writeExecutionError(c, err, "failed to fetch active tour execution")
		return
	}
	c.JSON(http.StatusOK, execution)
}

func (h *TourExecutionHandler) GetLatestExecution(c *gin.Context) {
	tourID, ok := parseObjectID(c, "tourId")
	if !ok {
		return
	}
	execution, err := h.executionService.GetLatestByTouristAndTour(c.Request.Context(), c.Param("touristId"), tourID)
	if err != nil {
		writeExecutionError(c, err, "failed to fetch latest tour execution")
		return
	}
	c.JSON(http.StatusOK, execution)
}

func (h *TourExecutionHandler) CheckLocation(c *gin.Context) {
	executionID, ok := parseExecutionID(c)
	if !ok {
		return
	}
	var req models.CheckTourExecutionLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := h.executionService.CheckLocation(c.Request.Context(), executionID, req)
	if err != nil {
		writeExecutionError(c, err, "failed to check tour location")
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *TourExecutionHandler) Complete(c *gin.Context) {
	executionID, ok := parseExecutionID(c)
	if !ok {
		return
	}
	execution, err := h.executionService.Complete(c.Request.Context(), executionID)
	if err != nil {
		writeExecutionError(c, err, "failed to complete tour execution")
		return
	}
	c.JSON(http.StatusOK, execution)
}

func (h *TourExecutionHandler) Abandon(c *gin.Context) {
	executionID, ok := parseExecutionID(c)
	if !ok {
		return
	}
	execution, err := h.executionService.Abandon(c.Request.Context(), executionID)
	if err != nil {
		writeExecutionError(c, err, "failed to abandon tour execution")
		return
	}
	c.JSON(http.StatusOK, execution)
}

func parseExecutionID(c *gin.Context) (primitive.ObjectID, bool) {
	id, err := primitive.ObjectIDFromHex(c.Param("executionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid executionId"})
		return primitive.NilObjectID, false
	}
	return id, true
}

func writeExecutionError(c *gin.Context, err error, fallback string) {
	if errors.Is(err, service.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationMessage(err)})
		return
	}
	if errors.Is(err, service.ErrPurchaseRequired) {
		c.JSON(http.StatusForbidden, gin.H{"error": "tour must be purchased before it can be started"})
		return
	}
	if errors.Is(err, service.ErrExecutionConflict) {
		c.JSON(http.StatusConflict, gin.H{"error": "tour execution is already active"})
		return
	}
	if errors.Is(err, service.ErrExecutionNotActive) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tour execution is not active"})
		return
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		c.JSON(http.StatusNotFound, gin.H{"error": "tour execution or tour not found"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": fallback})
}

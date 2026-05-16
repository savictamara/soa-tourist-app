package handlers

import (
	"errors"
	"log"
	"net/http"

	"tour-service/models"
	"tour-service/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TourHandler struct {
	tourService *service.TourService
}

func NewTourHandler(tourService *service.TourService) *TourHandler {
	return &TourHandler{tourService: tourService}
}

func (h *TourHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *TourHandler) CreateTour(c *gin.Context) {
	var req models.CreateTourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	log.Printf("Create tour payload: authorId=%s name=%s difficulty=%s tagsCount=%d", req.AuthorID, req.Name, req.Difficulty, len(req.Tags))

	tour, err := h.tourService.CreateTour(c.Request.Context(), req)
	if err != nil {
		log.Printf("Create tour error: %v", err)
		if errors.Is(err, service.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "authorId, name, description and difficulty are required"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tour"})
		return
	}
	log.Printf("Inserted tour id: %s", tour.ID.Hex())

	c.JSON(http.StatusCreated, tour)
}

func (h *TourHandler) GetTour(c *gin.Context) {
	tourID, ok := parseObjectID(c, "tourId")
	if !ok {
		return
	}

	tour, err := h.tourService.GetTourByID(c.Request.Context(), tourID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "tour not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tour"})
		return
	}

	c.JSON(http.StatusOK, tour)
}

func (h *TourHandler) GetToursByAuthor(c *gin.Context) {
	authorID := c.Param("authorId")
	tours, err := h.tourService.GetToursByAuthorID(c.Request.Context(), authorID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "authorId is required"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tours"})
		return
	}

	c.JSON(http.StatusOK, tours)
}

func (h *TourHandler) AddKeyPoint(c *gin.Context) {
	tourID, ok := parseObjectID(c, "tourId")
	if !ok {
		return
	}

	var req models.AddKeyPointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	log.Printf("Add key point payload: name=%s descriptionLen=%d hasImage=%t lat=%f lng=%f", req.Name, len(req.Description), req.ImageURL != "", req.Latitude, req.Longitude)
	log.Printf("Tour id: %s", tourID.Hex())

	keyPoint, err := h.tourService.AddKeyPoint(c.Request.Context(), tourID, req)
	if err != nil {
		log.Printf("Add key point error: %v", err)
		if errors.Is(err, service.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name, description, imageUrl, latitude and longitude are required and must be valid"})
			return
		}
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "tour not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add key point"})
		return
	}

	c.JSON(http.StatusCreated, keyPoint)
}

func (h *TourHandler) GetKeyPoints(c *gin.Context) {
	tourID, ok := parseObjectID(c, "tourId")
	if !ok {
		return
	}

	keyPoints, err := h.tourService.GetKeyPoints(c.Request.Context(), tourID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "tour not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch key points"})
		return
	}

	c.JSON(http.StatusOK, keyPoints)
}

func parseObjectID(c *gin.Context, paramName string) (primitive.ObjectID, bool) {
	idHex := c.Param(paramName)
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tourId"})
		return primitive.NilObjectID, false
	}
	return id, true
}

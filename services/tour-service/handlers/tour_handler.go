package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

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

func (h *TourHandler) GetTours(c *gin.Context) {
	tours, err := h.tourService.GetAllTours(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tours"})
		return
	}
	c.JSON(http.StatusOK, tours)
}

func (h *TourHandler) GetPublishedTours(c *gin.Context) {
	tours, err := h.tourService.GetPublishedTours(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch published tours"})
		return
	}
	c.JSON(http.StatusOK, tours)
}

func (h *TourHandler) UpdateDurations(c *gin.Context) {
	tourID, ok := parseObjectID(c, "tourId")
	if !ok {
		return
	}

	var req models.UpdateDurationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	tour, err := h.tourService.UpdateDurations(c.Request.Context(), tourID, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationMessage(err)})
			return
		}
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "tour not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update durations"})
		return
	}

	c.JSON(http.StatusOK, tour)
}

func (h *TourHandler) UpdatePrice(c *gin.Context) {
	tourID, ok := parseObjectID(c, "tourId")
	if !ok {
		return
	}

	var req models.UpdatePriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	tour, err := h.tourService.UpdatePrice(c.Request.Context(), tourID, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationMessage(err)})
			return
		}
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "tour not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update price"})
		return
	}

	c.JSON(http.StatusOK, tour)
}

func (h *TourHandler) PublishTour(c *gin.Context) {
	tourID, ok := parseObjectID(c, "tourId")
	if !ok {
		return
	}

	tour, err := h.tourService.PublishTour(c.Request.Context(), tourID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationMessage(err)})
			return
		}
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "tour not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish tour"})
		return
	}

	c.JSON(http.StatusOK, tour)
}

func (h *TourHandler) ArchiveTour(c *gin.Context) {
	tourID, ok := parseObjectID(c, "tourId")
	if !ok {
		return
	}

	tour, err := h.tourService.ArchiveTour(c.Request.Context(), tourID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationMessage(err)})
			return
		}
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "tour not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to archive tour"})
		return
	}

	c.JSON(http.StatusOK, tour)
}

func (h *TourHandler) ReactivateTour(c *gin.Context) {
	tourID, ok := parseObjectID(c, "tourId")
	if !ok {
		return
	}

	tour, err := h.tourService.ReactivateTour(c.Request.Context(), tourID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationMessage(err)})
			return
		}
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "tour not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reactivate tour"})
		return
	}

	c.JSON(http.StatusOK, tour)
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

func (h *TourHandler) UpdateKeyPoint(c *gin.Context) {
	tourID, ok := parseObjectID(c, "tourId")
	if !ok {
		return
	}
	keyPointID, ok := parseObjectID(c, "keyPointId")
	if !ok {
		return
	}

	var req models.UpdateKeyPointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	keyPoint, err := h.tourService.UpdateKeyPoint(c.Request.Context(), tourID, keyPointID, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name, description, imageUrl, latitude and longitude are required and must be valid"})
			return
		}
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "tour or key point not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update key point"})
		return
	}

	c.JSON(http.StatusOK, keyPoint)
}

func (h *TourHandler) DeleteKeyPoint(c *gin.Context) {
	tourID, ok := parseObjectID(c, "tourId")
	if !ok {
		return
	}
	keyPointID, ok := parseObjectID(c, "keyPointId")
	if !ok {
		return
	}

	if err := h.tourService.DeleteKeyPoint(c.Request.Context(), tourID, keyPointID); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "tour or key point not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete key point"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *TourHandler) AddReview(c *gin.Context) {
	tourID, ok := parseObjectID(c, "tourId")
	if !ok {
		return
	}

	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	review, err := h.tourService.AddReview(c.Request.Context(), tourID, req)
	if err != nil {
		if errors.Is(err, service.ErrFutureDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Visited date cannot be in the future."})
			return
		}
		if errors.Is(err, service.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "rating must be 1-5, comment, touristId, touristUsername and visitedDate are required"})
			return
		}
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "tour not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add review"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func (h *TourHandler) GetReviews(c *gin.Context) {
	tourID, ok := parseObjectID(c, "tourId")
	if !ok {
		return
	}

	reviews, err := h.tourService.GetReviews(c.Request.Context(), tourID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "tour not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
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

func validationMessage(err error) string {
	return strings.TrimPrefix(err.Error(), service.ErrInvalidInput.Error()+": ")
}

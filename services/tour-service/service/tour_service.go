package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"tour-service/models"
	"tour-service/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrInvalidInput = errors.New("invalid input")

type TourService struct {
	repo *repository.TourRepository
}

func NewTourService(repo *repository.TourRepository) *TourService {
	return &TourService{repo: repo}
}

func (s *TourService) CreateTour(ctx context.Context, req models.CreateTourRequest) (models.Tour, error) {
	if strings.TrimSpace(req.AuthorID) == "" ||
		strings.TrimSpace(req.Name) == "" ||
		strings.TrimSpace(req.Description) == "" ||
		strings.TrimSpace(req.Difficulty) == "" {
		return models.Tour{}, ErrInvalidInput
	}
	difficulty := strings.ToLower(strings.TrimSpace(req.Difficulty))
	if difficulty != "easy" && difficulty != "medium" && difficulty != "hard" {
		return models.Tour{}, ErrInvalidInput
	}

	tags := make([]string, 0, len(req.Tags))
	for _, t := range req.Tags {
		trimmed := strings.TrimSpace(t)
		if trimmed != "" {
			tags = append(tags, trimmed)
		}
	}

	now := time.Now().UTC()
	tour := models.Tour{
		AuthorID:    strings.TrimSpace(req.AuthorID),
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Difficulty:  difficulty,
		Tags:        tags,
		Status:      "draft",
		Price:       0,
		KeyPoints:   []models.KeyPoint{},
		Reviews:     []models.Review{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.repo.Create(ctx, tour)
}

func (s *TourService) GetTourByID(ctx context.Context, tourID primitive.ObjectID) (models.Tour, error) {
	return s.repo.GetByID(ctx, tourID)
}

func (s *TourService) GetToursByAuthorID(ctx context.Context, authorID string) ([]models.Tour, error) {
	if strings.TrimSpace(authorID) == "" {
		return []models.Tour{}, ErrInvalidInput
	}

	tours, err := s.repo.GetByAuthorID(ctx, strings.TrimSpace(authorID))
	if err != nil {
		return []models.Tour{}, err
	}
	if tours == nil {
		return []models.Tour{}, nil
	}
	return tours, nil
}

func (s *TourService) AddKeyPoint(ctx context.Context, tourID primitive.ObjectID, req models.AddKeyPointRequest) (models.KeyPoint, error) {
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Description) == "" || strings.TrimSpace(req.ImageURL) == "" {
		return models.KeyPoint{}, ErrInvalidInput
	}
	if req.Latitude < -90 || req.Latitude > 90 || req.Longitude < -180 || req.Longitude > 180 {
		return models.KeyPoint{}, ErrInvalidInput
	}

	tour, err := s.repo.GetByID(ctx, tourID)
	if err != nil {
		return models.KeyPoint{}, err
	}

	now := time.Now().UTC()
	keyPoint := models.KeyPoint{
		ID:          primitive.NewObjectID(),
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		ImageURL:    strings.TrimSpace(req.ImageURL),
		Order:       len(tour.KeyPoints) + 1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	updatedTour, err := s.repo.AddKeyPoint(ctx, tourID, keyPoint)
	if err != nil {
		return models.KeyPoint{}, err
	}
	if len(updatedTour.KeyPoints) == 0 {
		return models.KeyPoint{}, mongo.ErrNoDocuments
	}

	return updatedTour.KeyPoints[len(updatedTour.KeyPoints)-1], nil
}

func (s *TourService) GetKeyPoints(ctx context.Context, tourID primitive.ObjectID) ([]models.KeyPoint, error) {
	tour, err := s.repo.GetByID(ctx, tourID)
	if err != nil {
		return []models.KeyPoint{}, err
	}
	if tour.KeyPoints == nil {
		return []models.KeyPoint{}, nil
	}
	return tour.KeyPoints, nil
}

func (s *TourService) GetAllTours(ctx context.Context) ([]models.Tour, error) {
	tours, err := s.repo.GetAll(ctx)
	if err != nil {
		return []models.Tour{}, err
	}
	if tours == nil {
		return []models.Tour{}, nil
	}
	return tours, nil
}

func (s *TourService) AddReview(ctx context.Context, tourID primitive.ObjectID, req models.CreateReviewRequest) (models.Review, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return models.Review{}, fmt.Errorf("%w: rating must be between 1 and 5", ErrInvalidInput)
	}
	if strings.TrimSpace(req.Comment) == "" ||
		strings.TrimSpace(req.TouristID) == "" ||
		strings.TrimSpace(req.TouristUsername) == "" ||
		strings.TrimSpace(req.VisitedDate) == "" {
		return models.Review{}, ErrInvalidInput
	}

	now := time.Now().UTC()
	review := models.Review{
		ID:              primitive.NewObjectID(),
		Rating:          req.Rating,
		Comment:         strings.TrimSpace(req.Comment),
		TouristID:       strings.TrimSpace(req.TouristID),
		TouristUsername: strings.TrimSpace(req.TouristUsername),
		VisitedDate:     strings.TrimSpace(req.VisitedDate),
		CommentDate:     now,
		Images:          make([]string, 0, len(req.Images)),
	}
	for _, image := range req.Images {
		trimmed := strings.TrimSpace(image)
		if trimmed != "" {
			review.Images = append(review.Images, trimmed)
		}
	}

	if err := s.repo.AddReview(ctx, tourID, review, now); err != nil {
		return models.Review{}, err
	}

	return review, nil
}

func (s *TourService) GetReviews(ctx context.Context, tourID primitive.ObjectID) ([]models.Review, error) {
	tour, err := s.repo.GetByID(ctx, tourID)
	if err != nil {
		return []models.Review{}, err
	}
	if tour.Reviews == nil {
		return []models.Review{}, nil
	}
	return tour.Reviews, nil
}

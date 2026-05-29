package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"tour-service/models"
	"tour-service/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrInvalidInput = errors.New("invalid input")
var ErrFutureDate = errors.New("visited date cannot be in the future")

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
		LengthKm:    0,
		Durations:   []models.TourDuration{},
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
	updatedTour, err = s.repo.UpdateLength(ctx, tourID, calculateLengthKm(updatedTour.KeyPoints), now)
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

func (s *TourService) GetPublishedTours(ctx context.Context) ([]models.Tour, error) {
	tours, err := s.repo.GetPublished(ctx)
	if err != nil {
		return []models.Tour{}, err
	}
	if tours == nil {
		return []models.Tour{}, nil
	}
	for i := range tours {
		if len(tours[i].KeyPoints) > 1 {
			tours[i].KeyPoints = tours[i].KeyPoints[:1]
		}
	}
	return tours, nil
}

func (s *TourService) UpdateKeyPoint(ctx context.Context, tourID primitive.ObjectID, keyPointID primitive.ObjectID, req models.UpdateKeyPointRequest) (models.KeyPoint, error) {
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Description) == "" || strings.TrimSpace(req.ImageURL) == "" {
		return models.KeyPoint{}, ErrInvalidInput
	}
	if req.Latitude < -90 || req.Latitude > 90 || req.Longitude < -180 || req.Longitude > 180 {
		return models.KeyPoint{}, ErrInvalidInput
	}

	now := time.Now().UTC()
	tour, err := s.repo.UpdateKeyPoint(ctx, tourID, keyPointID, req, now)
	if err != nil {
		return models.KeyPoint{}, err
	}
	tour, err = s.repo.UpdateLength(ctx, tourID, calculateLengthKm(tour.KeyPoints), now)
	if err != nil {
		return models.KeyPoint{}, err
	}

	for _, kp := range tour.KeyPoints {
		if kp.ID == keyPointID {
			return kp, nil
		}
	}
	return models.KeyPoint{}, mongo.ErrNoDocuments
}

func (s *TourService) DeleteKeyPoint(ctx context.Context, tourID primitive.ObjectID, keyPointID primitive.ObjectID) error {
	now := time.Now().UTC()
	tour, err := s.repo.DeleteKeyPoint(ctx, tourID, keyPointID, now)
	if err != nil {
		return err
	}
	_, err = s.repo.UpdateLength(ctx, tourID, calculateLengthKm(tour.KeyPoints), now)
	return err
}

func (s *TourService) UpdateDurations(ctx context.Context, tourID primitive.ObjectID, req models.UpdateDurationsRequest) (models.Tour, error) {
	durations, err := validateDurations(req.Durations)
	if err != nil {
		return models.Tour{}, err
	}
	return s.repo.UpdateDurations(ctx, tourID, durations, time.Now().UTC())
}

func (s *TourService) UpdatePrice(ctx context.Context, tourID primitive.ObjectID, req models.UpdatePriceRequest) (models.Tour, error) {
	if req.Price < 0 {
		return models.Tour{}, fmt.Errorf("%w: price must be greater than or equal to 0", ErrInvalidInput)
	}
	return s.repo.UpdatePrice(ctx, tourID, req.Price, time.Now().UTC())
}

func (s *TourService) PublishTour(ctx context.Context, tourID primitive.ObjectID) (models.Tour, error) {
	tour, err := s.repo.GetByID(ctx, tourID)
	if err != nil {
		return models.Tour{}, err
	}
	if strings.TrimSpace(tour.Name) == "" {
		return models.Tour{}, fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	if strings.TrimSpace(tour.Description) == "" {
		return models.Tour{}, fmt.Errorf("%w: description is required", ErrInvalidInput)
	}
	if strings.TrimSpace(tour.Difficulty) == "" {
		return models.Tour{}, fmt.Errorf("%w: difficulty is required", ErrInvalidInput)
	}
	if len(tour.Tags) == 0 {
		return models.Tour{}, fmt.Errorf("%w: at least one tag is required", ErrInvalidInput)
	}
	if tour.Price <= 0 {
		return models.Tour{}, fmt.Errorf("%w: price must be greater than 0 before publishing", ErrInvalidInput)
	}
	if len(tour.KeyPoints) < 2 {
		return models.Tour{}, fmt.Errorf("%w: at least two key points are required", ErrInvalidInput)
	}
	durations, err := validateDurations(tour.Durations)
	if err != nil {
		return models.Tour{}, err
	}

	now := time.Now().UTC()
	update := bson.M{
		"$set": bson.M{
			"status":      "published",
			"publishedAt": now,
			"durations":   durations,
			"lengthKm":    calculateLengthKm(tour.KeyPoints),
			"updatedAt":   now,
		},
		"$unset": bson.M{"archivedAt": ""},
	}
	return s.repo.UpdateLifecycle(ctx, tourID, update)
}

func (s *TourService) ArchiveTour(ctx context.Context, tourID primitive.ObjectID) (models.Tour, error) {
	tour, err := s.repo.GetByID(ctx, tourID)
	if err != nil {
		return models.Tour{}, err
	}
	if tour.Status != "published" {
		return models.Tour{}, fmt.Errorf("%w: only published tours can be archived", ErrInvalidInput)
	}
	now := time.Now().UTC()
	update := bson.M{
		"$set": bson.M{
			"status":     "archived",
			"archivedAt": now,
			"updatedAt":  now,
		},
	}
	return s.repo.UpdateLifecycle(ctx, tourID, update)
}

func (s *TourService) ReactivateTour(ctx context.Context, tourID primitive.ObjectID) (models.Tour, error) {
	tour, err := s.repo.GetByID(ctx, tourID)
	if err != nil {
		return models.Tour{}, err
	}
	if tour.Status != "archived" {
		return models.Tour{}, fmt.Errorf("%w: only archived tours can be reactivated", ErrInvalidInput)
	}
	now := time.Now().UTC()
	update := bson.M{
		"$set": bson.M{
			"status":        "published",
			"reactivatedAt": now,
			"updatedAt":     now,
		},
		"$unset": bson.M{"archivedAt": ""},
	}
	return s.repo.UpdateLifecycle(ctx, tourID, update)
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

	visitedDate := strings.TrimSpace(req.VisitedDate)
	parsedDate, parseErr := time.Parse("2006-01-02", visitedDate)
	if parseErr != nil {
		return models.Review{}, fmt.Errorf("%w: visitedDate must be in YYYY-MM-DD format", ErrInvalidInput)
	}
	today := time.Now().UTC().Truncate(24 * time.Hour)
	if parsedDate.After(today) {
		return models.Review{}, ErrFutureDate
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

func validateDurations(input []models.TourDuration) ([]models.TourDuration, error) {
	if len(input) == 0 {
		return []models.TourDuration{}, fmt.Errorf("%w: at least one duration is required", ErrInvalidInput)
	}
	seen := map[string]bool{}
	durations := make([]models.TourDuration, 0, len(input))
	for _, duration := range input {
		transportType := strings.ToLower(strings.TrimSpace(duration.TransportType))
		if transportType != "walking" && transportType != "bicycle" && transportType != "car" {
			return []models.TourDuration{}, fmt.Errorf("%w: invalid duration transport type", ErrInvalidInput)
		}
		if seen[transportType] {
			return []models.TourDuration{}, fmt.Errorf("%w: duplicate duration transport type", ErrInvalidInput)
		}
		if duration.Minutes <= 0 {
			return []models.TourDuration{}, fmt.Errorf("%w: duration minutes must be greater than 0", ErrInvalidInput)
		}
		seen[transportType] = true
		durations = append(durations, models.TourDuration{
			TransportType: transportType,
			Minutes:       duration.Minutes,
		})
	}
	return durations, nil
}

func calculateLengthKm(keyPoints []models.KeyPoint) float64 {
	if len(keyPoints) < 2 {
		return 0
	}
	total := 0.0
	for i := 1; i < len(keyPoints); i++ {
		total += haversineKm(
			keyPoints[i-1].Latitude,
			keyPoints[i-1].Longitude,
			keyPoints[i].Latitude,
			keyPoints[i].Longitude,
		)
	}
	return math.Round(total*100) / 100
}

func haversineKm(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)
	rLat1 := degreesToRadians(lat1)
	rLat2 := degreesToRadians(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(rLat1)*math.Cos(rLat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

func degreesToRadians(value float64) float64 {
	return value * math.Pi / 180
}

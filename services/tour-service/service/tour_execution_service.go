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

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrPurchaseRequired = errors.New("tour purchase is required")
var ErrExecutionConflict = errors.New("tour execution already active")
var ErrExecutionNotActive = errors.New("tour execution is not active")

const keyPointReachThresholdKm = 0.3

type TourExecutionService struct {
	tourRepo       *repository.TourRepository
	executionRepo  *repository.TourExecutionRepository
	purchaseClient *PurchaseClient
}

func NewTourExecutionService(tourRepo *repository.TourRepository, executionRepo *repository.TourExecutionRepository, purchaseClient *PurchaseClient) *TourExecutionService {
	return &TourExecutionService{
		tourRepo:       tourRepo,
		executionRepo:  executionRepo,
		purchaseClient: purchaseClient,
	}
}

func (s *TourExecutionService) StartTour(ctx context.Context, tourID primitive.ObjectID, req models.StartTourExecutionRequest) (models.TourExecution, error) {
	touristID := strings.TrimSpace(req.TouristID)
	if touristID == "" || !validCoordinates(req.Latitude, req.Longitude) {
		return models.TourExecution{}, ErrInvalidInput
	}
	tour, err := s.tourRepo.GetByID(ctx, tourID)
	if err != nil {
		return models.TourExecution{}, err
	}
	if tour.Status != "published" && tour.Status != "archived" {
		return models.TourExecution{}, fmt.Errorf("%w: only published or archived tours can be started", ErrInvalidInput)
	}
	purchased, err := s.purchaseClient.IsPurchased(ctx, touristID, tourID.Hex())
	if err != nil {
		return models.TourExecution{}, err
	}
	if !purchased {
		return models.TourExecution{}, ErrPurchaseRequired
	}
	if _, err = s.executionRepo.GetActiveByTouristAndTour(ctx, touristID, tourID); err == nil {
		return models.TourExecution{}, ErrExecutionConflict
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return models.TourExecution{}, err
	}

	now := time.Now().UTC()
	execution := models.TourExecution{
		TourID:             tour.ID,
		TourName:           tour.Name,
		TouristID:          touristID,
		Status:             "active",
		StartedAt:          now,
		LastActivityAt:     now,
		StartLatitude:      req.Latitude,
		StartLongitude:     req.Longitude,
		CurrentLatitude:    req.Latitude,
		CurrentLongitude:   req.Longitude,
		CompletedKeyPoints: []models.CompletedKeyPoint{},
	}
	return s.executionRepo.Create(ctx, execution)
}

func (s *TourExecutionService) GetActiveByTouristID(ctx context.Context, touristID string) (models.TourExecution, error) {
	touristID = strings.TrimSpace(touristID)
	if touristID == "" {
		return models.TourExecution{}, ErrInvalidInput
	}
	return s.executionRepo.GetActiveByTouristID(ctx, touristID)
}

func (s *TourExecutionService) GetLatestByTouristAndTour(ctx context.Context, touristID string, tourID primitive.ObjectID) (models.TourExecution, error) {
	touristID = strings.TrimSpace(touristID)
	if touristID == "" {
		return models.TourExecution{}, ErrInvalidInput
	}
	return s.executionRepo.GetLatestByTouristAndTour(ctx, touristID, tourID)
}

func (s *TourExecutionService) CheckLocation(ctx context.Context, executionID primitive.ObjectID, req models.CheckTourExecutionLocationRequest) (models.CheckTourExecutionLocationResponse, error) {
	if !validCoordinates(req.Latitude, req.Longitude) {
		return models.CheckTourExecutionLocationResponse{}, ErrInvalidInput
	}
	execution, err := s.executionRepo.GetByID(ctx, executionID)
	if err != nil {
		return models.CheckTourExecutionLocationResponse{}, err
	}
	if execution.Status != "active" {
		return models.CheckTourExecutionLocationResponse{}, ErrExecutionNotActive
	}
	tour, err := s.tourRepo.GetByID(ctx, execution.TourID)
	if err != nil {
		return models.CheckTourExecutionLocationResponse{}, err
	}

	now := time.Now().UTC()
	execution.CurrentLatitude = req.Latitude
	execution.CurrentLongitude = req.Longitude
	execution.LastActivityAt = now

	completedIDs := map[primitive.ObjectID]bool{}
	for _, completed := range execution.CompletedKeyPoints {
		completedIDs[completed.KeyPointID] = true
	}

	var reached *models.CompletedKeyPoint
	closestDistanceKm := math.MaxFloat64
	for _, keyPoint := range tour.KeyPoints {
		if completedIDs[keyPoint.ID] {
			continue
		}
		distanceKm := haversineKm(req.Latitude, req.Longitude, keyPoint.Latitude, keyPoint.Longitude)
		if distanceKm <= keyPointReachThresholdKm && distanceKm < closestDistanceKm {
			completed := models.CompletedKeyPoint{
				KeyPointID:   keyPoint.ID,
				KeyPointName: keyPoint.Name,
				ReachedAt:    now,
			}
			reached = &completed
			closestDistanceKm = distanceKm
		}
	}

	if reached != nil {
		execution.CompletedKeyPoints = append(execution.CompletedKeyPoints, *reached)
	}
	execution, err = s.executionRepo.Save(ctx, execution)
	if err != nil {
		return models.CheckTourExecutionLocationResponse{}, err
	}

	distanceMeters := 0.0
	if closestDistanceKm != math.MaxFloat64 {
		distanceMeters = math.Round(closestDistanceKm * 1000)
	}
	return models.CheckTourExecutionLocationResponse{
		Execution:        execution,
		KeyPointReached:  reached != nil,
		ReachedKeyPoint:  reached,
		DistanceMeters:   distanceMeters,
		LastActivityAt:   execution.LastActivityAt,
		CompletedCount:   len(execution.CompletedKeyPoints),
		TotalKeyPointCnt: len(tour.KeyPoints),
	}, nil
}

func (s *TourExecutionService) Complete(ctx context.Context, executionID primitive.ObjectID) (models.TourExecution, error) {
	return s.finish(ctx, executionID, "completed")
}

func (s *TourExecutionService) Abandon(ctx context.Context, executionID primitive.ObjectID) (models.TourExecution, error) {
	return s.finish(ctx, executionID, "abandoned")
}

func (s *TourExecutionService) finish(ctx context.Context, executionID primitive.ObjectID, status string) (models.TourExecution, error) {
	execution, err := s.executionRepo.GetByID(ctx, executionID)
	if err != nil {
		return models.TourExecution{}, err
	}
	if execution.Status != "active" {
		return models.TourExecution{}, ErrExecutionNotActive
	}
	now := time.Now().UTC()
	execution.Status = status
	execution.LastActivityAt = now
	if status == "completed" {
		execution.CompletedAt = &now
	} else {
		execution.AbandonedAt = &now
	}
	return s.executionRepo.Save(ctx, execution)
}

func validCoordinates(latitude float64, longitude float64) bool {
	return latitude >= -90 && latitude <= 90 && longitude >= -180 && longitude <= 180
}

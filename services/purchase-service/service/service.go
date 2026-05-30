package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"purchase-service/models"
	"purchase-service/repository"
	"purchase-service/rpc"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrInvalidInput = errors.New("invalid input")

type PurchaseService struct {
	repo    *repository.PurchaseRepository
	tourRPC *rpc.TourRPCClient
}

func NewPurchaseService(repo *repository.PurchaseRepository, tourRPC *rpc.TourRPCClient) *PurchaseService {
	return &PurchaseService{
		repo:    repo,
		tourRPC: tourRPC,
	}
}

func (s *PurchaseService) GetCart(ctx context.Context, touristID string) (models.ShoppingCart, error) {
	touristID = strings.TrimSpace(touristID)
	if touristID == "" {
		return models.ShoppingCart{}, fmt.Errorf("%w: touristId is required", ErrInvalidInput)
	}
	return s.repo.GetOrCreateCart(ctx, touristID)
}

func (s *PurchaseService) AddItem(ctx context.Context, touristID string, item models.OrderItem) (models.ShoppingCart, error) {
	touristID = strings.TrimSpace(touristID)
	item.TourID = strings.TrimSpace(item.TourID)
	item.TourName = strings.TrimSpace(item.TourName)
	if touristID == "" || item.TourID == "" || item.TourName == "" {
		return models.ShoppingCart{}, fmt.Errorf("%w: touristId, tourId and tourName are required", ErrInvalidInput)
	}

	tour, err := s.getTour(ctx, item.TourID)
	if err != nil {
		return models.ShoppingCart{}, err
	}
	if !isPurchasable(tour) {
		return models.ShoppingCart{}, fmt.Errorf("%w: only published non-archived tours can be added to cart", ErrInvalidInput)
	}
	if tour.Price <= 0 {
		return models.ShoppingCart{}, fmt.Errorf("%w: tour price must be greater than 0", ErrInvalidInput)
	}
	owned, err := s.repo.HasToken(ctx, touristID, item.TourID)
	if err != nil {
		return models.ShoppingCart{}, err
	}
	if owned {
		return models.ShoppingCart{}, fmt.Errorf("%w: tour is already purchased", ErrInvalidInput)
	}

	cart, err := s.repo.GetOrCreateCart(ctx, touristID)
	if err != nil {
		return models.ShoppingCart{}, err
	}
	for _, existing := range cart.Items {
		if existing.TourID == item.TourID {
			return models.ShoppingCart{}, fmt.Errorf("%w: tour is already in cart", ErrInvalidInput)
		}
	}

	cart.Items = append(cart.Items, models.OrderItem{
		TourID:   item.TourID,
		TourName: tour.Name,
		Price:    tour.Price,
	})
	return s.repo.SaveCart(ctx, cart)
}

func (s *PurchaseService) RemoveItem(ctx context.Context, touristID string, tourID string) (models.ShoppingCart, error) {
	cart, err := s.GetCart(ctx, touristID)
	if err != nil {
		return models.ShoppingCart{}, err
	}
	tourID = strings.TrimSpace(tourID)
	items := make([]models.OrderItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		if item.TourID != tourID {
			items = append(items, item)
		}
	}
	cart.Items = items
	return s.repo.SaveCart(ctx, cart)
}

func (s *PurchaseService) Checkout(ctx context.Context, touristID string) ([]models.TourPurchaseToken, error) {
	cart, err := s.GetCart(ctx, touristID)
	if err != nil {
		return []models.TourPurchaseToken{}, err
	}
	if len(cart.Items) == 0 {
		return []models.TourPurchaseToken{}, fmt.Errorf("%w: cart is empty", ErrInvalidInput)
	}

	for _, item := range cart.Items {
		valid, message, err := s.tourRPC.ValidateTourPurchase(ctx, item.TourID)
		if err != nil {
			return []models.TourPurchaseToken{}, err
		}
		if !valid {
			return []models.TourPurchaseToken{}, fmt.Errorf("%w: tour %s is no longer purchasable: %s", ErrInvalidInput, item.TourName, message)
		}
	}

	now := time.Now().UTC()
	created := make([]models.TourPurchaseToken, 0, len(cart.Items))
	for _, item := range cart.Items {
		owned, err := s.repo.HasToken(ctx, touristID, item.TourID)
		if err != nil {
			return []models.TourPurchaseToken{}, err
		}
		if owned {
			continue
		}
		tokenValue, err := generateToken()
		if err != nil {
			return []models.TourPurchaseToken{}, err
		}
		token := models.TourPurchaseToken{
			ID:          primitive.NewObjectID(),
			TouristID:   touristID,
			TourID:      item.TourID,
			PurchasedAt: now,
			Token:       tokenValue,
		}
		if err = s.repo.CreateToken(ctx, token); err != nil {
			return []models.TourPurchaseToken{}, err
		}
		created = append(created, token)
	}

	cart.Items = []models.OrderItem{}
	if _, err = s.repo.SaveCart(ctx, cart); err != nil {
		return []models.TourPurchaseToken{}, err
	}
	return created, nil
}

func (s *PurchaseService) GetTokens(ctx context.Context, touristID string) ([]models.TourPurchaseToken, error) {
	touristID = strings.TrimSpace(touristID)
	if touristID == "" {
		return []models.TourPurchaseToken{}, fmt.Errorf("%w: touristId is required", ErrInvalidInput)
	}
	return s.repo.GetTokens(ctx, touristID)
}

func (s *PurchaseService) IsPurchased(ctx context.Context, touristID string, tourID string) (bool, error) {
	touristID = strings.TrimSpace(touristID)
	tourID = strings.TrimSpace(tourID)
	if touristID == "" || tourID == "" {
		return false, fmt.Errorf("%w: touristId and tourId are required", ErrInvalidInput)
	}
	return s.repo.HasToken(ctx, touristID, tourID)
}

func (s *PurchaseService) getTour(ctx context.Context, tourID string) (models.TourSnapshot, error) {
	return s.tourRPC.GetTourForPurchase(ctx, tourID)
}

func isPurchasable(tour models.TourSnapshot) bool {
	return tour.Status == "published" && !tour.Archived
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

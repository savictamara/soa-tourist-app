package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
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
	// SAGA T1: load cart
	cart, err := s.GetCart(ctx, touristID)
	if err != nil {
		return []models.TourPurchaseToken{}, err
	}
	if len(cart.Items) == 0 {
		return []models.TourPurchaseToken{}, fmt.Errorf("%w: cart is empty", ErrInvalidInput)
	}

	// SAGA T2: validate each tour via Tour Service RPC
	// Compensation C2: remove invalid items from cart so tourist can retry with valid ones
	validItems := make([]models.OrderItem, 0, len(cart.Items))
	invalidItems := make([]string, 0)
	for _, item := range cart.Items {
		valid, message, err := s.tourRPC.ValidateTourPurchase(ctx, item.TourID)
		if err != nil {
			return []models.TourPurchaseToken{}, err
		}
		if !valid {
			log.Printf("Checkout SAGA T2 tour no longer purchasable tourId=%s tourName=%s reason=%s — compensation C2: removing from cart", item.TourID, item.TourName, message)
			invalidItems = append(invalidItems, item.TourName)
		} else {
			validItems = append(validItems, item)
		}
	}

	// execute compensation C2: persist removal of invalid items
	if len(invalidItems) > 0 {
		cart.Items = validItems
		if _, err = s.repo.SaveCart(ctx, cart); err != nil {
			return []models.TourPurchaseToken{}, err
		}
		if len(validItems) == 0 {
			return []models.TourPurchaseToken{}, fmt.Errorf("%w: no purchasable tours in cart, removed: %s", ErrInvalidInput, strings.Join(invalidItems, ", "))
		}
		log.Printf("Checkout SAGA compensation C2 complete: removed %d invalid tours, continuing with %d valid", len(invalidItems), len(validItems))
	}

	// SAGA T3: create purchase tokens for each valid item
	// Compensation C3: if token creation fails, delete all tokens created so far in this checkout
	now := time.Now().UTC()
	createdTokenIDs := make([]primitive.ObjectID, 0, len(validItems))
	created := make([]models.TourPurchaseToken, 0, len(validItems))

	for _, item := range validItems {
		owned, err := s.repo.HasToken(ctx, touristID, item.TourID)
		if err != nil {
			s.compensateTokens(ctx, createdTokenIDs)
			return []models.TourPurchaseToken{}, err
		}
		if owned {
			continue
		}
		tokenValue, err := generateToken()
		if err != nil {
			s.compensateTokens(ctx, createdTokenIDs)
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
			s.compensateTokens(ctx, createdTokenIDs)
			return []models.TourPurchaseToken{}, err
		}
		createdTokenIDs = append(createdTokenIDs, token.ID)
		created = append(created, token)
	}

	// SAGA T4: clear cart
	// Compensation C4: if cart clearing fails, delete all created tokens
	cart.Items = []models.OrderItem{}
	if _, err = s.repo.SaveCart(ctx, cart); err != nil {
		log.Printf("Checkout SAGA T4 failed touristId=%s — executing compensation C4: deleting %d tokens", touristID, len(createdTokenIDs))
		s.compensateTokens(ctx, createdTokenIDs)
		return []models.TourPurchaseToken{}, err
	}

	log.Printf("Checkout SAGA completed touristId=%s tokens=%d", touristID, len(created))
	return created, nil
}

func (s *PurchaseService) compensateTokens(ctx context.Context, ids []primitive.ObjectID) {
	if len(ids) == 0 {
		return
	}
	log.Printf("Checkout SAGA compensation C3/C4: deleting %d tokens", len(ids))
	if err := s.repo.DeleteTokensBatch(ctx, ids); err != nil {
		log.Printf("Checkout SAGA compensation failed to delete tokens: %v", err)
	}
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

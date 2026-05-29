package repository

import (
	"context"
	"time"

	"purchase-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PurchaseRepository struct {
	carts  *mongo.Collection
	tokens *mongo.Collection
}

func NewPurchaseRepository(carts *mongo.Collection, tokens *mongo.Collection) *PurchaseRepository {
	return &PurchaseRepository{carts: carts, tokens: tokens}
}

func (r *PurchaseRepository) GetOrCreateCart(ctx context.Context, touristID string) (models.ShoppingCart, error) {
	var cart models.ShoppingCart
	err := r.carts.FindOne(ctx, bson.M{"touristId": touristID}).Decode(&cart)
	if err == nil {
		normalizeCart(&cart)
		return cart, nil
	}
	if err != mongo.ErrNoDocuments {
		return models.ShoppingCart{}, err
	}

	now := time.Now().UTC()
	cart = models.ShoppingCart{
		TouristID:  touristID,
		Items:      []models.OrderItem{},
		TotalPrice: 0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	_, err = r.carts.InsertOne(ctx, cart)
	if err != nil {
		return models.ShoppingCart{}, err
	}
	return r.GetOrCreateCart(ctx, touristID)
}

func (r *PurchaseRepository) SaveCart(ctx context.Context, cart models.ShoppingCart) (models.ShoppingCart, error) {
	cart.TotalPrice = calculateTotal(cart.Items)
	cart.UpdatedAt = time.Now().UTC()
	_, err := r.carts.UpdateOne(ctx, bson.M{"touristId": cart.TouristID}, bson.M{"$set": bson.M{
		"items":      cart.Items,
		"totalPrice": cart.TotalPrice,
		"updatedAt":  cart.UpdatedAt,
	}})
	if err != nil {
		return models.ShoppingCart{}, err
	}
	return r.GetOrCreateCart(ctx, cart.TouristID)
}

func (r *PurchaseRepository) HasToken(ctx context.Context, touristID string, tourID string) (bool, error) {
	count, err := r.tokens.CountDocuments(ctx, bson.M{"touristId": touristID, "tourId": tourID})
	return count > 0, err
}

func (r *PurchaseRepository) CreateToken(ctx context.Context, token models.TourPurchaseToken) error {
	_, err := r.tokens.InsertOne(ctx, token)
	if mongo.IsDuplicateKeyError(err) {
		return nil
	}
	return err
}

func (r *PurchaseRepository) GetTokens(ctx context.Context, touristID string) ([]models.TourPurchaseToken, error) {
	cur, err := r.tokens.Find(ctx, bson.M{"touristId": touristID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var tokens []models.TourPurchaseToken
	if err = cur.All(ctx, &tokens); err != nil {
		return nil, err
	}
	if tokens == nil {
		return []models.TourPurchaseToken{}, nil
	}
	return tokens, nil
}

func calculateTotal(items []models.OrderItem) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Price
	}
	return total
}

func normalizeCart(cart *models.ShoppingCart) {
	if cart.Items == nil {
		cart.Items = []models.OrderItem{}
	}
	cart.TotalPrice = calculateTotal(cart.Items)
}

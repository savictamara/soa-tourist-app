package repository

import (
	"context"

	"tour-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TourExecutionRepository struct {
	collection *mongo.Collection
}

func NewTourExecutionRepository(collection *mongo.Collection) *TourExecutionRepository {
	return &TourExecutionRepository{collection: collection}
}

func (r *TourExecutionRepository) Create(ctx context.Context, execution models.TourExecution) (models.TourExecution, error) {
	res, err := r.collection.InsertOne(ctx, execution)
	if err != nil {
		return models.TourExecution{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		execution.ID = oid
	}
	normalizeExecution(&execution)
	return execution, nil
}

func (r *TourExecutionRepository) GetByID(ctx context.Context, id primitive.ObjectID) (models.TourExecution, error) {
	var execution models.TourExecution
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&execution)
	if err == nil {
		normalizeExecution(&execution)
	}
	return execution, err
}

func (r *TourExecutionRepository) GetActiveByTouristID(ctx context.Context, touristID string) (models.TourExecution, error) {
	findOpts := options.FindOne().SetSort(bson.D{{Key: "startedAt", Value: -1}})
	var execution models.TourExecution
	err := r.collection.FindOne(ctx, bson.M{"touristId": touristID, "status": "active"}, findOpts).Decode(&execution)
	if err == nil {
		normalizeExecution(&execution)
	}
	return execution, err
}

func (r *TourExecutionRepository) GetActiveByTouristAndTour(ctx context.Context, touristID string, tourID primitive.ObjectID) (models.TourExecution, error) {
	var execution models.TourExecution
	err := r.collection.FindOne(ctx, bson.M{"touristId": touristID, "tourId": tourID, "status": "active"}).Decode(&execution)
	if err == nil {
		normalizeExecution(&execution)
	}
	return execution, err
}

func (r *TourExecutionRepository) GetLatestByTouristAndTour(ctx context.Context, touristID string, tourID primitive.ObjectID) (models.TourExecution, error) {
	findOpts := options.FindOne().SetSort(bson.D{{Key: "startedAt", Value: -1}})
	var execution models.TourExecution
	err := r.collection.FindOne(ctx, bson.M{"touristId": touristID, "tourId": tourID}, findOpts).Decode(&execution)
	if err == nil {
		normalizeExecution(&execution)
	}
	return execution, err
}

func (r *TourExecutionRepository) Save(ctx context.Context, execution models.TourExecution) (models.TourExecution, error) {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": execution.ID}, execution)
	if err != nil {
		return models.TourExecution{}, err
	}
	return r.GetByID(ctx, execution.ID)
}

func normalizeExecution(execution *models.TourExecution) {
	if execution.CompletedKeyPoints == nil {
		execution.CompletedKeyPoints = []models.CompletedKeyPoint{}
	}
}

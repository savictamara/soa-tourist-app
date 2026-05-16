package repository

import (
	"context"
	"log"
	"time"

	"tour-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TourRepository struct {
	collection *mongo.Collection
}

func NewTourRepository(collection *mongo.Collection) *TourRepository {
	return &TourRepository{collection: collection}
}

func (r *TourRepository) Create(ctx context.Context, tour models.Tour) (models.Tour, error) {
	res, err := r.collection.InsertOne(ctx, tour)
	if err != nil {
		return models.Tour{}, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		tour.ID = oid
	}

	return tour, nil
}

func (r *TourRepository) GetByID(ctx context.Context, id primitive.ObjectID) (models.Tour, error) {
	var tour models.Tour
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&tour)
	if err == nil {
		if tour.Tags == nil {
			tour.Tags = []string{}
		}
		if tour.KeyPoints == nil {
			tour.KeyPoints = []models.KeyPoint{}
		}
		if tour.Reviews == nil {
			tour.Reviews = []models.Review{}
		}
	}
	return tour, err
}

func (r *TourRepository) GetByAuthorID(ctx context.Context, authorID string) ([]models.Tour, error) {
	findOpts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cur, err := r.collection.Find(ctx, bson.M{"authorId": authorID}, findOpts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var tours []models.Tour
	if err = cur.All(ctx, &tours); err != nil {
		return nil, err
	}
	if tours == nil {
		return []models.Tour{}, nil
	}
	for i := range tours {
		if tours[i].Tags == nil {
			tours[i].Tags = []string{}
		}
		if tours[i].KeyPoints == nil {
			tours[i].KeyPoints = []models.KeyPoint{}
		}
		if tours[i].Reviews == nil {
			tours[i].Reviews = []models.Review{}
		}
	}
	return tours, nil
}

func (r *TourRepository) GetAll(ctx context.Context) ([]models.Tour, error) {
	findOpts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cur, err := r.collection.Find(ctx, bson.M{}, findOpts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var tours []models.Tour
	if err = cur.All(ctx, &tours); err != nil {
		return nil, err
	}
	if tours == nil {
		return []models.Tour{}, nil
	}
	for i := range tours {
		if tours[i].Tags == nil {
			tours[i].Tags = []string{}
		}
		if tours[i].KeyPoints == nil {
			tours[i].KeyPoints = []models.KeyPoint{}
		}
		if tours[i].Reviews == nil {
			tours[i].Reviews = []models.Review{}
		}
	}
	return tours, nil
}

func (r *TourRepository) AddKeyPoint(ctx context.Context, tourID primitive.ObjectID, keyPoint models.KeyPoint) (models.Tour, error) {
	update := bson.M{
		"$push": bson.M{
			"keyPoints": keyPoint,
		},
		"$set": bson.M{
			"updatedAt": keyPoint.UpdatedAt,
		},
	}

	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": tourID}, update)
	if err != nil {
		return models.Tour{}, err
	}
	log.Printf("Add key point update result: tourId=%s MatchedCount=%d ModifiedCount=%d", tourID.Hex(), res.MatchedCount, res.ModifiedCount)
	if res.MatchedCount == 0 {
		log.Printf("Add key point update: tour not found for id=%s", tourID.Hex())
		return models.Tour{}, mongo.ErrNoDocuments
	}

	return r.GetByID(ctx, tourID)
}

func (r *TourRepository) AddReview(ctx context.Context, tourID primitive.ObjectID, review models.Review, updatedAt time.Time) error {
	update := bson.M{
		"$push": bson.M{
			"reviews": review,
		},
		"$set": bson.M{
			"updatedAt": updatedAt,
		},
	}
	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": tourID}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

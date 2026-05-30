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
		normalizeTour(&tour)
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
	normalizeTours(tours)
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
	normalizeTours(tours)
	return tours, nil
}

func (r *TourRepository) GetPublished(ctx context.Context) ([]models.Tour, error) {
	findOpts := options.Find().SetSort(bson.D{{Key: "publishedAt", Value: -1}, {Key: "createdAt", Value: -1}})
	cur, err := r.collection.Find(ctx, bson.M{"status": "published"}, findOpts)
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
	normalizeTours(tours)
	return tours, nil
}

func (r *TourRepository) GetAvailableForTourists(ctx context.Context) ([]models.Tour, error) {
	findOpts := options.Find().SetSort(bson.D{{Key: "publishedAt", Value: -1}, {Key: "createdAt", Value: -1}})
	cur, err := r.collection.Find(ctx, bson.M{"status": bson.M{"$in": []string{"published", "archived"}}}, findOpts)
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
	normalizeTours(tours)
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

func (r *TourRepository) UpdateKeyPoint(ctx context.Context, tourID primitive.ObjectID, keyPointID primitive.ObjectID, req models.UpdateKeyPointRequest, now time.Time) (models.Tour, error) {
	update := bson.M{
		"$set": bson.M{
			"keyPoints.$[kp].name":        req.Name,
			"keyPoints.$[kp].description": req.Description,
			"keyPoints.$[kp].latitude":    req.Latitude,
			"keyPoints.$[kp].longitude":   req.Longitude,
			"keyPoints.$[kp].imageUrl":    req.ImageURL,
			"keyPoints.$[kp].updatedAt":   now,
			"updatedAt":                   now,
		},
	}
	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"kp._id": keyPointID}},
	}
	opts := options.Update().SetArrayFilters(arrayFilters)

	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": tourID}, update, opts)
	if err != nil {
		return models.Tour{}, err
	}
	if res.MatchedCount == 0 {
		return models.Tour{}, mongo.ErrNoDocuments
	}
	return r.GetByID(ctx, tourID)
}

func (r *TourRepository) DeleteKeyPoint(ctx context.Context, tourID primitive.ObjectID, keyPointID primitive.ObjectID, now time.Time) (models.Tour, error) {
	tour, err := r.GetByID(ctx, tourID)
	if err != nil {
		return models.Tour{}, err
	}

	remaining := make([]models.KeyPoint, 0, len(tour.KeyPoints))
	for _, kp := range tour.KeyPoints {
		if kp.ID != keyPointID {
			remaining = append(remaining, kp)
		}
	}
	for i := range remaining {
		remaining[i].Order = i + 1
	}

	update := bson.M{
		"$set": bson.M{
			"keyPoints": remaining,
			"updatedAt": now,
		},
	}
	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": tourID}, update)
	if err != nil {
		return models.Tour{}, err
	}
	if res.MatchedCount == 0 {
		return models.Tour{}, mongo.ErrNoDocuments
	}
	return r.GetByID(ctx, tourID)
}

func (r *TourRepository) UpdateLength(ctx context.Context, tourID primitive.ObjectID, lengthKm float64, updatedAt time.Time) (models.Tour, error) {
	update := bson.M{
		"$set": bson.M{
			"lengthKm":  lengthKm,
			"updatedAt": updatedAt,
		},
	}
	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": tourID}, update)
	if err != nil {
		return models.Tour{}, err
	}
	if res.MatchedCount == 0 {
		return models.Tour{}, mongo.ErrNoDocuments
	}
	return r.GetByID(ctx, tourID)
}

func (r *TourRepository) UpdateDurations(ctx context.Context, tourID primitive.ObjectID, durations []models.TourDuration, updatedAt time.Time) (models.Tour, error) {
	update := bson.M{
		"$set": bson.M{
			"durations": durations,
			"updatedAt": updatedAt,
		},
	}
	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": tourID}, update)
	if err != nil {
		return models.Tour{}, err
	}
	if res.MatchedCount == 0 {
		return models.Tour{}, mongo.ErrNoDocuments
	}
	return r.GetByID(ctx, tourID)
}

func (r *TourRepository) UpdatePrice(ctx context.Context, tourID primitive.ObjectID, price float64, updatedAt time.Time) (models.Tour, error) {
	update := bson.M{
		"$set": bson.M{
			"price":     price,
			"updatedAt": updatedAt,
		},
	}
	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": tourID}, update)
	if err != nil {
		return models.Tour{}, err
	}
	if res.MatchedCount == 0 {
		return models.Tour{}, mongo.ErrNoDocuments
	}
	return r.GetByID(ctx, tourID)
}

func (r *TourRepository) UpdateLifecycle(ctx context.Context, tourID primitive.ObjectID, update bson.M) (models.Tour, error) {
	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": tourID}, update)
	if err != nil {
		return models.Tour{}, err
	}
	if res.MatchedCount == 0 {
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

func normalizeTours(tours []models.Tour) {
	for i := range tours {
		normalizeTour(&tours[i])
	}
}

func normalizeTour(tour *models.Tour) {
	if tour.Tags == nil {
		tour.Tags = []string{}
	}
	if tour.KeyPoints == nil {
		tour.KeyPoints = []models.KeyPoint{}
	}
	if tour.Reviews == nil {
		tour.Reviews = []models.Review{}
	}
	if tour.Durations == nil {
		tour.Durations = []models.TourDuration{}
	}
	if tour.Status == "" {
		tour.Status = "draft"
	}
}

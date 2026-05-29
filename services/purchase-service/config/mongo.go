package config

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultMongoURI = "mongodb://localhost:27018"
	databaseName    = "tourist_app_purchases"
)

type MongoConfig struct {
	Client         *mongo.Client
	Database       *mongo.Database
	Carts          *mongo.Collection
	PurchaseTokens *mongo.Collection
}

func ConnectMongo() (*MongoConfig, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = defaultMongoURI
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	db := client.Database(databaseName)
	cfg := &MongoConfig{
		Client:         client,
		Database:       db,
		Carts:          db.Collection("shopping_carts"),
		PurchaseTokens: db.Collection("purchase_tokens"),
	}
	return cfg, ensureIndexes(ctx, cfg)
}

func ensureIndexes(ctx context.Context, cfg *MongoConfig) error {
	_, err := cfg.PurchaseTokens.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "touristId", Value: 1}, {Key: "tourId", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("ux_tourist_tour"),
	})
	if err != nil {
		return err
	}
	_, err = cfg.Carts.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "touristId", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("ux_cart_tourist"),
	})
	return err
}

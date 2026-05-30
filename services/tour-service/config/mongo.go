package config

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultMongoURI          = "mongodb://localhost:27017"
	databaseName             = "tourist_app_tours"
	collectionName           = "tours"
	executionsCollectionName = "tour_executions"
)

type MongoConfig struct {
	Client               *mongo.Client
	Database             *mongo.Database
	Collection           *mongo.Collection
	ExecutionsCollection *mongo.Collection
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
	return &MongoConfig{
		Client:               client,
		Database:             db,
		Collection:           db.Collection(collectionName),
		ExecutionsCollection: db.Collection(executionsCollectionName),
	}, nil
}

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KeyPoint struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Latitude    float64            `bson:"latitude" json:"latitude"`
	Longitude   float64            `bson:"longitude" json:"longitude"`
	ImageURL    string             `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"`
	Order       int                `bson:"order" json:"order"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Tour struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthorID    string             `bson:"authorId" json:"authorId"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Difficulty  string             `bson:"difficulty" json:"difficulty"`
	Tags        []string           `bson:"tags" json:"tags"`
	Status      string             `bson:"status" json:"status"`
	Price       float64            `bson:"price" json:"price"`
	KeyPoints   []KeyPoint         `bson:"keyPoints" json:"keyPoints"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type CreateTourRequest struct {
	AuthorID    string   `json:"authorId"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Difficulty  string   `json:"difficulty"`
	Tags        []string `json:"tags"`
}

type AddKeyPointRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ImageURL    string  `json:"imageUrl"`
}

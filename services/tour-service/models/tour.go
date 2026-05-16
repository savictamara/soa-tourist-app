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
	Reviews     []Review           `bson:"reviews" json:"reviews"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Review struct {
	ID              primitive.ObjectID `bson:"id" json:"id"`
	Rating          int                `bson:"rating" json:"rating"`
	Comment         string             `bson:"comment" json:"comment"`
	TouristID       string             `bson:"touristId" json:"touristId"`
	TouristUsername string             `bson:"touristUsername" json:"touristUsername"`
	VisitedDate     string             `bson:"visitedDate" json:"visitedDate"`
	CommentDate     time.Time          `bson:"commentDate" json:"commentDate"`
	Images          []string           `bson:"images" json:"images"`
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

type CreateReviewRequest struct {
	Rating          int      `json:"rating"`
	Comment         string   `json:"comment"`
	TouristID       string   `json:"touristId"`
	TouristUsername string   `json:"touristUsername"`
	VisitedDate     string   `json:"visitedDate"`
	Images          []string `json:"images"`
}

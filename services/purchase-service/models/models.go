package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShoppingCart struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TouristID  string             `bson:"touristId" json:"touristId"`
	Items      []OrderItem        `bson:"items" json:"items"`
	TotalPrice float64            `bson:"totalPrice" json:"totalPrice"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type OrderItem struct {
	TourID   string  `bson:"tourId" json:"tourId"`
	TourName string  `bson:"tourName" json:"tourName"`
	Price    float64 `bson:"price" json:"price"`
}

type TourPurchaseToken struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TouristID   string             `bson:"touristId" json:"touristId"`
	TourID      string             `bson:"tourId" json:"tourId"`
	PurchasedAt time.Time          `bson:"purchasedAt" json:"purchasedAt"`
	Token       string             `bson:"token" json:"token"`
}

type TourSnapshot struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Status   string  `json:"status"`
	Archived bool    `json:"archived"`
	Price    float64 `json:"price"`
}

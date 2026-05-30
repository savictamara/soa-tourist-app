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
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthorID      string             `bson:"authorId" json:"authorId"`
	Name          string             `bson:"name" json:"name"`
	Description   string             `bson:"description" json:"description"`
	Difficulty    string             `bson:"difficulty" json:"difficulty"`
	Tags          []string           `bson:"tags" json:"tags"`
	Status        string             `bson:"status" json:"status"`
	PublishedAt   *time.Time         `bson:"publishedAt,omitempty" json:"publishedAt,omitempty"`
	ArchivedAt    *time.Time         `bson:"archivedAt,omitempty" json:"archivedAt,omitempty"`
	ReactivatedAt *time.Time         `bson:"reactivatedAt,omitempty" json:"reactivatedAt,omitempty"`
	LengthKm      float64            `bson:"lengthKm" json:"lengthKm"`
	Durations     []TourDuration     `bson:"durations" json:"durations"`
	Price         float64            `bson:"price" json:"price"`
	KeyPoints     []KeyPoint         `bson:"keyPoints" json:"keyPoints"`
	Reviews       []Review           `bson:"reviews" json:"reviews"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type TourDuration struct {
	TransportType string `bson:"transportType" json:"transportType"`
	Minutes       int    `bson:"minutes" json:"minutes"`
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

type CompletedKeyPoint struct {
	KeyPointID   primitive.ObjectID `bson:"keyPointId" json:"keyPointId"`
	KeyPointName string             `bson:"keyPointName" json:"keyPointName"`
	ReachedAt    time.Time          `bson:"reachedAt" json:"reachedAt"`
}

type TourExecution struct {
	ID                 primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	TourID             primitive.ObjectID  `bson:"tourId" json:"tourId"`
	TourName           string              `bson:"tourName" json:"tourName"`
	TouristID          string              `bson:"touristId" json:"touristId"`
	Status             string              `bson:"status" json:"status"`
	StartedAt          time.Time           `bson:"startedAt" json:"startedAt"`
	CompletedAt        *time.Time          `bson:"completedAt,omitempty" json:"completedAt,omitempty"`
	AbandonedAt        *time.Time          `bson:"abandonedAt,omitempty" json:"abandonedAt,omitempty"`
	LastActivityAt     time.Time           `bson:"lastActivityAt" json:"lastActivityAt"`
	StartLatitude      float64             `bson:"startLatitude" json:"startLatitude"`
	StartLongitude     float64             `bson:"startLongitude" json:"startLongitude"`
	CurrentLatitude    float64             `bson:"currentLatitude" json:"currentLatitude"`
	CurrentLongitude   float64             `bson:"currentLongitude" json:"currentLongitude"`
	CompletedKeyPoints []CompletedKeyPoint `bson:"completedKeyPoints" json:"completedKeyPoints"`
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

type UpdateKeyPointRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ImageURL    string  `json:"imageUrl"`
}

type UpdateDurationsRequest struct {
	Durations []TourDuration `json:"durations"`
}

type UpdatePriceRequest struct {
	Price float64 `json:"price"`
}

type CreateReviewRequest struct {
	Rating          int      `json:"rating"`
	Comment         string   `json:"comment"`
	TouristID       string   `json:"touristId"`
	TouristUsername string   `json:"touristUsername"`
	VisitedDate     string   `json:"visitedDate"`
	Images          []string `json:"images"`
}

type StartTourExecutionRequest struct {
	TouristID string  `json:"touristId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CheckTourExecutionLocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CheckTourExecutionLocationResponse struct {
	Execution        TourExecution      `json:"execution"`
	KeyPointReached  bool               `json:"keyPointReached"`
	ReachedKeyPoint  *CompletedKeyPoint `json:"reachedKeyPoint,omitempty"`
	DistanceMeters   float64            `json:"distanceMeters"`
	LastActivityAt   time.Time          `json:"lastActivityAt"`
	CompletedCount   int                `json:"completedCount"`
	TotalKeyPointCnt int                `json:"totalKeyPointCount"`
}

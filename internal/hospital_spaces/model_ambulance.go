package hospital_spaces

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ambulance represents an ambulance in the system
type Ambulance struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	AmbulanceID string             `json:"ambulance_id" bson:"ambulance_id"`
	Name        string             `json:"name" bson:"name" binding:"required"`
	Location    string             `json:"location" bson:"location" binding:"required"`
	Status      string             `json:"status" bson:"status"`
	Type        string             `json:"type" bson:"type" binding:"required"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// AmbulanceCreateRequest represents the request for creating a new ambulance
type AmbulanceCreateRequest struct {
	Name     string `json:"name" bson:"name" binding:"required"`
	Type     string `json:"type" bson:"type" binding:"required"`
	Location string `json:"location" bson:"location" binding:"required"`
}

// NewAmbulance creates a new Ambulance with default values
func NewAmbulance(req AmbulanceCreateRequest) *Ambulance {
	now := time.Now()
	return &Ambulance{
		AmbulanceID: uuid.New().String(),
		Name:        req.Name,
		Type:        req.Type,
		Location:    req.Location,
		Status:      "available",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

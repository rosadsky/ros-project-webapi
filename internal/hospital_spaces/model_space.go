package hospital_spaces

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Space represents a hospital space/room
type Space struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SpaceID      string             `json:"space_id" bson:"space_id"`
	Name         string             `json:"name" bson:"name" binding:"required"`
	Type         string             `json:"type" bson:"type" binding:"required"`
	Floor        int                `json:"floor" bson:"floor" binding:"required"`
	Capacity     int                `json:"capacity" bson:"capacity" binding:"required"`
	Status       string             `json:"status" bson:"status"`
	AssignedTo   *string            `json:"assigned_to,omitempty" bson:"assigned_to,omitempty"`
	AssignedType *string            `json:"assigned_type,omitempty" bson:"assigned_type,omitempty"`
	AssignedID   *string            `json:"assigned_id,omitempty" bson:"assigned_id,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

// SpaceCreateRequest represents the request for creating a new space
type SpaceCreateRequest struct {
	Name     string `json:"name" bson:"name" binding:"required"`
	Type     string `json:"type" bson:"type" binding:"required"`
	Floor    int    `json:"floor" bson:"floor" binding:"required"`
	Capacity int    `json:"capacity" bson:"capacity" binding:"required"`
}

// SpaceUpdateRequest represents the request for updating a space
type SpaceUpdateRequest struct {
	AssignedTo   *string `json:"assigned_to,omitempty" bson:"assigned_to,omitempty"`
	AssignedType *string `json:"assigned_type,omitempty" bson:"assigned_type,omitempty"`
	AssignedID   *string `json:"assigned_id,omitempty" bson:"assigned_id,omitempty"`
}

// NewSpace creates a new Space with default values
func NewSpace(req SpaceCreateRequest) *Space {
	now := time.Now()
	return &Space{
		SpaceID:   uuid.New().String(),
		Name:      req.Name,
		Type:      req.Type,
		Floor:     req.Floor,
		Capacity:  req.Capacity,
		Status:    "available",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateAssignment updates the space assignment
func (s *Space) UpdateAssignment(req SpaceUpdateRequest) {
	if req.AssignedTo != nil && *req.AssignedTo != "" {
		s.AssignedTo = req.AssignedTo
		s.AssignedType = req.AssignedType
		s.AssignedID = req.AssignedID
		s.Status = "occupied"
	} else {
		s.AssignedTo = nil
		s.AssignedType = nil
		s.AssignedID = nil
		s.Status = "available"
	}
	s.UpdatedAt = time.Now()
}

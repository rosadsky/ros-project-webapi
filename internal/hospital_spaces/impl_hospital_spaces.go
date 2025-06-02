package hospital_spaces

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rosadsky/ros-project-backend/internal/db_service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionSpaces     = "spaces"
	collectionAmbulances = "ambulances"
	ErrNoDocuments       = "no documents found"
)

// SpaceServiceImpl implements the space service operations
type SpaceServiceImpl struct {
	dbService *db_service.DbService
}

// NewSpaceServiceImpl creates a new space service implementation
func NewSpaceServiceImpl(dbService *db_service.DbService) *SpaceServiceImpl {
	return &SpaceServiceImpl{
		dbService: dbService,
	}
}

// CreateSpace creates a new hospital space
// @Summary Create a new hospital space
// @Description Create a new hospital space with the specified details
// @Tags Spaces
// @Accept json
// @Produce json
// @Param space body SpaceCreateRequest true "Space creation details"
// @Success 201 {object} Space "Space created successfully"
// @Failure 400 {object} map[string]string "Bad request - invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/spaces [post]
func (s *SpaceServiceImpl) CreateSpace(c *gin.Context) {
	var request SpaceCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	space := NewSpace(request)
	collection := s.dbService.GetCollection(collectionSpaces)
	ctx, cancel := s.dbService.CreateContext()
	defer cancel()

	result, err := collection.InsertOne(ctx, space)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create space: %v", err)})
		return
	}

	space.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, space)
}

// GetSpaces retrieves all hospital spaces
// @Summary Get all hospital spaces
// @Description Retrieve a list of all hospital spaces with their current status and assignments
// @Tags Spaces
// @Accept json
// @Produce json
// @Success 200 {array} Space "List of hospital spaces"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/spaces [get]
func (s *SpaceServiceImpl) GetSpaces(c *gin.Context) {
	collection := s.dbService.GetCollection(collectionSpaces)
	ctx, cancel := s.dbService.CreateContext()
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve spaces: %v", err)})
		return
	}
	defer cursor.Close(ctx)

	var spaces []Space
	if err := cursor.All(ctx, &spaces); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to decode spaces: %v", err)})
		return
	}

	if spaces == nil {
		spaces = []Space{}
	}

	c.JSON(http.StatusOK, spaces)
}

// UpdateSpace updates a hospital space assignment
// @Summary Update a hospital space
// @Description Update space assignment details such as assigned entity, type, and ID
// @Tags Spaces
// @Accept json
// @Produce json
// @Param id path string true "The unique space ID (UUID format)" format(uuid)
// @Param space body SpaceUpdateRequest true "Space update details"
// @Success 200 {object} Space "Space updated successfully"
// @Failure 400 {object} map[string]string "Bad request - invalid space ID or input"
// @Failure 404 {object} map[string]string "Space not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/spaces/{id} [put]
func (s *SpaceServiceImpl) UpdateSpace(c *gin.Context) {
	spaceIDStr := c.Param("id")
	// Validate that it's a valid UUID format
	if _, err := uuid.Parse(spaceIDStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid space ID"})
		return
	}

	var request SpaceUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := s.dbService.GetCollection(collectionSpaces)
	ctx, cancel := s.dbService.CreateContext()
	defer cancel()

	// Find the space first
	var space Space
	filter := bson.M{"space_id": spaceIDStr}
	err := collection.FindOne(ctx, filter).Decode(&space)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Space not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to find space: %v", err)})
		return
	}

	// Update the assignment
	space.UpdateAssignment(request)

	// Update in database with correct field names
	update := bson.M{
		"$set": bson.M{
			"assigned_to":   space.AssignedTo,
			"assigned_type": space.AssignedType,
			"assigned_id":   space.AssignedID,
			"status":        space.Status,
			"updated_at":    space.UpdatedAt,
		},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update space: %v", err)})
		return
	}

	c.JSON(http.StatusOK, space)
}

// DeleteSpace deletes a hospital space
// @Summary Delete a hospital space
// @Description Remove a hospital space from the system
// @Tags Spaces
// @Accept json
// @Produce json
// @Param id path string true "The unique space ID (UUID format)" format(uuid)
// @Success 204 "Space deleted successfully"
// @Failure 400 {object} map[string]string "Bad request - invalid space ID"
// @Failure 404 {object} map[string]string "Space not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/spaces/{id} [delete]
func (s *SpaceServiceImpl) DeleteSpace(c *gin.Context) {
	spaceIDStr := c.Param("id")
	// Validate that it's a valid UUID format
	if _, err := uuid.Parse(spaceIDStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid space ID"})
		return
	}

	collection := s.dbService.GetCollection(collectionSpaces)
	ctx, cancel := s.dbService.CreateContext()
	defer cancel()

	filter := bson.M{"space_id": spaceIDStr}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete space: %v", err)})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Space not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateAmbulance creates a new ambulance
// @Summary Create a new ambulance
// @Description Register a new ambulance in the system
// @Tags Ambulances
// @Accept json
// @Produce json
// @Param ambulance body AmbulanceCreateRequest true "Ambulance creation details"
// @Success 201 {object} Ambulance "Ambulance created successfully"
// @Failure 400 {object} map[string]string "Bad request - invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/ambulances [post]
func (s *SpaceServiceImpl) CreateAmbulance(c *gin.Context) {
	var request AmbulanceCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ambulance := NewAmbulance(request)
	collection := s.dbService.GetCollection(collectionAmbulances)
	ctx, cancel := s.dbService.CreateContext()
	defer cancel()

	result, err := collection.InsertOne(ctx, ambulance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create ambulance: %v", err)})
		return
	}

	ambulance.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, ambulance)
}

// GetAmbulances retrieves all ambulances
// @Summary Get all ambulances
// @Description Retrieve a list of all ambulances in the system
// @Tags Ambulances
// @Accept json
// @Produce json
// @Success 200 {array} Ambulance "List of ambulances"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/ambulances [get]
func (s *SpaceServiceImpl) GetAmbulances(c *gin.Context) {
	collection := s.dbService.GetCollection(collectionAmbulances)
	ctx, cancel := s.dbService.CreateContext()
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve ambulances: %v", err)})
		return
	}
	defer cursor.Close(ctx)

	var ambulances []Ambulance
	if err := cursor.All(ctx, &ambulances); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to decode ambulances: %v", err)})
		return
	}

	if ambulances == nil {
		ambulances = []Ambulance{}
	}

	c.JSON(http.StatusOK, ambulances)
}

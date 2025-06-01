package hospital_spaces

import (
	"github.com/gin-gonic/gin"
	"github.com/rosadsky/ros-project-backend/internal/db_service"
)

// SpaceAPIRouter provides space API routing
type SpaceAPIRouter struct {
	spaceService *SpaceServiceImpl
}

// NewSpaceAPIRouter creates a new space API router
func NewSpaceAPIRouter(dbService *db_service.DbService) *SpaceAPIRouter {
	return &SpaceAPIRouter{
		spaceService: NewSpaceServiceImpl(dbService),
	}
}

// RegisterRoutes registers all space-related routes
func (router *SpaceAPIRouter) RegisterRoutes(engine *gin.Engine) {
	api := engine.Group("/api")
	{
		// Space routes - 4 simple CRUD endpoints
		spaces := api.Group("/spaces")
		{
			spaces.POST("", router.spaceService.CreateSpace)       // CREATE
			spaces.GET("", router.spaceService.GetSpaces)          // READ (all)
			spaces.PUT("/:id", router.spaceService.UpdateSpace)    // UPDATE
			spaces.DELETE("/:id", router.spaceService.DeleteSpace) // DELETE
		}

		// Ambulance routes for assignment support
		ambulances := api.Group("/ambulances")
		{
			ambulances.POST("", router.spaceService.CreateAmbulance)
			ambulances.GET("", router.spaceService.GetAmbulances)
		}
	}

	// Health check endpoint
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "hospital-spaces-api",
		})
	})
}

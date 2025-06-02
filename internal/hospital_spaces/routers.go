package hospital_spaces

import (
	"github.com/gin-gonic/gin"
	"github.com/rosadsky/ros-project-backend/internal/db_service"
)

type SpaceAPIRouter struct {
	spaceService *SpaceServiceImpl
}

func NewSpaceAPIRouter(dbService *db_service.DbService) *SpaceAPIRouter {
	return &SpaceAPIRouter{
		spaceService: NewSpaceServiceImpl(dbService),
	}
}

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

		ambulances := api.Group("/ambulances")
		{
			ambulances.POST("", router.spaceService.CreateAmbulance)
			ambulances.GET("", router.spaceService.GetAmbulances)
		}
	}

	// Health check endpoint
	// @Summary Health check
	// @Description Check the health status of the API service
	// @Tags Health
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string "Service is healthy"
	// @Router /api/health [get]
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "hospital-spaces-api",
		})
	})
}

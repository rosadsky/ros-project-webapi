// @title Hospital Spaces API
// @version 1.0.0
// @description RESTful API for managing hospital spaces and ambulances. This service provides comprehensive management of hospital room assignments, space allocation, and ambulance tracking.
// @description
// @description ## Features
// @description - Hospital space management (CRUD operations)
// @description - Ambulance management
// @description - Space assignment and status tracking
// @description - Health monitoring
// @contact.name ROS Project Backend
// @contact.url https://github.com/rosadsky/ros-project-backend
// @license.name MIT
// @host localhost:8080
// @BasePath /
// @schemes http
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rosadsky/ros-project-backend/internal/db_service"
	"github.com/rosadsky/ros-project-backend/internal/hospital_spaces"
	"github.com/rs/zerolog"

	// Swagger imports
	_ "github.com/rosadsky/ros-project-backend/docs" // Import generated docs
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Initialize logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Initialize database service
	dbService := db_service.NewDbService()
	defer func() {
		if err := dbService.Disconnect(); err != nil {
			logger.Error().Err(err).Msg("Failed to disconnect from database")
		}
	}()

	// Make index creation optional for development
	if err := dbService.EnsureIndexes(); err != nil {
		log.Printf("Warning: Failed to create database indexes: %v", err)
		log.Println("Continuing without indexes for development...")
		// Don't exit, just continue
	} else {
		log.Println("Database indexes created successfully")
	}

	// Create Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3333",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	}))

	// Add other middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Initialize and register routes
	spaceRouter := hospital_spaces.NewSpaceAPIRouter(dbService)
	spaceRouter.RegisterRoutes(router)

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Get port from environment
	port := os.Getenv("AMBULANCE_API_PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info().Str("port", port).Msg("Starting Hospital Spaces API server")
		logger.Info().Msg("Swagger UI available at: http://localhost:8080/swagger/index.html")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info().Msg("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server exited")
}

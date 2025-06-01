package db_service

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DbService provides database connection and operations
type DbService struct {
	Client   *mongo.Client
	Database *mongo.Database
	timeout  time.Duration
}

var (
	dbService *DbService
	once      sync.Once
)

// NewDbService creates a new database service instance
func NewDbService() *DbService {
	once.Do(func() {
		dbService = &DbService{
			timeout: 10 * time.Second,
		}
		dbService.connect()
	})
	return dbService
}

// connect establishes connection to MongoDB
func (db *DbService) connect() {
	mongoURI := os.Getenv("AMBULANCE_API_MONGODB_URI")
	if mongoURI == "" {
		// Build MongoDB URI from individual environment variables
		mongoHost := os.Getenv("AMBULANCE_API_MONGODB_HOST")
		if mongoHost == "" {
			mongoHost = "localhost"
		}

		mongoPort := os.Getenv("AMBULANCE_API_MONGODB_PORT")
		if mongoPort == "" {
			mongoPort = "27017"
		}

		mongoUser := os.Getenv("AMBULANCE_API_MONGODB_USERNAME")
		mongoPassword := os.Getenv("AMBULANCE_API_MONGODB_PASSWORD")

		if mongoUser != "" && mongoPassword != "" {
			mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%s", mongoUser, mongoPassword, mongoHost, mongoPort)
		} else {
			mongoURI = fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort)
		}
	}

	databaseName := os.Getenv("AMBULANCE_API_MONGODB_DATABASE")
	if databaseName == "" {
		databaseName = "hospital-spaces" // default database name
	}

	log.Printf("Connecting to MongoDB at %s", mongoURI)

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	db.Client = client
	db.Database = client.Database(databaseName)
	log.Printf("Successfully connected to MongoDB database: %s", databaseName)
}

// GetCollection returns a MongoDB collection
func (db *DbService) GetCollection(collectionName string) *mongo.Collection {
	return db.Database.Collection(collectionName)
}

// Disconnect closes the database connection
func (db *DbService) Disconnect() error {
	if db.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
		defer cancel()
		return db.Client.Disconnect(ctx)
	}
	return nil
}

// CreateContext creates a context with timeout for database operations
func (db *DbService) CreateContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), db.timeout)
}

// EnsureIndexes creates necessary indexes for collections
func (db *DbService) EnsureIndexes() error {
	ctx, cancel := db.CreateContext()
	defer cancel()

	// Create indexes for spaces collection
	spacesCollection := db.GetCollection("spaces")
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "name", Value: 1},
				{Key: "type", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "floor", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "space_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := spacesCollection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		// Log warning but don't fail the application
		log.Printf("Warning: failed to create spaces indexes: %v", err)
		return nil // Don't return error, just warn
	}

	// Create indexes for ambulances collection
	ambulancesCollection := db.GetCollection("ambulances")
	ambulanceIndexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "name", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "ambulance_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err = ambulancesCollection.Indexes().CreateMany(ctx, ambulanceIndexModels)
	if err != nil {
		// Log warning but don't fail the application
		log.Printf("Warning: failed to create ambulances indexes: %v", err)
		return nil // Don't return error, just warn
	}

	log.Println("Database indexes created successfully")
	return nil
}

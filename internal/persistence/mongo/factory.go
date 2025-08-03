package mongo

import (
	"context"
	"fmt"

	"github.com/thisausername99/pantry_butler/config"
)

// // StartMongo connects to MongoDB and returns the client and database handle.
// func StartMongo() (*mongo.Client, *mongo.Database, error) {
// 	logger := logging.GetLogger()
// 	// mongoHost := os.Getenv("MONGO_HOST")
// 	// mongoPort := os.Getenv("MONGO_PORT")
// 	// mongoUser := os.Getenv("MONGO_USER")
// 	// mongoPassword := os.Getenv("MONGO_PASSWORD")
// 	// mongoDb := os.Getenv("MONGO_DB")

// 	// if mongoHost == "" {
// 	// 	mongoHost = "localhost"
// 	// }
// 	dbName := os.Getenv("MONGODB_DB")
// 	if dbName == "" {
// 		dbName = "pantry_butler_dev"
// 	}
// 	// Form the URI connection string
// 	// If you change MONGO_DB or MONGODB_DB, make sure your docker-compose and environment variables are correct and rebuild with --no-cache.
// 	uri := os.Getenv("MONGO_URI")
// 	logger.Info("Connecting to MongoDB...", zap.String("uri", uri))
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
// 	if err != nil {
// 		logger.Error("Error connecting to MongoDB", zap.Error(err))
// 		return nil, nil, err
// 	}
// 	// Ping to verify connection
// 	if err := client.Ping(ctx, nil); err != nil {
// 		logger.Error("Error pinging MongoDB", zap.Error(err))
// 		return nil, nil, err
// 	}
// 	logger.Info("Connected to MongoDB", zap.String("dbName", dbName))
// 	db := client.Database(dbName)
// 	return client, db, nil
// }

func NewMongoConnection(cfg *config.MongoDBConfig) (MongoDB, error) {
	// Create client with configuration options
	client, err := NewMongoClient(cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}

	// Test connection with configured timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	return client, nil
}

package mongo

import (
	"context"
	"os"
	"time"

	"github.com/thisausername99/pantry-butler/internal/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// StartMongo connects to MongoDB and returns the client and database handle.
func StartMongo() (*mongo.Client, *mongo.Database, error) {
	logger := logging.GetLogger()
	// mongoHost := os.Getenv("MONGO_HOST")
	// mongoPort := os.Getenv("MONGO_PORT")
	// mongoUser := os.Getenv("MONGO_USER")
	// mongoPassword := os.Getenv("MONGO_PASSWORD")
	// mongoDb := os.Getenv("MONGO_DB")

	// if mongoHost == "" {
	// 	mongoHost = "localhost"
	// }
	dbName := os.Getenv("MONGODB_DB")
	if dbName == "" {
		dbName = "pantry_butler_dev"
	}
	// Form the URI connection string
	// If you change MONGO_DB or MONGODB_DB, make sure your docker-compose and environment variables are correct and rebuild with --no-cache.
	uri := os.Getenv("MONGO_URI")
	logger.Info("Connecting to MongoDB...", zap.String("uri", uri))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		logger.Error("Error connecting to MongoDB", zap.Error(err))
		return nil, nil, err
	}
	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		logger.Error("Error pinging MongoDB", zap.Error(err))
		return nil, nil, err
	}
	logger.Info("Connected to MongoDB", zap.String("dbName", dbName))
	db := client.Database(dbName)
	return client, db, nil
}

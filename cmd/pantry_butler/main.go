package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/thisausername99/pantry_butler/config"
	"github.com/thisausername99/pantry_butler/internal/delivery/http"
	"github.com/thisausername99/pantry_butler/internal/persistence/mongo"
	"github.com/thisausername99/pantry_butler/internal/usecase"
	"github.com/thisausername99/pantry_butler/pkg/logging"
	"go.uber.org/zap"
)

const defaultPort = "8080"

func main() {
	log := logging.GetLogger()
	log.Info("Starting Pantry Butler server...")

	// Connect to MongoDB
	config := config.Load()
	mongoClient, err := mongo.NewMongoConnection(&config.MongoDB)

	if err != nil {
		log.Error("error connecting to db", zap.Error(err))
		os.Exit(1)
	} else {
		log.Info("Connected to MongoDB successfully!")
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Error("Error disconnecting MongoDB client", zap.Error(err))
		} else {
			log.Info("MongoDB client disconnected.")
		}
	}()

	// Setup collections
	pantryEntryCollection := mongoClient.Database(config.MongoDB.Database).Collection("pantries")
	recipeCollection := mongoClient.Database(config.MongoDB.Database).Collection("recipes")
	userCollection := mongoClient.Database(config.MongoDB.Database).Collection("users")

	// Setup use case
	uc := &usecase.Usecase{
		Logger: log,
		RepoWrapper: usecase.RepoWrapper{
			RecipeRepo: &mongo.RecipeRepo{Collection: recipeCollection, Logger: log},
			PantryRepo: &mongo.PantryEntryRepo{Collection: pantryEntryCollection, Logger: log},
			UserRepo:   &mongo.UserRepo{Collection: userCollection, Logger: log},
		},
	}

	// Create HTTP server with Gin
	server := http.NewServer(log, uc)

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Info("Starting HTTP server", zap.String("port", defaultPort))
		if err := server.Start(defaultPort); err != nil {
			log.Error("Failed to start server", zap.Error(err))
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	<-quit
	log.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		log.Error("Error during server shutdown", zap.Error(err))
	}

	log.Info("Server stopped")
}

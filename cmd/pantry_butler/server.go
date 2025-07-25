package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/thisausername99/pantry-butler/internal/adapter/delivery/graphql"
	"github.com/thisausername99/pantry-butler/internal/adapter/persistence/mongo"
	"github.com/thisausername99/pantry-butler/internal/logging"
	"go.uber.org/zap"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/thisausername99/pantry-butler/internal/usecase"
)

const defaultPort = "5000"

func main() {
	log := logging.GetLogger()
	log.Info("Starting Pantry Butler server...")

	client, db, err := mongo.StartMongo()
	if err != nil {
		log.Error("error connecting to db", zap.Error(err))
		os.Exit(1)
	} else {
		log.Info("Connected to MongoDB successfully!")
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Error("Error disconnecting MongoDB client", zap.Error(err))
		} else {
			log.Info("MongoDB client disconnected.")
		}
	}()

	pantryEntryCollection := db.Collection("pantry_entries")
	recipeCollection := db.Collection("recipes")

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	uc := &usecase.Usecase{
		Logger: log,
		RepoWrapper: usecase.RepoWrapper{
			RecipeRepo:      &mongo.RecipeRepo{Collection: recipeCollection, Logger: log},
			PantryEntryRepo: &mongo.PantryEntryRepo{Collection: pantryEntryCollection, Logger: log},
		},
	}

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{UseCase: *uc}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Info(fmt.Sprintf("connect to http://localhost:%s/ for GraphQL playground", port))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Failed to start HTTP server", zap.Error(err))
	}
}

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
)

const defaultPort = "5000"

func main() {
	log := logging.GetLogger()
	client, db, err := mongo.StartMongo()
	if err != nil {
		panic(fmt.Errorf("error connecting to db: %w", err))
	} else {
		fmt.Printf("Connected successfully!\n")
	}
	defer client.Disconnect(context.Background())

	pantryEntryCollection := db.Collection("pantry_entries")
	recipeCollection := db.Collection("recipes")

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{
		RecipeRepo:      mongo.RecipeRepo{Collection: recipeCollection},
		PantryEntryRepo: mongo.PantryEntryRepo{Collection: pantryEntryCollection},
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Info("connect to http://localhost:%s/ for GraphQL playground", zap.String("port", port))
	http.ListenAndServe(":"+port, nil)
}

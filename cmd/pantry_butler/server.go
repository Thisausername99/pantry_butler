package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/thisausername99/pantry-butler/internal/adapter/delivery/graphql"
	"github.com/thisausername99/pantry-butler/internal/adapter/persistence/mongo"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	DB, err := mongo.StartDB()
	if err != nil {
		panic(fmt.Errorf("error connecting to db"))
	} else {
		fmt.Printf("Connected succesfully!")
	}

	defer DB.Close()

	DB.AddQueryHook(mongo.DBLogger{})

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{
		RecipeRepo:      mongo.RecipeRepo{DB: DB},
		PantryEntryRepo: mongo.PantryEntryRepo{DB: DB},
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

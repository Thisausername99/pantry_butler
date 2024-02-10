package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/thisausername99/recipes-api/graph"
	"github.com/thisausername99/recipes-api/postgres"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	DB, err := postgres.StartDB()
	if err != nil {
		panic(fmt.Errorf("Error connecting to DB"))
	}

	defer DB.Close()

	DB.AddQueryHook(postgres.DBLogger{})

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		RecipesRepo:       postgres.RecipesRepo{DB: DB},
		PantryEntriesRepo: postgres.PantryEntriesRepo{DB: DB},
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

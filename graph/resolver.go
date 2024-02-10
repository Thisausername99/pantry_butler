package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/thisausername99/recipes-api/postgres"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// TODO:
// Add database or http client here

type Resolver struct {
	RecipesRepo       postgres.RecipesRepo
	PantryEntriesRepo postgres.PantryEntriesRepo
}

package graphql

import "github.com/thisausername99/pantry-butler/internal/adapter/persistence/mongo"

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	RecipeRepo      mongo.RecipeRepo
	PantryEntryRepo mongo.PantryEntryRepo
}

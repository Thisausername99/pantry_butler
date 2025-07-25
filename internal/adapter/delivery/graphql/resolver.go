package graphql

import (
	"github.com/thisausername99/pantry-butler/internal/usecase"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UseCase usecase.Usecase // Add the MatchUsecase here for use in resolvers
}

package usecase

import (
	context "context"

	entity "github.com/thisausername99/pantry-butler/internal/domain"
	"go.uber.org/zap"
)

type RecipeRepository interface {
	GetRecipes(ctx context.Context) ([]*entity.Recipe, error)
	GetRecipesByCuisine(ctx context.Context, cuisine string) ([]*entity.Recipe, error)
}

type PantryEntryRepository interface {
	GetPantryEntries(ctx context.Context) ([]*entity.PantryEntry, error)
	InsertPantryEntry(ctx context.Context, entry *entity.PantryEntryInput) (*entity.PantryEntry, error)
}

type RepoWrapper struct {
	RecipeRepo      RecipeRepository
	PantryEntryRepo PantryEntryRepository
	// Add more repositories as needed
}

type Usecase struct {
	RepoWrapper RepoWrapper
	Logger      *zap.Logger
}

package repository

import (
	"context"

	"github.com/thisausername99/pantry_butler/internal/domain/entity"
)

type RecipeRepository interface {
	GetRecipes(ctx context.Context) ([]entity.Recipe, error)
	GetRecipesByCuisine(ctx context.Context, cuisine string) ([]entity.Recipe, error)
}

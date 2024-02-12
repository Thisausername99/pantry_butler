package graph

import (
	"context"

	"github.com/thisausername99/recipes-api/models"
)

type queryResolver struct{ *Resolver }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Recipe is the resolver for the recipe field.
func (r *queryResolver) Recipe(ctx context.Context) ([]*models.Recipe, error) {
	return r.RecipeRepo.GetRecipes()
}

// RecipeByCuisine is the resolver for the recipeByCuisine field.
func (r *queryResolver) RecipeByCuisine(ctx context.Context, cuisine string) ([]*models.Recipe, error) {
	return r.RecipeRepo.GetRecipesByCuisine(cuisine)
}

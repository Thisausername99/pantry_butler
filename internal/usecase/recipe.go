package usecase

import (
	"context"

	entity "github.com/thisausername99/pantry-butler/internal/domain"
)

func (u *Usecase) GetAllRecipes(ctx context.Context) ([]*entity.Recipe, error) {
	return u.RepoWrapper.RecipeRepo.GetRecipes(ctx)
}

func (u *Usecase) GetRecipeByCuisine(ctx context.Context, cuisine string) ([]*entity.Recipe, error) {
	return u.RepoWrapper.RecipeRepo.GetRecipesByCuisine(ctx, cuisine)
}

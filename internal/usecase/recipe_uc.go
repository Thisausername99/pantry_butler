package usecase

import (
	"context"

	"github.com/thisausername99/pantry_butler/internal/domain/entity"
)

func (u *Usecase) GetAllRecipes(ctx context.Context) ([]*entity.Recipe, error) {
	return u.RepoWrapper.RecipeRepo.GetRecipes(ctx)
}

func (u *Usecase) GetRecipeByCuisine(ctx context.Context, cuisine string) ([]*entity.Recipe, error) {
	return u.RepoWrapper.RecipeRepo.GetRecipesByCuisine(ctx, cuisine)
}

func (u *Usecase) FindMatchingRecipes(ctx context.Context, pantryID string) ([]*entity.Recipe, error) {
	recipes, err := u.RepoWrapper.RecipeRepo.GetRecipes(ctx)
	if err != nil {
		return nil, err
	}

	pantryEntries, err := u.RepoWrapper.PantryRepo.GetPantryEntries(ctx, pantryID)
	if err != nil {
		return nil, err
	}

	pantrySet := make(map[string]struct{})
	for _, entry := range pantryEntries {
		pantrySet[entry.Name] = struct{}{}
	}

	var matches []*entity.Recipe
	for _, recipe := range recipes {
		allIngredientsPresent := true
		for ingredient := range recipe.Ingredients {
			if _, ok := pantrySet[ingredient]; !ok {
				allIngredientsPresent = false
				break
			}
		}
		if allIngredientsPresent {
			matches = append(matches, recipe)
		}
	}
	return matches, nil
}

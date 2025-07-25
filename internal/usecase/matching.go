package usecase

import (
	"context"

	entity "github.com/thisausername99/pantry-butler/internal/domain"
)

func (u *Usecase) FindMatchingRecipes(ctx context.Context) ([]*entity.Recipe, error) {
	recipes, err := u.RepoWrapper.RecipeRepo.GetRecipes(ctx)
	if err != nil {
		return nil, err
	}
	pantryEntries, err := u.RepoWrapper.PantryEntryRepo.GetPantryEntries(ctx)
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

package usecase

import (
	context "context"

	entity "github.com/thisausername99/pantry-butler/internal/domain"
)

type RecipeRepository interface {
	GetRecipes(ctx context.Context) ([]*entity.Recipe, error)
}

type PantryEntryRepository interface {
	GetPantryEntries(ctx context.Context) ([]*entity.PantryEntry, error)
}

type MatchUsecase struct {
	RecipeRepo      RecipeRepository
	PantryEntryRepo PantryEntryRepository
}

// FindMatchingRecipes returns recipes that can be made with the available pantry entries
func (u *MatchUsecase) FindMatchingRecipes(ctx context.Context) ([]*entity.Recipe, error) {
	recipes, err := u.RecipeRepo.GetRecipes(ctx)
	if err != nil {
		return nil, err
	}
	pantryEntries, err := u.PantryEntryRepo.GetPantryEntries(ctx)
	if err != nil {
		return nil, err
	}

	// Build a set of available pantry item names
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

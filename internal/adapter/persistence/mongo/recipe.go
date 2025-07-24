package mongo

import (
	"github.com/go-pg/pg/v10"
	entity "github.com/thisausername99/pantry-butler/internal/domain"
)

type RecipeRepo struct {
	DB *pg.DB
}

func (m *RecipeRepo) GetRecipes() ([]*entity.Recipe, error) {
	var recipe []*entity.Recipe
	err := m.DB.Model(&recipe).Select()
	if err != nil {
		return nil, err
	}
	return recipe, nil
}

func (m *RecipeRepo) GetRecipesByCuisine(cuisine string) ([]*entity.Recipe, error) {
	var recipe []*entity.Recipe
	err := m.DB.Model(&recipe).Where("cuisine = ?", cuisine).Select()
	if err != nil {
		return nil, err
	}
	return recipe, nil
}

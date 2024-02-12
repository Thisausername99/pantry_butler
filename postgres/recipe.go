package postgres

import (
	"github.com/go-pg/pg/v10"
	"github.com/thisausername99/recipes-api/models"
)

type RecipeRepo struct {
	DB *pg.DB
}

func (m *RecipeRepo) GetRecipes() ([]*models.Recipe, error) {
	var recipe []*models.Recipe
	err := m.DB.Model(&recipe).Select()
	if err != nil {
		return nil, err
	}
	return recipe, nil
}

func (m *RecipeRepo) GetRecipesByCuisine(cuisine string) ([]*models.Recipe, error) {
	var recipe []*models.Recipe
	err := m.DB.Model(&recipe).Where("cuisine = ?", cuisine).Select()
	if err != nil {
		return nil, err
	}
	return recipe, nil
}

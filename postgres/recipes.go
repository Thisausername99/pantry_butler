package postgres

import (
	"github.com/go-pg/pg/v10"
	"github.com/thisausername99/recipes-api/graph/model"
)

type RecipesRepo struct {
	DB *pg.DB
}

func (m *RecipesRepo) GetRecipes() ([]*model.Recipe, error) {
	var recipes []*model.Recipe
	err := m.DB.Model(&recipes).Select()
	if err != nil {
		return nil, err
	}
	return recipes, nil
}

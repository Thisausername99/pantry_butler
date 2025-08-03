package mongo

import (
	"context"

	"github.com/thisausername99/pantry_butler/internal/domain/entity"
	"github.com/thisausername99/pantry_butler/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type RecipeRepo struct {
	Collection MongoCollection
	Logger     *zap.Logger
}

// Ensure it implements the interface
var _ repository.RecipeRepository = (*RecipeRepo)(nil)

func (m *RecipeRepo) GetRecipes(ctx context.Context) ([]entity.Recipe, error) {
	var recipes []entity.Recipe
	cursor, err := m.Collection.Find(ctx, bson.M{}) // Empty filter = get ALL documents
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var recipe entity.Recipe
		if err := cursor.Decode(&recipe); err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return recipes, nil
}

func (m *RecipeRepo) GetRecipesByCuisine(ctx context.Context, cuisine string) ([]entity.Recipe, error) {
	var recipes []entity.Recipe
	filter := bson.M{"cuisine": cuisine}
	cursor, err := m.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var recipe entity.Recipe
		if err := cursor.Decode(&recipe); err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return recipes, nil
}

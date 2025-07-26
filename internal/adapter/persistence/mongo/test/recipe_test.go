package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	mongorepo "github.com/thisausername99/pantry-butler/internal/adapter/persistence/mongo"
)

func setupRecipeCollection(t *testing.T) *mongo.Collection {
	collection := testDB.Collection("recipes")
	// Clear the collection before each test
	_, err := collection.DeleteMany(context.Background(), map[string]interface{}{})
	require.NoError(t, err)
	return collection
}

func TestRecipeRepo_GetRecipes(t *testing.T) {
	collection := setupRecipeCollection(t)
	logger := zap.NewNop()
	repo := &mongorepo.RecipeRepo{
		Collection: collection,
		Logger:     logger,
	}

	// Insert test data
	testRecipes := []interface{}{
		map[string]interface{}{
			"name":        "Apple Pie",
			"cuisine":     "American",
			"ingredients": map[string]interface{}{"apples": 5, "flour": "2 cups"},
			"difficulty":  3,
			"description": "A delicious apple pie",
		},
		map[string]interface{}{
			"name":        "Pasta Carbonara",
			"cuisine":     "Italian",
			"ingredients": map[string]interface{}{"pasta": "500g", "eggs": 4},
			"difficulty":  2,
			"description": "Classic Italian pasta dish",
		},
	}

	_, err := collection.InsertMany(context.Background(), testRecipes)
	require.NoError(t, err)

	// Test retrieval
	ctx := context.Background()
	recipes, err := repo.GetRecipes(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, recipes, 2)
	assert.Equal(t, "Apple Pie", recipes[0].Name)
	assert.Equal(t, "Pasta Carbonara", recipes[1].Name)
	assert.Equal(t, "American", *recipes[0].Cuisine)
	assert.Equal(t, "Italian", *recipes[1].Cuisine)
}

func TestRecipeRepo_GetRecipesByCuisine(t *testing.T) {
	collection := setupRecipeCollection(t)
	logger := zap.NewNop()
	repo := &mongorepo.RecipeRepo{
		Collection: collection,
		Logger:     logger,
	}

	// Insert test data
	testRecipes := []interface{}{
		map[string]interface{}{
			"name":        "Apple Pie",
			"cuisine":     "American",
			"ingredients": map[string]interface{}{"apples": 5, "flour": "2 cups"},
			"difficulty":  3,
			"description": "A delicious apple pie",
		},
		map[string]interface{}{
			"name":        "Pasta Carbonara",
			"cuisine":     "Italian",
			"ingredients": map[string]interface{}{"pasta": "500g", "eggs": 4},
			"difficulty":  2,
			"description": "Classic Italian pasta dish",
		},
		map[string]interface{}{
			"name":        "Grilled Cheese",
			"cuisine":     "American",
			"ingredients": map[string]interface{}{"bread": 2, "cheese": "4 slices"},
			"difficulty":  1,
			"description": "Simple grilled cheese sandwich",
		},
	}

	_, err := collection.InsertMany(context.Background(), testRecipes)
	require.NoError(t, err)

	// Test retrieval by cuisine
	ctx := context.Background()
	americanRecipes, err := repo.GetRecipesByCuisine(ctx, "American")

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, americanRecipes, 2)
	assert.Equal(t, "Apple Pie", americanRecipes[0].Name)
	assert.Equal(t, "Grilled Cheese", americanRecipes[1].Name)

	// Test Italian cuisine
	italianRecipes, err := repo.GetRecipesByCuisine(ctx, "Italian")
	assert.NoError(t, err)
	assert.Len(t, italianRecipes, 1)
	assert.Equal(t, "Pasta Carbonara", italianRecipes[0].Name)
}

func TestRecipeRepo_GetRecipes_Empty(t *testing.T) {
	collection := setupRecipeCollection(t)
	logger := zap.NewNop()
	repo := &mongorepo.RecipeRepo{
		Collection: collection,
		Logger:     logger,
	}

	// Test retrieval from empty collection
	ctx := context.Background()
	recipes, err := repo.GetRecipes(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, recipes, 0)
}

func TestRecipeRepo_GetRecipesByCuisine_NotFound(t *testing.T) {
	collection := setupRecipeCollection(t)
	logger := zap.NewNop()
	repo := &mongorepo.RecipeRepo{
		Collection: collection,
		Logger:     logger,
	}

	// Test retrieval for non-existent cuisine
	ctx := context.Background()
	recipes, err := repo.GetRecipesByCuisine(ctx, "French")

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, recipes, 0)
}

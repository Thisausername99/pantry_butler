package test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	entity "github.com/thisausername99/pantry-butler/internal/domain"
	"github.com/thisausername99/pantry-butler/internal/usecase"
)

func TestGetAllPantryEntries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repository
	mockPantryRepo := NewMockPantryEntryRepository(ctrl)
	mockRecipeRepo := NewMockRecipeRepository(ctrl)

	// Create usecase instance
	usecaseInstance := &usecase.Usecase{
		RepoWrapper: usecase.RepoWrapper{
			PantryEntryRepo: mockPantryRepo,
			RecipeRepo:      mockRecipeRepo,
		},
		Logger: zap.NewNop(),
	}

	// Test data
	expectedEntries := []*entity.PantryEntry{
		{
			ID:         "1",
			Name:       "Apples",
			Quantity:   intPtr(5),
			Expiration: timePtr(time.Now().AddDate(0, 0, 7)),
		},
		{
			ID:         "2",
			Name:       "Milk",
			Quantity:   intPtr(1),
			Expiration: timePtr(time.Now().AddDate(0, 0, 3)),
		},
	}

	// Set up expectations
	ctx := context.Background()
	mockPantryRepo.EXPECT().
		GetPantryEntries(ctx).
		Return(expectedEntries, nil).
		Times(1)

	// Execute test
	result, err := usecaseInstance.GetAllPantryEntries(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedEntries, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "Apples", result[0].Name)
	assert.Equal(t, "Milk", result[1].Name)
}

func TestGetAllPantryEntries_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repository
	mockPantryRepo := NewMockPantryEntryRepository(ctrl)
	mockRecipeRepo := NewMockRecipeRepository(ctrl)

	// Create usecase instance
	usecaseInstance := &usecase.Usecase{
		RepoWrapper: usecase.RepoWrapper{
			PantryEntryRepo: mockPantryRepo,
			RecipeRepo:      mockRecipeRepo,
		},
		Logger: zap.NewNop(),
	}

	// Set up expectations for error case
	ctx := context.Background()
	expectedError := assert.AnError
	mockPantryRepo.EXPECT().
		GetPantryEntries(ctx).
		Return(nil, expectedError).
		Times(1)

	// Execute test
	result, err := usecaseInstance.GetAllPantryEntries(ctx)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func timePtr(t time.Time) *time.Time {
	return &t
}

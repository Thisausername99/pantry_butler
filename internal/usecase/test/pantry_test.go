package test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/thisausername99/pantry_butler/internal/domain/entity"
	m "github.com/thisausername99/pantry_butler/internal/mocks"
	"github.com/thisausername99/pantry_butler/internal/usecase"
)

var testPantryID = "testPantryID"

// Global test variables
var (
	mockCtrl        *gomock.Controller
	mockPantryRepo  *m.MockPantryRepository
	mockRecipeRepo  *m.MockRecipeRepository
	usecaseInstance *usecase.Usecase
)

// setupTest initializes the global test instance
func setupTest(t *testing.T) {
	mockCtrl = gomock.NewController(t)
	mockPantryRepo = m.NewMockPantryRepository(mockCtrl)
	mockRecipeRepo = m.NewMockRecipeRepository(mockCtrl)

	usecaseInstance = &usecase.Usecase{
		RepoWrapper: usecase.RepoWrapper{
			PantryRepo: mockPantryRepo,
			RecipeRepo: mockRecipeRepo,
		},
		Logger: zap.NewNop(),
	}
}

// teardownTest cleans up after tests
func teardownTest() {
	if mockCtrl != nil {
		mockCtrl.Finish()
	}
}

func TestGetAllPantryEntries(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	// Test data
	expectedEntries := []entity.PantryEntry{
		{
			ID:         "1",
			Name:       "Apples",
			Quantity:   float64Ptr(5),
			Expiration: timePtr(time.Now().AddDate(0, 0, 7)),
		},
		{
			ID:         "2",
			Name:       "Milk",
			Quantity:   float64Ptr(1),
			Expiration: timePtr(time.Now().AddDate(0, 0, 3)),
		},
	}

	// Set up expectations
	ctx := context.Background()
	mockPantryRepo.EXPECT().
		GetPantryEntries(ctx, testPantryID).
		Return(expectedEntries, nil).
		Times(1)

	// Execute test
	result, err := usecaseInstance.GetAllPantryEntries(ctx, testPantryID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedEntries, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "Apples", result[0].Name)
	assert.Equal(t, "Milk", result[1].Name)
}

func TestGetAllPantryEntries_Error(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	// Set up expectations for error case
	ctx := context.Background()
	expectedError := assert.AnError
	mockPantryRepo.EXPECT().
		GetPantryEntries(ctx, testPantryID).
		Return(nil, expectedError).
		Times(1)

	// Execute test
	result, err := usecaseInstance.GetAllPantryEntries(ctx, testPantryID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestGetAllPantryEntries_Empty(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	// Set up expectations for empty result
	ctx := context.Background()
	mockPantryRepo.EXPECT().
		GetPantryEntries(ctx, testPantryID).
		Return([]entity.PantryEntry{}, nil).
		Times(1)

	// Execute test
	result, err := usecaseInstance.GetAllPantryEntries(ctx, testPantryID)

	// Assertions
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Len(t, result, 0)
}

// Helper functions
func float64Ptr(f float64) *float64 {
	return &f
}

func timePtr(t time.Time) *time.Time {
	return &t
}

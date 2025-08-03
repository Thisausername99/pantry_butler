package test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/thisausername99/pantry_butler/internal/domain/entity"
	"github.com/thisausername99/pantry_butler/internal/mocks"
	mongo "github.com/thisausername99/pantry_butler/internal/persistence/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func TestPantryEntryRepo_WithMockedMongo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock MongoDB components
	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockSingleResult := mocks.NewMockMongoSingleResult(ctrl)

	// Create repository with mocked collection
	repo := &mongo.PantryEntryRepo{
		Collection: mockCollection,
		Logger:     zap.NewNop(),
	}

	// Set up expectations for GetPantryEntries
	ctx := context.Background()
	pantryID := "test-pantry-id"
	filter := bson.M{"id": pantryID}

	// Mock the FindOne operation
	mockCollection.EXPECT().
		FindOne(ctx, filter).
		Return(mockSingleResult).
		Times(1)

	// Mock single result behavior
	mockSingleResult.EXPECT().Decode(gomock.Any()).DoAndReturn(func(val interface{}) error {
		// Simulate decoding a pantry document with entries
		pantry := val.(*entity.Pantry)
		pantry.ID = pantryID
		pantry.Name = "Test Pantry"
		pantry.OwnerID = "user1"
		pantry.CreatedAt = time.Now()

		// Create entries
		entries := []entity.PantryEntry{
			{
				ID:       "1",
				Name:     "Apples",
				Quantity: float64Ptr(5),
			},
			{
				ID:       "2",
				Name:     "Milk",
				Quantity: float64Ptr(1),
			},
		}
		pantry.Entries = &entries
		return nil
	}).Times(1)

	// Execute test
	result, err := repo.GetPantryEntries(ctx, pantryID)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Apples", result[0].Name)
	assert.Equal(t, "Milk", result[1].Name)
}

// Helper functions
// func intPtr(i int) *int {
// 	return &i
// }

func float64Ptr(f float64) *float64 {
	return &f
}

func stringPtr(s string) *string {
	return &s
}

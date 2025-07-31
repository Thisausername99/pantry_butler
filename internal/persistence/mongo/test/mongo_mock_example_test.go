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
	mockCursor := mocks.NewMockMongoCursor(ctrl)

	// Create repository with mocked collection
	repo := &mongo.PantryEntryRepo{
		Collection: mockCollection,
		Logger:     zap.NewNop(),
	}

	// Set up expectations for GetPantryEntries
	ctx := context.Background()
	pantryID := "test-pantry-id"
	filter := bson.M{"pantryId": pantryID}

	// Mock the Find operation
	mockCollection.EXPECT().
		Find(ctx, filter, gomock.Any()).
		Return(mockCursor, nil).
		Times(1)

	// Mock cursor behavior
	mockCursor.EXPECT().Next(ctx).Return(true).Times(2)  // Two documents
	mockCursor.EXPECT().Next(ctx).Return(false).Times(1) // End of cursor
	mockCursor.EXPECT().Decode(gomock.Any()).DoAndReturn(func(val interface{}) error {
		// Simulate decoding the first document
		entry := val.(*entity.PantryEntry)
		entry.ID = "1"
		entry.Name = "Apples"
		entry.Quantity = intPtr(5)
		return nil
	}).Times(1)
	mockCursor.EXPECT().Decode(gomock.Any()).DoAndReturn(func(val interface{}) error {
		// Simulate decoding the second document
		entry := val.(*entity.PantryEntry)
		entry.ID = "2"
		entry.Name = "Milk"
		entry.Quantity = intPtr(1)
		return nil
	}).Times(1)
	mockCursor.EXPECT().Close(ctx).Return(nil).Times(1)
	mockCursor.EXPECT().Err().Return(nil).Times(1)

	// Execute test
	result, err := repo.GetPantryEntries(ctx, pantryID)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Apples", result[0].Name)
	assert.Equal(t, "Milk", result[1].Name)
}

func TestPantryEntryRepo_InsertEntry_WithMockedMongo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock MongoDB components
	mockCollection := mocks.NewMockMongoCollection(ctrl)
	mockInsertOneResult := mocks.NewMockMongoInsertOneResult(ctrl)

	// Create repository with mocked collection
	repo := &mongo.PantryEntryRepo{
		Collection: mockCollection,
		Logger:     zap.NewNop(),
	}

	// Test data
	entry := &entity.PantryEntryInput{
		Name:         "Bananas",
		Quantity:     intPtr(3),
		QuantityType: stringPtr("bunch"),
	}

	// Set up expectations for InsertPantryEntry
	ctx := context.Background()
	pantryID := "test-pantry-id"
	document := bson.M{"pantryId": pantryID, "entry": entry}

	// Mock the InsertOne operation
	mockCollection.EXPECT().
		InsertOne(ctx, document).
		Return(mockInsertOneResult, nil).
		Times(1)

	mockInsertOneResult.EXPECT().
		InsertedID().
		Return("new-entry-id").
		Times(1)

	// Execute test
	err := repo.InsertPantryEntry(ctx, pantryID, entry)

	// Assertions
	assert.NoError(t, err)
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

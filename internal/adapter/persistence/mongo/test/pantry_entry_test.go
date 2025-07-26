package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	mongorepo "github.com/thisausername99/pantry-butler/internal/adapter/persistence/mongo"
	entity "github.com/thisausername99/pantry-butler/internal/domain"
)

var testClient *mongo.Client
var testDB *mongo.Database

func TestMain(m *testing.M) {
	// Setup test database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use test database
	mongoURI := "mongodb://root:password@localhost:27017/pantry_butler_test?authSource=admin"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	testClient = client
	testDB = client.Database("pantry_butler_test")

	// Run tests
	code := m.Run()

	// Cleanup
	if err := client.Disconnect(context.Background()); err != nil {
		panic(err)
	}

	os.Exit(code)
}

func setupTestCollection(t *testing.T) *mongo.Collection {
	collection := testDB.Collection("pantry_entries")
	// Clear the collection before each test
	_, err := collection.DeleteMany(context.Background(), map[string]interface{}{})
	require.NoError(t, err)
	return collection
}

func TestPantryEntryRepo_GetPantryEntries(t *testing.T) {
	collection := setupTestCollection(t)
	logger := zap.NewNop()
	repo := &mongorepo.PantryEntryRepo{
		Collection: collection,
		Logger:     logger,
	}

	// Insert test data
	testEntries := []interface{}{
		map[string]interface{}{
			"name":       "Apples",
			"quantity":   5,
			"expiration": time.Now().AddDate(0, 0, 7),
			"created_at": time.Now(),
			"updated_at": time.Now(),
		},
		map[string]interface{}{
			"name":       "Milk",
			"quantity":   1,
			"expiration": time.Now().AddDate(0, 0, 3),
			"created_at": time.Now(),
			"updated_at": time.Now(),
		},
	}

	_, err := collection.InsertMany(context.Background(), testEntries)
	require.NoError(t, err)

	// Test retrieval
	ctx := context.Background()
	entries, err := repo.GetPantryEntries(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, "Apples", entries[0].Name)
	assert.Equal(t, "Milk", entries[1].Name)
}

func TestPantryEntryRepo_InsertPantryEntry(t *testing.T) {
	collection := setupTestCollection(t)
	logger := zap.NewNop()
	repo := &mongorepo.PantryEntryRepo{
		Collection: collection,
		Logger:     logger,
	}

	// Test data
	quantity := 5
	expiration := time.Now().AddDate(0, 0, 7)
	input := &entity.PantryEntryInput{
		Name:       "Test Apples",
		Quantity:   &quantity,
		Expiration: &expiration,
	}

	// Test insertion
	ctx := context.Background()
	entry, err := repo.InsertPantryEntry(ctx, input)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, "Test Apples", entry.Name)
	assert.Equal(t, &quantity, entry.Quantity)
	assert.Equal(t, &expiration, entry.Expiration)

	// Verify it was actually inserted
	count, err := collection.CountDocuments(ctx, map[string]interface{}{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestPantryEntryRepo_InsertPantryEntry_EmptyName(t *testing.T) {
	collection := setupTestCollection(t)
	logger := zap.NewNop()
	repo := &mongorepo.PantryEntryRepo{
		Collection: collection,
		Logger:     logger,
	}

	// Test data with empty name
	input := &entity.PantryEntryInput{
		Name: "",
	}

	// Test insertion
	ctx := context.Background()
	entry, err := repo.InsertPantryEntry(ctx, input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, entry)
	assert.Contains(t, err.Error(), "entry does not have a name")
}

func TestPantryEntryRepo_GetPantryEntries_Empty(t *testing.T) {
	collection := setupTestCollection(t)
	logger := zap.NewNop()
	repo := &mongorepo.PantryEntryRepo{
		Collection: collection,
		Logger:     logger,
	}

	// Test retrieval from empty collection
	ctx := context.Background()
	entries, err := repo.GetPantryEntries(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, entries, 0)
}

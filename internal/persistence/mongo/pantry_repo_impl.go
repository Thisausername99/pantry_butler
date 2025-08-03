package mongo

import (
	"context"
	"errors"

	"github.com/thisausername99/pantry_butler/internal/domain/entity"
	"github.com/thisausername99/pantry_butler/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type PantryEntryRepo struct {
	Collection MongoCollection
	Logger     *zap.Logger
}

// Ensure it implements the interface
var _ repository.PantryRepository = (*PantryEntryRepo)(nil)

func (m *PantryEntryRepo) GetPantryEntries(ctx context.Context, pantryID string) ([]entity.PantryEntry, error) {
	// var entries []*entity.PantryEntry
	var pantry entity.Pantry
	m.Logger.Info("PantryEntryRepo.GetPantryEntries: Getting pantry entries", zap.String("pantryID", pantryID))
	err := m.Collection.FindOne(ctx, bson.M{"id": pantryID}).Decode(&pantry)
	if err != nil {
		if m.Logger != nil {
			m.Logger.Error("PantryEntryRepo.GetPantryEntries: Failed to find pantry", zap.Error(err))
		}
		return nil, err
	}
	// m.Logger.Info("Pantry object", zap.Any("pantry", pantry))
	if m.Logger != nil {
		if pantry.Entries != nil {

			m.Logger.Info("PantryEntryRepo.GetPantryEntries: Successfully retrieved pantry entries", zap.Int("count", len(*pantry.Entries)))
			m.Logger.Info("Pantry entries", zap.Any("entries", pantry.Entries))
		} else {
			m.Logger.Info("Successfully retrieved pantry entries", zap.Int("count", 0))
			m.Logger.Info("Pantry entries", zap.Any("entries", nil))
		}
	}

	// Handle nil entries
	if pantry.Entries == nil {
		return []entity.PantryEntry{}, nil
	}

	return *pantry.Entries, nil
}

func (m *PantryEntryRepo) InsertPantryEntry(ctx context.Context, pantryID string, entry *entity.PantryEntry) error {
	if entry.Name == "" {
		return errors.New("entry does not have a name")
	}

	//! TODO: Add expiration date or generate one

	// Add the entry to the pantry's pantry_entries array where id == pantryID
	filter := bson.M{"id": pantryID}
	update := bson.M{
		"$push": bson.M{
			"pantry_entries": entry,
		},
	}
	_, err := m.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		m.Logger.Error("Failed to insert pantry entry", zap.Error(err))
		return err
	}
	m.Logger.Info("Inserted pantry entry", zap.Any("entry", entry))
	return nil
}

func (m *PantryEntryRepo) DeletePantryEntry(ctx context.Context, pantryID string, entryID string) error {
	_, err := m.Collection.DeleteOne(ctx, bson.M{"pantryId": pantryID, "entryId": entryID})
	if err != nil {
		m.Logger.Error("Failed to delete pantry entry", zap.Error(err))
		return err
	}
	m.Logger.Info("Deleted pantry entry", zap.String("entryId", entryID))
	return nil
}

func (m *PantryEntryRepo) CreateNewPantry(ctx context.Context, pantry *entity.Pantry) error {
	_, err := m.Collection.InsertOne(ctx, pantry)
	if err != nil {
		if m.Logger != nil {
			m.Logger.Error("Failed to create new pantry", zap.Error(err))
		}
		return err
	}
	return nil
}

func (m *PantryEntryRepo) DeletePantry(ctx context.Context, pantryID string) error {
	_, err := m.Collection.DeleteOne(ctx, bson.M{"id": pantryID})
	m.Logger.Info("PantryEntryRepo.DeletePantry: Deleting pantry", zap.String("pantryID", pantryID))
	if err != nil {
		m.Logger.Error("Failed to delete pantry", zap.Error(err))
		return err
	}
	return nil
}

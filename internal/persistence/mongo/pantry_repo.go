package mongo

import (
	"context"
	"errors"

	"github.com/thisausername99/pantry_butler/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type PantryEntryRepo struct {
	Collection MongoCollection
	Logger     *zap.Logger
}

func (m *PantryEntryRepo) GetPantryEntries(ctx context.Context, pantryID string) ([]*entity.PantryEntry, error) {
	var entries []*entity.PantryEntry
	cursor, err := m.Collection.Find(ctx, bson.M{"pantryId": pantryID})
	if err != nil {
		if m.Logger != nil {
			m.Logger.Error("Failed to find pantry entries", zap.Error(err))
		}
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var entry entity.PantryEntry
		if err := cursor.Decode(&entry); err != nil {
			if m.Logger != nil {
				m.Logger.Error("Failed to decode pantry entry", zap.Error(err))
			}
			return nil, err
		}
		entries = append(entries, &entry)
	}
	if err := cursor.Err(); err != nil {
		if m.Logger != nil {
			m.Logger.Error("Cursor error after iterating pantry entries", zap.Error(err))
		}
		return nil, err
	}
	if m.Logger != nil {
		m.Logger.Info("Successfully retrieved pantry entries", zap.Int("count", len(entries)))
	}
	return entries, nil
}

func (m *PantryEntryRepo) InsertPantryEntry(ctx context.Context, pantryID string, entry *entity.PantryEntryInput) error {
	pantryEntry := &entity.PantryEntry{}
	if entry.Name == "" {
		return errors.New("entry does not have a name")
	}

	pantryEntry.Name = entry.Name

	if entry.Quantity != nil {
		pantryEntry.Quantity = entry.Quantity
	}

	//! TODO: Add expiration date or generate one

	_, err := m.Collection.InsertOne(ctx, bson.M{"pantryId": pantryID, "entry": pantryEntry})
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

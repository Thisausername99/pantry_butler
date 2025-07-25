package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	entity "github.com/thisausername99/pantry-butler/internal/domain"
)

type PantryEntryRepo struct {
	Collection *mongo.Collection
	Logger     *zap.Logger
}

func (m *PantryEntryRepo) GetPantryEntries(ctx context.Context) ([]*entity.PantryEntry, error) {
	var entries []*entity.PantryEntry
	cursor, err := m.Collection.Find(ctx, bson.M{})
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

func (m *PantryEntryRepo) InsertPantryEntry(ctx context.Context, entry *entity.PantryEntryInput) (*entity.PantryEntry, error) {
	pantryEntry := &entity.PantryEntry{}
	if entry.Name == "" {
		return nil, errors.New("entry does not have a name")
	}

	pantryEntry.Name = entry.Name

	if entry.Quantity != nil {
		pantryEntry.Quantity = entry.Quantity
	}

	if entry.Expiration != nil {
		pantryEntry.Expiration = entry.Expiration
	}

	_, err := m.Collection.InsertOne(ctx, entry)
	if err != nil {
		m.Logger.Error("Failed to insert pantry entry", zap.Error(err))
		return nil, err
	}
	m.Logger.Info("Inserted pantry entry", zap.Any("entry", entry))
	return pantryEntry, nil
}

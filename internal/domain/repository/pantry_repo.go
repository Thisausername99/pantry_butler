package repository

import (
	"context"

	"github.com/thisausername99/pantry_butler/internal/domain/entity"
)

type PantryRepository interface {
	GetPantryEntries(ctx context.Context, pantryID string) ([]entity.PantryEntry, error)
	InsertPantryEntry(ctx context.Context, pantryID string, entry *entity.PantryEntry) error
	DeletePantryEntry(ctx context.Context, pantryID string, entryID string) error
	CreateNewPantry(ctx context.Context, pantry *entity.Pantry) error
	DeletePantry(ctx context.Context, pantryID string) error
}

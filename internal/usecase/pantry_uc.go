package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/thisausername99/pantry_butler/internal/domain/entity"
)

func (u *Usecase) GetAllPantryEntries(ctx context.Context, pantryID string) ([]entity.PantryEntry, error) {
	return u.RepoWrapper.PantryRepo.GetPantryEntries(ctx, pantryID)
}

func (u *Usecase) InsertPantryEntry(ctx context.Context, pantryID string, pantryEntryInput *entity.PantryEntryInput) error {
	entry := &entity.PantryEntry{
		ID:           uuid.New().String(),
		Name:         pantryEntryInput.Name,
		Quantity:     pantryEntryInput.Quantity,
		QuantityType: pantryEntryInput.QuantityType,
	}

	return u.RepoWrapper.PantryRepo.InsertPantryEntry(ctx, pantryID, entry)
}

func (u *Usecase) DeletePantryEntry(ctx context.Context, pantryID string, entryID string) error {
	return u.RepoWrapper.PantryRepo.DeletePantryEntry(ctx, pantryID, entryID)
}

func (u *Usecase) CreateNewPantry(ctx context.Context, pantry *entity.Pantry) error {
	return u.RepoWrapper.PantryRepo.CreateNewPantry(ctx, pantry)
}

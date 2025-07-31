package usecase

import (
	"context"

	"github.com/thisausername99/pantry_butler/internal/domain/entity"
)

func (u *Usecase) GetAllPantryEntries(ctx context.Context, pantryID string) ([]*entity.PantryEntry, error) {
	return u.RepoWrapper.PantryRepo.GetPantryEntries(ctx, pantryID)
}

func (u *Usecase) InsertPantryEntry(ctx context.Context, pantryID string, entry *entity.PantryEntryInput) error {
	return u.RepoWrapper.PantryRepo.InsertPantryEntry(ctx, pantryID, entry)
}

func (u *Usecase) DeletePantryEntry(ctx context.Context, pantryID string, entryID string) error {
	return u.RepoWrapper.PantryRepo.DeletePantryEntry(ctx, pantryID, entryID)
}

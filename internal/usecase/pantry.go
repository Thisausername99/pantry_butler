package usecase

import (
	"context"

	entity "github.com/thisausername99/pantry-butler/internal/domain"
)

func (u *Usecase) GetAllPantryEntries(ctx context.Context) ([]*entity.PantryEntry, error) {
	return u.RepoWrapper.PantryEntryRepo.GetPantryEntries(ctx)
}

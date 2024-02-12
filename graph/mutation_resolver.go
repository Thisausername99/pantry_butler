package graph

import (
	"context"
	"errors"

	"github.com/thisausername99/recipes-api/models"
)

type mutationResolver struct{ *Resolver }

// UpsertEntry is the resolver for the upsertEntry field.
func (r *mutationResolver) InsertEntry(ctx context.Context, entry models.PantryEntryInput) (*models.PantryEntry, error) {
	pantryEntry := &models.PantryEntry{}
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

	return r.PantryEntryRepo.InsertPantryEntry(pantryEntry)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

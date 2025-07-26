//go:generate mockgen -destination=test/mocks.go -package=test github.com/thisausername99/pantry-butler/internal/usecase RecipeRepository,PantryEntryRepository

package usecase

// This file is used to generate mocks for testing 
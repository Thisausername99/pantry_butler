package usecase

import (
	repo "github.com/thisausername99/pantry_butler/internal/domain/repository"
	"go.uber.org/zap"
)

type RepoWrapper struct {
	RecipeRepo repo.RecipeRepository
	PantryRepo repo.PantryRepository
	UserRepo   repo.UserRepository
	// Add more repositories as needed
}

type Usecase struct {
	RepoWrapper RepoWrapper
	Logger      *zap.Logger
}

func NewUsecase(repoWrapper RepoWrapper, logger *zap.Logger) *Usecase {
	return &Usecase{
		RepoWrapper: repoWrapper,
		Logger:      logger,
	}
}

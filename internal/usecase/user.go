package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/thisausername99/pantry_butler/internal/domain/entity"
	"github.com/thisausername99/pantry_butler/pkg/security"
	"go.uber.org/zap"
)

func (u *Usecase) GetUser(ctx context.Context, id string) (*entity.User, error) {
	return u.RepoWrapper.UserRepo.GetUser(ctx, id)
}

func (u *Usecase) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	// Check if user already exists
	existingUser, err := u.RepoWrapper.UserRepo.GetUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		u.Logger.Error("user already exists", zap.String("id", user.ID))
		return nil, errors.New("user already exists")
	}
	// Create user
	user.ID = uuid.New().String()
	hashedPassword, err := security.HashPassword(user.Password)
	if err != nil {
		u.Logger.Error("error hashing password", zap.Error(err))
	}
	user.Password = hashedPassword
	user.CreatedAt = time.Now()

	return u.RepoWrapper.UserRepo.CreateUser(ctx, user)
}

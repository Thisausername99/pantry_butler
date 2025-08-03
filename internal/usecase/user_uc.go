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

func (u *Usecase) RegisterUser(ctx context.Context, input *entity.UserRegisterInput) (*entity.User, error) {
	// Check if user already exists by email
	existingUser, err := u.RepoWrapper.UserRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		u.Logger.Error("error checking existing user", zap.Error(err))
		return nil, err
	}
	if existingUser != nil {
		u.Logger.Error("user already exists", zap.String("email", input.Email))
		return nil, errors.New("user already exists")
	}

	// Create user object with business logic
	user := &entity.User{
		ID:        uuid.New().String(),
		Email:     input.Email,
		CreatedAt: time.Now(),
		Pantries:  []string{}, // Initialize empty pantries array
	}

	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}

	// Hash password
	hashedPassword, err := security.HashPassword(input.Password)
	if err != nil {
		u.Logger.Error("error hashing password", zap.Error(err))
		return nil, err
	}
	user.Password = hashedPassword

	// Persist user to database
	err = u.RepoWrapper.UserRepo.CreateUser(ctx, user)
	if err != nil {
		u.Logger.Error("error creating user", zap.Error(err))
		return nil, err
	}

	u.Logger.Info("user created successfully", zap.String("id", user.ID))
	return user, nil
}

func (u *Usecase) UpdateUserWithPantry(ctx context.Context, userID string, name string) error {
	entries := &[]entity.PantryEntry{}
	newPantry := &entity.Pantry{
		ID:        uuid.New().String(),
		Name:      name,
		OwnerID:   userID,
		CreatedAt: time.Now(),
		Entries:   entries,
	}
	err := u.RepoWrapper.PantryRepo.CreateNewPantry(ctx, newPantry)
	if err != nil {
		return err
	}
	err = u.RepoWrapper.UserRepo.UpdateUserWithPantry(ctx, userID, newPantry.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) RemoveUserPantry(ctx context.Context, userID string, pantryID string) error {
	err := u.RepoWrapper.UserRepo.DeletePantryFromUser(ctx, userID, pantryID)
	if err != nil {
		u.Logger.Error("error deleting pantry from user", zap.Error(err))
		return err
	}
	err = u.RepoWrapper.PantryRepo.DeletePantry(ctx, pantryID)
	if err != nil {
		u.Logger.Error("error deleting pantry", zap.Error(err))
		return err
	}
	return nil
}

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

func (u *Usecase) RegisterUser(ctx context.Context, userInput *entity.UserInput) error {
	// Check if user already exists
	newUser := &entity.User{
		ID:        uuid.New().String(),
		Email:     user.Email,
		UserName:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: time.Now(),
	}
	existingUser, err := u.RepoWrapper.UserRepo.GetUserByEmail(ctx, newUser.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		u.Logger.Error("user already exists", zap.String("id", userInput.Email))
		return errors.New("user already exists")
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

func (u *Usecase) UpdateUserWithPantry(ctx context.Context, userID string, name string) error {
	entries := &[]entity.PantryEntry{}
	newPantry := &entity.Pantry{
		ID:        uuid.New().String(),
		Name:      name,
		Owner:     userID,
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

func (u *Usecase) RegisterUser(ctx context.Context, input *entity.UserInput) (*entity.User, error) {
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
		FirstName: input.FirstName,
		LastName:  input.LastName,
		CreatedAt: time.Now(),
		Pantries:  []string{}, // Initialize empty pantries array
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

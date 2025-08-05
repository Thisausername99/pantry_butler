package repository

import (
	"context"

	"github.com/thisausername99/pantry_butler/internal/domain/entity"
)

type UserRepository interface {
	// User CRUD operations
	CreateUser(ctx context.Context, user *entity.User) error
	GetUser(ctx context.Context, userID string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateUser(ctx context.Context, userID string, user *entity.User) error
	DeleteUser(ctx context.Context, userID string) error

	// User authentication & security
	UpdatePassword(ctx context.Context, userID string, hashedPassword string) error
	GetUserByCredentials(ctx context.Context, email, hashedPassword string) (*entity.User, error)

	// User profile management
	UpdateProfile(ctx context.Context, userID string, firstName, lastName string) error
	UpdateEmail(ctx context.Context, userID string, email string) error

	// User pantry management
	UpdateUserWithPantry(ctx context.Context, userID string, pantryID string) error
	DeletePantryFromUser(ctx context.Context, userID string, pantryID string) error
	GetUserPantries(ctx context.Context, userID string) ([]string, error)

	// User search & listing
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
	// GetUsersByRole(ctx context.Context, role string) ([]*entity.User, error)
	// SearchUsers(ctx context.Context, query string) ([]*entity.User, error)

	// User status & verification
	// UpdateUserStatus(ctx context.Context, userID string, status string) error
	// VerifyUserEmail(ctx context.Context, userID string) error
	// SetEmailVerificationToken(ctx context.Context, userID string, token string) error

	// User preferences & settings
	// UpdateUserPreferences(ctx context.Context, userID string, preferences map[string]interface{}) error
	// GetUserPreferences(ctx context.Context, userID string) (map[string]interface{}, error)
}

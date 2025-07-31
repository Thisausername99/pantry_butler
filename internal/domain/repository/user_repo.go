package repository

import (
	"context"

	"github.com/thisausername99/pantry_butler/internal/domain/entity"
)

type UserRepository interface {
	GetUser(ctx context.Context, userID string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
}

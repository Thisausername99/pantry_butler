package config

import (
	"time"

	"github.com/thisausername99/pantry_butler/internal/usecase"
	"go.uber.org/zap"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret     string
	TokenDuration time.Duration
	UseCase       *usecase.Usecase
	Logger        *zap.Logger
}

// NewAuthConfig creates a new AuthConfig with default values
func NewAuthConfig(useCase *usecase.Usecase, logger *zap.Logger) *AuthConfig {
	return &AuthConfig{
		JWTSecret:     getEnv("JWT_SECRET", "default-secret-key-change-in-production"),
		TokenDuration: getDurationEnv("JWT_DURATION", 24*time.Hour),
		UseCase:       useCase,
		Logger:        logger,
	}
}

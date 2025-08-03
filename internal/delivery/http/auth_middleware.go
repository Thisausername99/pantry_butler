package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"github.com/thisausername99/pantry_butler/config"
	"github.com/thisausername99/pantry_butler/internal/delivery/types"
	"github.com/thisausername99/pantry_butler/internal/domain/entity"
)

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware(authConfig *config.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth for public endpoints
		if isPublicEndpoint(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			authConfig.Logger.Warn("Missing authorization header",
				zap.String("requestID", getRequestID(c)),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
				"code":  "MISSING_AUTH_HEADER",
			})
			c.Abort()
			return
		}

		// Parse Bearer token
		tokenString := extractBearerToken(authHeader)
		if tokenString == "" {
			authConfig.Logger.Warn("Invalid authorization header format",
				zap.String("requestID", getRequestID(c)),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
				"code":  "INVALID_AUTH_FORMAT",
			})
			c.Abort()
			return
		}

		// Validate and parse JWT token
		claims, err := validateJWTToken(tokenString, authConfig.JWTSecret)
		if err != nil {
			authConfig.Logger.Warn("Invalid JWT token",
				zap.String("requestID", getRequestID(c)),
				zap.String("path", c.Request.URL.Path),
				zap.Error(err),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
				"code":  "INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		// Fetch user from database to ensure they still exist
		user, err := authConfig.UseCase.RepoWrapper.UserRepo.GetUser(context.Background(), claims.UserID)
		if err != nil {
			authConfig.Logger.Warn("User not found in database",
				zap.String("requestID", getRequestID(c)),
				zap.String("userID", claims.UserID),
				zap.Error(err),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
				"code":  "USER_NOT_FOUND",
			})
			c.Abort()
			return
		}

		// Create current user context
		currentUser := &types.CurrentUser{
			ID:       user.ID,
			UserName: &user.UserName,
			Email:    user.Email,
		}

		// Add user to context
		ctx := context.WithValue(c.Request.Context(), types.CurrentUserKey, currentUser)
		ctx = context.WithValue(ctx, types.UserIDKey, user.ID)
		ctx = context.WithValue(ctx, types.AuthenticatedKey, true)
		c.Request = c.Request.WithContext(ctx)

		// Add user info to Gin context for easy access
		c.Set("currentUser", currentUser)
		c.Set("userID", user.ID)

		authConfig.Logger.Info("User authenticated",
			zap.String("requestID", getRequestID(c)),
			zap.String("userID", user.ID),
			zap.String("userName", user.UserName),
			zap.String("path", c.Request.URL.Path),
		)

		c.Next()
	}
}

// OptionalAuthMiddleware allows optional authentication
func OptionalAuthMiddleware(authConfig *config.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No token provided, continue without authentication
			c.Next()
			return
		}

		tokenString := extractBearerToken(authHeader)
		if tokenString == "" {
			// Invalid token format, continue without authentication
			c.Next()
			return
		}

		claims, err := validateJWTToken(tokenString, authConfig.JWTSecret)
		if err != nil {
			// Invalid token, continue without authentication
			c.Next()
			return
		}

		// Token is valid, fetch user and set context
		user, err := authConfig.UseCase.RepoWrapper.UserRepo.GetUser(context.Background(), claims.UserID)
		if err != nil {
			// User not found, continue without authentication
			c.Next()
			return
		}

		currentUser := &types.CurrentUser{
			ID:       user.ID,
			UserName: &user.UserName,
			Email:    user.Email,
		}

		ctx := context.WithValue(c.Request.Context(), types.CurrentUserKey, currentUser)
		ctx = context.WithValue(ctx, types.UserIDKey, user.ID)
		ctx = context.WithValue(ctx, types.AuthenticatedKey, true)
		c.Request = c.Request.WithContext(ctx)

		c.Set("currentUser", currentUser)
		c.Set("userID", user.ID)

		c.Next()
	}
}

// GenerateJWTToken creates a new JWT token for a user
func GenerateJWTToken(user *entity.User, secret string, duration time.Duration) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		UserName: user.UserName,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "pantry-butler",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// validateJWTToken validates and parses a JWT token
func validateJWTToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// extractBearerToken extracts the token from Authorization header
func extractBearerToken(authHeader string) string {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(authHeader, "Bearer ")
}

// isPublicEndpoint checks if the endpoint is public (no auth required)
func isPublicEndpoint(path string) bool {
	publicPaths := []string{
		"/",
		"/health",
		"/api/info",
		"/login",
		"/register",
		"/graphql", // Allow GraphQL playground
	}

	for _, publicPath := range publicPaths {
		if path == publicPath {
			return true
		}
	}

	return false
}

// GetCurrentUser extracts the current user from Gin context
func GetCurrentUser(c *gin.Context) (*types.CurrentUser, bool) {
	if user, exists := c.Get("currentUser"); exists {
		if currentUser, ok := user.(*types.CurrentUser); ok {
			return currentUser, true
		}
	}
	return nil, false
}

// GetCurrentUserFromContext extracts the current user from request context
func GetCurrentUserFromContext(ctx context.Context) (*types.CurrentUser, bool) {
	if user, ok := ctx.Value(types.CurrentUserKey).(*types.CurrentUser); ok {
		return user, true
	}
	return nil, false
}

// IsAuthenticated checks if the request is authenticated
func IsAuthenticated(c *gin.Context) bool {
	if authenticated, exists := c.Get("authenticated"); exists {
		if isAuth, ok := authenticated.(bool); ok {
			return isAuth
		}
	}
	return false
}

// RequireAuth middleware that requires authentication
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsAuthenticated(c) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
				"code":  "AUTH_REQUIRED",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

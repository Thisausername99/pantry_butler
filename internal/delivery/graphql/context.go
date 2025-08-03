package graphql

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/thisausername99/pantry_butler/internal/delivery/types"
)

// GetGraphQLContext extracts GraphQL context from the request context
func GetGraphQLContext(ctx context.Context) *types.GraphQLContext {
	gqlCtx := &types.GraphQLContext{}

	// Extract current user
	if user, ok := ctx.Value(types.CurrentUserKey).(*types.CurrentUser); ok {
		gqlCtx.CurrentUser = user
		gqlCtx.UserID = user.ID
		gqlCtx.IsAuthenticated = true
	}

	// Extract request ID
	if requestID, ok := ctx.Value(types.RequestIDKey).(string); ok {
		gqlCtx.RequestID = requestID
	}

	// Extract Gin context for additional info
	if ginCtx, ok := ctx.Value("ginContext").(*gin.Context); ok {
		gqlCtx.ClientIP = ginCtx.ClientIP()
		gqlCtx.UserAgent = ginCtx.GetHeader("User-Agent")
		gqlCtx.Method = ginCtx.Request.Method
		gqlCtx.Path = ginCtx.Request.URL.Path
	}

	return gqlCtx
}

// GetCurrentUserFromGraphQLContext extracts the current user from GraphQL context
func GetCurrentUserFromGraphQLContext(ctx context.Context) (*types.CurrentUser, bool) {
	gqlCtx := GetGraphQLContext(ctx)
	return gqlCtx.CurrentUser, gqlCtx.IsAuthenticated
}

// IsAuthenticatedInGraphQL checks if the GraphQL request is authenticated
func IsAuthenticatedInGraphQL(ctx context.Context) bool {
	gqlCtx := GetGraphQLContext(ctx)
	return gqlCtx.IsAuthenticated
}

// RequireAuthInGraphQL panics if the request is not authenticated
func RequireAuthInGraphQL(ctx context.Context) *types.CurrentUser {
	user, authenticated := GetCurrentUserFromGraphQLContext(ctx)
	if !authenticated {
		panic("Authentication required")
	}
	return user
}

// GetUserIDFromGraphQLContext extracts the user ID from GraphQL context
func GetUserIDFromGraphQLContext(ctx context.Context) (string, bool) {
	gqlCtx := GetGraphQLContext(ctx)
	if gqlCtx.IsAuthenticated {
		return gqlCtx.UserID, true
	}
	return "", false
}

// GetRequestInfoFromGraphQLContext extracts request information from GraphQL context
func GetRequestInfoFromGraphQLContext(ctx context.Context) map[string]interface{} {
	info := make(map[string]interface{})

	gqlCtx := GetGraphQLContext(ctx)

	if gqlCtx.IsAuthenticated {
		info["userID"] = gqlCtx.UserID
		info["userName"] = gqlCtx.CurrentUser.UserName
		info["email"] = gqlCtx.CurrentUser.Email
	}

	info["requestID"] = gqlCtx.RequestID
	info["authenticated"] = gqlCtx.IsAuthenticated
	info["clientIP"] = gqlCtx.ClientIP
	info["userAgent"] = gqlCtx.UserAgent
	info["method"] = gqlCtx.Method
	info["path"] = gqlCtx.Path

	return info
}

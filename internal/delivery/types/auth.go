package types

// Context key types to avoid linter warnings
type contextKey string

const (
	CurrentUserKey   contextKey = "currentUser"
	UserIDKey        contextKey = "userID"
	AuthenticatedKey contextKey = "authenticated"
	RequestIDKey     contextKey = "requestID"
)

// CurrentUser represents the current authenticated user
type CurrentUser struct {
	ID       string  `json:"id"`
	UserName *string `json:"userName"`
	Email    string  `json:"email"`
}

// GraphQLContext represents the context available in GraphQL resolvers
type GraphQLContext struct {
	CurrentUser     *CurrentUser
	UserID          string
	IsAuthenticated bool
	RequestID       string
	ClientIP        string
	UserAgent       string
	Method          string
	Path            string
}

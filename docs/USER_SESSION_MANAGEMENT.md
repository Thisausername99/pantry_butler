# User Session Management Guide

This guide explains how to identify and manage the current user session in your Pantry Butler application.

## 🔍 How to Identify the Current User

### 1. **In HTTP Handlers (Gin)**

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/thisausername99/pantry_butler/internal/delivery/http"
    "github.com/thisausername99/pantry_butler/internal/delivery/types"
)

func getUserProfile(c *gin.Context) {
    // Method 1: Get current user from Gin context
    currentUser, exists := http.GetCurrentUser(c)
    if !exists {
        c.JSON(401, gin.H{"error": "User not authenticated"})
        return
    }
    
    // Access user information
    userID := currentUser.ID
    userName := currentUser.UserName
    email := currentUser.Email
    
    c.JSON(200, gin.H{
        "user": currentUser,
        "message": "Profile retrieved successfully",
    })
}

func checkAuthentication(c *gin.Context) {
    // Method 2: Check if user is authenticated
    if !http.IsAuthenticated(c) {
        c.JSON(401, gin.H{"error": "Authentication required"})
        return
    }
    
    // User is authenticated, proceed
    c.Next()
}
```

### 2. **In GraphQL Resolvers**

```go
package graphql

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"
    "github.com/thisausername99/pantry_butler/internal/delivery/graphql"
    "github.com/thisausername99/pantry_butler/internal/delivery/types"
    "github.com/thisausername99/pantry_butler/internal/entity"
)

// Example resolver that requires authentication
func (r *mutationResolver) CreatePantry(ctx context.Context, input CreatePantryInput) (*Pantry, error) {
    // Method 1: Get current user from GraphQL context
    currentUser, authenticated := graphql.GetCurrentUserFromGraphQLContext(ctx)
    if !authenticated {
        return nil, fmt.Errorf("authentication required")
    }
    
    // Use current user information
    userID := currentUser.ID
    userName := currentUser.UserName
    
    // Create pantry with user association
    pantry := &entity.Pantry{
        Name: input.Name,
        UserID: userID,
        // ... other fields
    }
    
    return r.UseCase.RepoWrapper.PantryRepo.CreatePantry(ctx, pantry)
}

// Example resolver with optional authentication
func (r *queryResolver) GetPantries(ctx context.Context) ([]*Pantry, error) {
    // Method 2: Check if user is authenticated
    if graphql.IsAuthenticatedInGraphQL(ctx) {
        // User is authenticated, return their pantries
        userID, _ := graphql.GetUserIDFromGraphQLContext(ctx)
        return r.UseCase.RepoWrapper.PantryRepo.GetPantriesByUserID(ctx, userID)
    }
    
    // User is not authenticated, return public pantries
    return r.UseCase.RepoWrapper.PantryRepo.GetPublicPantries(ctx)
}

// Example resolver that requires authentication (with panic)
func (r *mutationResolver) UpdateUserProfile(ctx context.Context, input UpdateUserInput) (*User, error) {
    // Method 3: Require authentication (panics if not authenticated)
    currentUser := graphql.RequireAuthInGraphQL(ctx)
    
    // Update user profile
    user := &entity.User{
        ID: currentUser.ID,
        FirstName: input.FirstName,
        LastName: input.LastName,
        // ... other fields
    }
    
    return r.UseCase.RepoWrapper.UserRepo.UpdateUser(ctx, currentUser.ID, user)
}
```

### 3. **In Use Case Layer**

```go
package usecase

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"
    "github.com/thisausername99/pantry_butler/internal/delivery/types"
    "github.com/thisausername99/pantry_butler/internal/entity"
)

type PantryUseCase struct {
    // ... dependencies
}

func (uc *PantryUseCase) CreatePantry(ctx context.Context, name string) (*entity.Pantry, error) {
    // Get current user from context
    currentUser, ok := ctx.Value("currentUser").(*types.CurrentUser)
    if !ok {
        return nil, fmt.Errorf("user not authenticated")
    }
    
    // Create pantry with user association
    pantry := &entity.Pantry{
        Name: name,
        UserID: currentUser.ID,
        CreatedAt: time.Now(),
    }
    
    return uc.pantryRepo.CreatePantry(ctx, pantry)
}

func (uc *PantryUseCase) GetUserPantries(ctx context.Context) ([]*entity.Pantry, error) {
    // Get user ID from context
    userID, ok := ctx.Value("userID").(string)
    if !ok {
        return nil, fmt.Errorf("user not authenticated")
    }
    
    return uc.pantryRepo.GetPantriesByUserID(ctx, userID)
}
```

### 4. **In Middleware**

```go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/thisausername99/pantry_butler/internal/delivery/http"
    "log"
)

// Custom middleware that logs user actions
func UserActionLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get current user
        currentUser, exists := http.GetCurrentUser(c)
        
        if exists {
            // Log authenticated user action
            log.Printf("User %s (%s) accessed %s", 
                currentUser.UserName, 
                currentUser.ID, 
                c.Request.URL.Path)
        } else {
            // Log anonymous user action
            log.Printf("Anonymous user accessed %s", c.Request.URL.Path)
        }
        
        c.Next()
    }
}

// Middleware that requires specific user role
func RequireAdmin() gin.HandlerFunc {
    return func(c *gin.Context) {
        currentUser, exists := http.GetCurrentUser(c)
        if !exists {
            c.JSON(401, gin.H{"error": "Authentication required"})
            c.Abort()
            return
        }
        
        // Check if user is admin (you can add role field to user entity)
        if currentUser.UserName != "admin" {
            c.JSON(403, gin.H{"error": "Admin access required"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

## 🔧 Setting Up Authentication

### 1. **Configure JWT Secret**

Add to your environment variables:
```bash
JWT_SECRET=your-super-secret-jwt-key-here
JWT_DURATION=24h
```

### 2. **Update Server Configuration**

```go
// In your main.go or server setup
func setupAuthMiddleware(useCase *usecase.Usecase, logger *zap.Logger) gin.HandlerFunc {
    config := http.AuthConfig{
        JWTSecret:     os.Getenv("JWT_SECRET"),
        TokenDuration: 24 * time.Hour, // or parse from env
        UseCase:       useCase,
        Logger:        logger,
    }
    
    return http.AuthMiddleware(config)
}

// Apply middleware in your server setup
func (s *Server) setupMiddleware() {
    // ... other middleware
    
    // Apply authentication middleware
    s.router.Use(setupAuthMiddleware(s.useCase, s.logger))
}
```

### 3. **Create Login Endpoint**

```go
func loginHandler(c *gin.Context) {
    var loginRequest struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&loginRequest); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }
    
    // Authenticate user (implement your logic)
    user, err := authenticateUser(loginRequest.Email, loginRequest.Password)
    if err != nil {
        c.JSON(401, gin.H{"error": "Invalid credentials"})
        return
    }
    
    // Generate JWT token
    token, err := http.GenerateJWTToken(user, os.Getenv("JWT_SECRET"), 24*time.Hour)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to generate token"})
        return
    }
    
    c.JSON(200, gin.H{
        "token": token,
        "user": gin.H{
            "id": user.ID,
            "userName": user.UserName,
            "email": user.Email,
        },
    })
}
```

## 📊 Request Information Available

### **Current User Information**
```go
type CurrentUser struct {
    ID        string    `json:"id"`
    UserName  string    `json:"userName"`
    Email     string    `json:"email"`
    FirstName string    `json:"firstName"`
    LastName  string    `json:"lastName"`
    CreatedAt time.Time `json:"createdAt"`
    Pantries  []string  `json:"pantries"` // Associated pantry IDs
}
```

### **Request Context Information**
```go
// Available in GraphQL context
gqlCtx := graphql.GetGraphQLContext(ctx)
info := map[string]interface{}{
    "userID": gqlCtx.UserID,
    "userName": gqlCtx.CurrentUser.UserName,
    "email": gqlCtx.CurrentUser.Email,
    "requestID": gqlCtx.RequestID,
    "authenticated": gqlCtx.IsAuthenticated,
    "clientIP": gqlCtx.ClientIP,
    "userAgent": gqlCtx.UserAgent,
    "method": gqlCtx.Method,
    "path": gqlCtx.Path,
}
```

## 🛡️ Security Best Practices

### 1. **Always Validate Authentication**
```go
// ✅ Good: Check authentication before accessing user data
func getPrivateData(c *gin.Context) {
    currentUser, exists := http.GetCurrentUser(c)
    if !exists {
        c.JSON(401, gin.H{"error": "Authentication required"})
        return
    }
    
    // Only return data for authenticated user
    data := getDataForUser(currentUser.ID)
    c.JSON(200, data)
}

// ❌ Bad: Assume user is authenticated
func getPrivateData(c *gin.Context) {
    userID := c.GetHeader("X-User-ID") // Can be spoofed!
    data := getDataForUser(userID)
    c.JSON(200, data)
}
```

### 2. **Use Context for User Information**
```go
// ✅ Good: Get user from context (validated by middleware)
currentUser, _ := http.GetCurrentUser(c)
userID := currentUser.ID

// ❌ Bad: Get user from headers (can be spoofed)
userID := c.GetHeader("X-User-ID")
```

### 3. **Validate User Permissions**
```go
func updatePantry(c *gin.Context) {
    currentUser, _ := http.GetCurrentUser(c)
    pantryID := c.Param("id")
    
    // Check if user owns this pantry
    pantry, err := getPantry(pantryID)
    if err != nil || pantry.UserID != currentUser.ID {
        c.JSON(403, gin.H{"error": "Access denied"})
        return
    }
    
    // Proceed with update
    updatePantryData(pantryID, c.Request.Body)
}
```

## 🔍 Debugging User Sessions

### 1. **Check Authentication Status**
```go
func debugAuth(c *gin.Context) {
    currentUser, exists := http.GetCurrentUser(c)
    
    debugInfo := gin.H{
        "authenticated": exists,
        "requestID": getRequestID(c),
        "clientIP": c.ClientIP(),
        "userAgent": c.GetHeader("User-Agent"),
    }
    
    if exists {
        debugInfo["user"] = gin.H{
            "id": currentUser.ID,
            "userName": currentUser.UserName,
            "email": currentUser.Email,
        }
    }
    
    c.JSON(200, debugInfo)
}
```

### 2. **Log User Actions**
```go
func logUserAction(c *gin.Context) {
    currentUser, exists := http.GetCurrentUser(c)
    
    if exists {
        log.Printf("[AUTH] User %s (%s) performed action: %s %s",
            currentUser.UserName,
            currentUser.ID,
            c.Request.Method,
            c.Request.URL.Path)
    } else {
        log.Printf("[AUTH] Anonymous user performed action: %s %s",
            c.Request.Method,
            c.Request.URL.Path)
    }
    
    c.Next()
}
```

## 📝 Example Usage in Your Application

### **GraphQL Query with User Context**
```graphql
query GetMyPantries {
  pantries {
    id
    name
    items {
      id
      name
      quantity
    }
  }
}
```

### **HTTP Request with Authentication**
```bash
# Login to get token
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'

# Use token for authenticated requests
curl -X GET http://localhost:8080/api/pantries \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

This comprehensive system allows you to identify the current user throughout your application while maintaining security and providing flexibility for both authenticated and anonymous users. 
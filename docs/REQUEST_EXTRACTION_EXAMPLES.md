# Request Extraction Examples for Pantry Butler

This document shows practical examples of how to extract and handle requests from the frontend in your current GraphQL setup.

## 🎯 Current Request Flow

```
Frontend → GraphQL Query/Mutation → Resolver → UseCase → Repository → MongoDB
```

## 📝 Basic Request Extraction

### **1. Simple Query Request**
```javascript
// Frontend sends:
query {
  getRecipe {
    id
    name
    cuisine
  }
}
```

```go
// Resolver automatically receives:
func (r *queryResolver) GetRecipe(ctx context.Context) ([]*entity.Recipe, error) {
    // No parameters to extract - just return all recipes
    return r.UseCase.RepoWrapper.RecipeRepo.GetRecipes(ctx)
}
```

### **2. Query with Parameters**
```javascript
// Frontend sends:
query {
  getRecipeByCuisine(cuisine: "Italian") {
    id
    name
    cuisine
  }
}
```

```go
// Resolver automatically receives:
func (r *queryResolver) GetRecipeByCuisine(ctx context.Context, cuisine string) ([]*entity.Recipe, error) {
    // 'cuisine' parameter is automatically extracted
    // cuisine = "Italian"
    return r.UseCase.RepoWrapper.RecipeRepo.GetRecipesByCuisine(ctx, cuisine)
}
```

### **3. Mutation with Complex Input**
```javascript
// Frontend sends:
mutation {
  insertEntry(entry: {
    name: "Apples"
    quantity: 5
    expiration: "2024-08-15T00:00:00Z"
    quantityType: "pieces"
  }) {
    id
    name
    quantity
  }
}
```

```go
// Resolver automatically receives:
func (r *mutationResolver) InsertEntry(ctx context.Context, entry entity.PantryEntryInput) (*entity.PantryEntry, error) {
    // 'entry' object is automatically extracted and parsed
    // entry = PantryEntryInput{
    //   Name: "Apples",
    //   Quantity: 5,
    //   Expiration: "2024-08-15T00:00:00Z",
    //   QuantityType: "pieces"
    // }
    
    // Add validation
    if err := validateEntry(entry); err != nil {
        return nil, err
    }
    
    return r.UseCase.RepoWrapper.PantryEntryRepo.InsertPantryEntry(ctx, &entry)
}

func validateEntry(entry entity.PantryEntryInput) error {
    if entry.Name == "" {
        return fmt.Errorf("name is required")
    }
    if entry.Quantity != nil && *entry.Quantity <= 0 {
        return fmt.Errorf("quantity must be positive")
    }
    return nil
}
```

## 🔧 Enhanced Request Extraction

### **1. Adding Request Validation**
```go
// Enhanced resolver with validation
func (r *mutationResolver) InsertEntry(ctx context.Context, entry entity.PantryEntryInput) (*entity.PantryEntry, error) {
    // Validate input
    if err := validatePantryEntry(entry); err != nil {
        return nil, fmt.Errorf("validation error: %w", err)
    }
    
    // Log the request
    log.Printf("Inserting pantry entry: %+v", entry)
    
    // Process request
    result, err := r.UseCase.RepoWrapper.PantryEntryRepo.InsertPantryEntry(ctx, &entry)
    if err != nil {
        log.Printf("Failed to insert pantry entry: %v", err)
        return nil, err
    }
    
    log.Printf("Successfully inserted pantry entry: %s", result.Name)
    return result, nil
}

func validatePantryEntry(entry entity.PantryEntryInput) error {
    if entry.Name == "" {
        return fmt.Errorf("name is required")
    }
    
    if entry.Quantity != nil && *entry.Quantity <= 0 {
        return fmt.Errorf("quantity must be positive")
    }
    
    if entry.QuantityType != nil && *entry.QuantityType == "" {
        return fmt.Errorf("quantity type cannot be empty if provided")
    }
    
    return nil
}
```

### **2. Adding Context Information**
```go
// Enhanced resolver with context
func (r *queryResolver) GetRecipe(ctx context.Context) ([]*entity.Recipe, error) {
    // Extract context information (if available)
    userID, _ := ctx.Value("userID").(string)
    requestID, _ := ctx.Value("requestID").(string)
    
    log.Printf("Getting recipes - User: %s, Request: %s", userID, requestID)
    
    recipes, err := r.UseCase.RepoWrapper.RecipeRepo.GetRecipes(ctx)
    if err != nil {
        log.Printf("Failed to get recipes: %v", err)
        return nil, err
    }
    
    log.Printf("Retrieved %d recipes", len(recipes))
    return recipes, nil
}
```

### **3. Adding Error Handling**
```go
// Enhanced resolver with error handling
func (r *mutationResolver) InsertEntry(ctx context.Context, entry entity.PantryEntryInput) (*entity.PantryEntry, error) {
    // Validate input
    if err := validatePantryEntry(entry); err != nil {
        return nil, &ValidationError{
            Field:   "entry",
            Message: err.Error(),
        }
    }
    
    // Process request
    result, err := r.UseCase.RepoWrapper.PantryEntryRepo.InsertPantryEntry(ctx, &entry)
    if err != nil {
        // Log the error
        log.Printf("Database error: %v", err)
        
        // Return user-friendly error
        return nil, &DatabaseError{
            Message: "Failed to save pantry entry",
            Code:    "DB_ERROR",
        }
    }
    
    return result, nil
}

// Custom error types
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error in field %s: %s", e.Field, e.Message)
}

type DatabaseError struct {
    Message string `json:"message"`
    Code    string `json:"code"`
}

func (e *DatabaseError) Error() string {
    return e.Message
}
```

## 🚀 Adding Middleware for Request Context

### **1. Create Middleware**
```go
// internal/adapter/delivery/graphql/middleware.go
package graphql

import (
    "context"
    "net/http"
    "github.com/google/uuid"
)

// RequestContext adds request information to context
func RequestContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Generate request ID
        requestID := uuid.New().String()
        
        // Extract user information from headers
        userID := r.Header.Get("X-User-ID")
        authToken := r.Header.Get("Authorization")
        
        // Create context with request information
        ctx := context.WithValue(r.Context(), "requestID", requestID)
        ctx = context.WithValue(ctx, "userID", userID)
        ctx = context.WithValue(ctx, "authToken", authToken)
        ctx = context.WithValue(ctx, "userAgent", r.Header.Get("User-Agent"))
        ctx = context.WithValue(ctx, "remoteAddr", r.RemoteAddr)
        
        // Call next handler with enhanced context
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### **2. Apply Middleware to Server**
```go
// cmd/pantry_butler/server.go
func main() {
    // ... existing setup ...
    
    // Create GraphQL handler
    srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{
        Resolvers: &graphql.Resolver{UseCase: *uc}
    }))
    
    // Apply middleware
    graphqlHandler := graphql.RequestContext(srv)
    
    // Setup routes
    http.Handle("/query", graphqlHandler)
    http.Handle("/", playground.Handler("GraphQL playground", "/query"))
    
    // ... rest of setup ...
}
```

### **3. Use Context in Resolvers**
```go
// Enhanced resolver using context from middleware
func (r *mutationResolver) InsertEntry(ctx context.Context, entry entity.PantryEntryInput) (*entity.PantryEntry, error) {
    // Extract context information
    requestID, _ := ctx.Value("requestID").(string)
    userID, _ := ctx.Value("userID").(string)
    userAgent, _ := ctx.Value("userAgent").(string)
    
    log.Printf("Request %s - User %s inserting pantry entry via %s", requestID, userID, userAgent)
    
    // Validate and process
    if err := validatePantryEntry(entry); err != nil {
        log.Printf("Request %s - Validation failed: %v", requestID, err)
        return nil, err
    }
    
    result, err := r.UseCase.RepoWrapper.PantryEntryRepo.InsertPantryEntry(ctx, &entry)
    if err != nil {
        log.Printf("Request %s - Database error: %v", requestID, err)
        return nil, err
    }
    
    log.Printf("Request %s - Successfully inserted entry: %s", requestID, result.Name)
    return result, nil
}
```

## 📊 Request Extraction Summary

| Request Component | Extraction Method | Example |
|-------------------|-------------------|---------|
| **GraphQL Arguments** | Automatic via resolver parameters | `func(ctx, arg string)` |
| **GraphQL Input Objects** | Automatic via resolver parameters | `func(ctx, input InputType)` |
| **Context Values** | `ctx.Value("key")` | `userID := ctx.Value("userID")` |
| **Request Headers** | Via middleware + context | `r.Header.Get("X-User-ID")` |
| **Request Metadata** | Via middleware + context | `requestID, userAgent, etc.` |

## 🎯 Best Practices

1. **Validate Early**: Validate requests in resolvers before passing to use cases
2. **Log Requests**: Log important request information for debugging
3. **Use Context**: Pass request context through all layers
4. **Handle Errors**: Return meaningful error messages
5. **Keep Resolvers Thin**: Extract data and delegate to use cases
6. **Use Middleware**: Add authentication, logging, and request context via middleware

## 🔄 Testing Request Extraction

```go
// Test request extraction
func TestInsertEntryRequestExtraction(t *testing.T) {
    // Create test context with values
    ctx := context.WithValue(context.Background(), "userID", "test-user")
    ctx = context.WithValue(ctx, "requestID", "test-request")
    
    // Create test input
    input := entity.PantryEntryInput{
        Name:     "Test Item",
        Quantity: intPtr(5),
    }
    
    // Test resolver
    resolver := &mutationResolver{}
    result, err := resolver.InsertEntry(ctx, input)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "Test Item", result.Name)
}
``` 
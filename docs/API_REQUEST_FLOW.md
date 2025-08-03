# API Request Flow in Clean Architecture

This document explains how API requests from the frontend are extracted and processed in the Pantry Butler clean architecture setup.

## 🏗️ Architecture Layers

```
Frontend → GraphQL Handler → Resolver → UseCase → Repository → Database
```

## 📡 Request Flow Overview

### 1. **Frontend Request**
```javascript
// Frontend GraphQL query
const GET_RECIPES = `
  query GetRecipes {
    getRecipe {
      id
      name
      cuisine
      description
      ingredients
    }
  }
`;

// Frontend mutation
const INSERT_PANTRY_ENTRY = `
  mutation InsertEntry($entry: PantryEntryInput!) {
    insertEntry(entry: $entry) {
      id
      name
      quantity
      expiration
    }
  }
`;
```

### 2. **HTTP Server (Entry Point)**
```go
// cmd/pantry_butler/server.go
func main() {
    // GraphQL handler setup
    srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{
        Resolvers: &graphql.Resolver{UseCase: *uc}
    }))
    
    // Route handling
    http.Handle("/query", srv)  // GraphQL endpoint
    http.Handle("/", playground.Handler("GraphQL playground", "/query"))
}
```

### 3. **GraphQL Resolver (Request Extraction)**
```go
// internal/adapter/delivery/graphql/schema.resolvers.go

// Query resolver - extracts request parameters
func (r *queryResolver) GetRecipe(ctx context.Context) ([]*entity.Recipe, error) {
    // Request is automatically extracted from GraphQL query
    // No additional parameters needed for this query
    return r.UseCase.RepoWrapper.RecipeRepo.GetRecipes(ctx)
}

// Query with parameters - extracts request parameters
func (r *queryResolver) GetRecipeByCuisine(ctx context.Context, cuisine string) ([]*entity.Recipe, error) {
    // 'cuisine' parameter is automatically extracted from GraphQL query
    return r.UseCase.RepoWrapper.RecipeRepo.GetRecipesByCuisine(ctx, cuisine)
}

// Mutation resolver - extracts request body
func (r *mutationResolver) InsertEntry(ctx context.Context, entry entity.PantryEntryInput) (*entity.PantryEntry, error) {
    // 'entry' object is automatically extracted from GraphQL mutation
    return r.UseCase.RepoWrapper.PantryEntryRepo.InsertPantryEntry(ctx, &entry)
}
```

## 🔍 Request Extraction Examples

### **Query Request Extraction**
```go
// Frontend sends:
query {
  getRecipeByCuisine(cuisine: "Italian") {
    id
    name
    cuisine
  }
}

// Resolver automatically receives:
func (r *queryResolver) GetRecipeByCuisine(ctx context.Context, cuisine string) ([]*entity.Recipe, error) {
    // cuisine = "Italian" (automatically extracted)
    return r.UseCase.RepoWrapper.RecipeRepo.GetRecipesByCuisine(ctx, cuisine)
}
```

### **Mutation Request Extraction**
```go
// Frontend sends:
mutation {
  insertEntry(entry: {
    name: "Apples"
    quantity: 5
    expiration: "2024-08-15T00:00:00Z"
  }) {
    id
    name
  }
}

// Resolver automatically receives:
func (r *mutationResolver) InsertEntry(ctx context.Context, entry entity.PantryEntryInput) (*entity.PantryEntry, error) {
    // entry = PantryEntryInput{
    //   Name: "Apples",
    //   Quantity: 5,
    //   Expiration: "2024-08-15T00:00:00Z"
    // }
    return r.UseCase.RepoWrapper.PantryEntryRepo.InsertPantryEntry(ctx, &entry)
}
```

## 🛠️ Adding Custom Request Extraction

### **1. Adding HTTP Headers Extraction**
```go
// Enhanced resolver with header extraction
func (r *mutationResolver) InsertEntry(ctx context.Context, entry entity.PantryEntryInput) (*entity.PantryEntry, error) {
    // Extract custom headers
    if gc := graphql.GetFieldContext(ctx); gc != nil {
        if req := gc.Request; req != nil {
            userID := req.Header.Get("X-User-ID")
            authToken := req.Header.Get("Authorization")
            
            // Use extracted headers in business logic
            log.Printf("User ID: %s, Auth: %s", userID, authToken)
        }
    }
    
    return r.UseCase.RepoWrapper.PantryEntryRepo.InsertPantryEntry(ctx, &entry)
}
```

### **2. Adding Query Parameters Extraction**
```go
// Enhanced resolver with query parameter extraction
func (r *queryResolver) GetRecipe(ctx context.Context) ([]*entity.Recipe, error) {
    // Extract query parameters
    if gc := graphql.GetFieldContext(ctx); gc != nil {
        if req := gc.Request; req != nil {
            limit := req.URL.Query().Get("limit")
            offset := req.URL.Query().Get("offset")
            
            // Use extracted parameters
            log.Printf("Limit: %s, Offset: %s", limit, offset)
        }
    }
    
    return r.UseCase.RepoWrapper.RecipeRepo.GetRecipes(ctx)
}
```

### **3. Adding Request Validation**
```go
// Enhanced resolver with validation
func (r *mutationResolver) InsertEntry(ctx context.Context, entry entity.PantryEntryInput) (*entity.PantryEntry, error) {
    // Validate request data
    if entry.Name == "" {
        return nil, fmt.Errorf("name is required")
    }
    
    if entry.Quantity != nil && *entry.Quantity <= 0 {
        return nil, fmt.Errorf("quantity must be positive")
    }
    
    return r.UseCase.RepoWrapper.PantryEntryRepo.InsertPantryEntry(ctx, &entry)
}
```

## 🔧 Adding REST API Support

If you want to add REST API endpoints alongside GraphQL:

### **1. Create REST Handler**
```go
// internal/adapter/delivery/rest/handler.go
package rest

import (
    "encoding/json"
    "net/http"
    "strconv"
    
    "github.com/thisausername99/pantry-butler/internal/usecase"
    entity "github.com/thisausername99/pantry-butler/internal/domain"
)

type PantryHandler struct {
    useCase *usecase.Usecase
}

func NewPantryHandler(uc *usecase.Usecase) *PantryHandler {
    return &PantryHandler{useCase: uc}
}

// GET /api/pantry-entries
func (h *PantryHandler) GetPantryEntries(w http.ResponseWriter, r *http.Request) {
    // Extract query parameters
    limitStr := r.URL.Query().Get("limit")
    offsetStr := r.URL.Query().Get("offset")
    
    // Extract headers
    userID := r.Header.Get("X-User-ID")
    
    // Process request
    entries, err := h.useCase.GetAllPantryEntries(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Return response
    json.NewEncoder(w).Encode(entries)
}

// POST /api/pantry-entries
func (h *PantryHandler) CreatePantryEntry(w http.ResponseWriter, r *http.Request) {
    // Extract request body
    var entry entity.PantryEntryInput
    if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Extract headers
    userID := r.Header.Get("X-User-ID")
    
    // Process request
    result, err := h.useCase.RepoWrapper.PantryEntryRepo.InsertPantryEntry(r.Context(), &entry)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Return response
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(result)
}
```

### **2. Add REST Routes to Server**
```go
// cmd/pantry_butler/server.go
func main() {
    // ... existing setup ...
    
    // Add REST handlers
    pantryHandler := rest.NewPantryHandler(uc)
    
    // REST routes
    http.HandleFunc("/api/pantry-entries", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            pantryHandler.GetPantryEntries(w, r)
        case "POST":
            pantryHandler.CreatePantryEntry(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })
    
    // ... existing GraphQL setup ...
}
```

## 📋 Request Extraction Summary

| Request Type | Extraction Method | Location |
|--------------|-------------------|----------|
| GraphQL Query Parameters | Automatic via resolver args | `schema.resolvers.go` |
| GraphQL Mutation Body | Automatic via resolver args | `schema.resolvers.go` |
| HTTP Headers | `graphql.GetFieldContext(ctx).Request.Header` | Resolver |
| Query Parameters | `req.URL.Query().Get()` | Resolver |
| Request Body (REST) | `json.NewDecoder(r.Body).Decode()` | REST Handler |
| Path Parameters | `mux.Vars(r)` or `chi.URLParam()` | REST Handler |

## 🎯 Best Practices

1. **Keep resolvers thin** - Extract data and delegate to use cases
2. **Validate early** - Validate requests in the delivery layer
3. **Use context** - Pass request context through all layers
4. **Handle errors gracefully** - Return appropriate HTTP status codes
5. **Log requests** - Log important request data for debugging
6. **Use middleware** - Add authentication, logging, CORS middleware 
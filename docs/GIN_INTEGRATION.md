# Gin Integration with GraphQL

This document explains how to use Gin Gonic as the middleware layer for your GraphQL server in the Pantry Butler application.

## 🏗️ Architecture Overview

```
Frontend → Gin Router → Middleware Stack → GraphQL Handler → Resolver → UseCase → Repository → MongoDB
```

## 🚀 Key Features

### **1. Comprehensive Middleware Stack**
- **Request Context**: Adds request ID, user info, and metadata
- **Request Logging**: Structured logging with request details
- **CORS**: Cross-origin resource sharing support
- **Rate Limiting**: Basic in-memory rate limiting
- **Recovery**: Panic recovery with proper error handling
- **Authentication**: Token-based authentication (optional)
- **Health Checks**: Built-in health check endpoint

### **2. Enhanced Request Extraction**
- **HTTP Headers**: Access to all request headers
- **Query Parameters**: URL query parameter extraction
- **Request Metadata**: Client IP, user agent, request ID
- **Context Propagation**: Request context flows through all layers

### **3. GraphQL Integration**
- **Seamless Integration**: GraphQL works with Gin middleware
- **Context Access**: Resolvers can access Gin context
- **Request Logging**: GraphQL requests are logged with details
- **Error Handling**: Proper error handling and logging

## 📋 Available Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | GraphQL Playground |
| `/query` | POST/GET | GraphQL endpoint |
| `/health` | GET | Health check |
| `/api/info` | GET | API information |

## 🔧 Middleware Configuration

### **Request Context Middleware**
```go
// Adds request ID, user info, and metadata to context
s.router.Use(RequestContext(logger))
```

**Features:**
- Generates unique request ID for each request
- Extracts user ID from `X-User-ID` header
- Extracts authorization token from `Authorization` header
- Adds request ID to response headers
- Logs request start with structured data

### **Request Logging Middleware**
```go
// Logs request completion with performance metrics
s.router.Use(RequestLogger(logger))
```

**Features:**
- Logs request completion with status code
- Includes request latency
- Includes request ID for correlation
- Structured logging with Zap

### **CORS Middleware**
```go
// Handles cross-origin requests
s.router.Use(CORSMiddleware())
```

**Features:**
- Allows all origins (`*`)
- Supports credentials
- Handles preflight requests
- Configurable headers and methods

### **Rate Limiting Middleware**
```go
// Basic rate limiting (100 requests per minute per IP)
s.router.Use(RateLimitMiddleware(logger))
```

**Features:**
- In-memory rate limiting
- 100 requests per minute per client IP
- Automatic cleanup of old requests
- Configurable limits

### **Authentication Middleware**
```go
// Token-based authentication (optional)
s.router.Use(AuthMiddleware(logger))
```

**Features:**
- Validates `Authorization` header
- Skips authentication for GraphQL playground
- Adds authentication status to context
- Returns proper error responses

## 🎯 Request Extraction Examples

### **1. Basic Request Extraction**
```go
// In your resolver
func (r *queryResolver) GetRecipe(ctx context.Context) ([]*entity.Recipe, error) {
    // Extract request information
    requestInfo := ExtractRequestInfo(ctx)
    
    // Access request details
    requestID := requestInfo["requestID"].(string)
    userID := requestInfo["userID"].(string)
    userAgent := requestInfo["userAgent"].(string)
    
    log.Printf("Request %s - User %s getting recipes via %s", 
        requestID, userID, userAgent)
    
    return r.UseCase.RepoWrapper.RecipeRepo.GetRecipes(ctx)
}
```

### **2. Enhanced Request Extraction with Gin**
```go
// In your enhanced resolver
func (r *queryResolver) GinEnhancedGetRecipe(ctx context.Context) ([]*entity.Recipe, error) {
    // Extract Gin context
    ginCtx := extractGinContext(ctx)
    if ginCtx == nil {
        return nil, fmt.Errorf("gin context not available")
    }

    // Extract headers
    userID := ginCtx.GetHeader("X-User-ID")
    authToken := ginCtx.GetHeader("Authorization")
    
    // Extract query parameters
    limit := ginCtx.Query("limit")
    offset := ginCtx.Query("offset")
    
    // Extract request information
    requestInfo := extractRequestInfoFromGin(ginCtx)
    
    log.Printf("Request %s - User %s getting recipes with limit=%s, offset=%s", 
        requestInfo["requestID"], userID, limit, offset)
    
    return r.UseCase.RepoWrapper.RecipeRepo.GetRecipes(ctx)
}
```

### **3. Mutation with Validation**
```go
// In your mutation resolver
func (r *mutationResolver) GinEnhancedInsertEntry(ctx context.Context, entry entity.PantryEntryInput) (*entity.PantryEntry, error) {
    // Extract Gin context
    ginCtx := extractGinContext(ctx)
    if ginCtx == nil {
        return nil, fmt.Errorf("gin context not available")
    }

    // Extract request information
    requestInfo := extractRequestInfoFromGin(ginCtx)
    userID := ginCtx.GetHeader("X-User-ID")
    
    // Validate input
    if err := validatePantryEntryInput(entry); err != nil {
        log.Printf("Validation failed - RequestID: %s, User: %s, Error: %v", 
            requestInfo["requestID"], userID, err)
        return nil, fmt.Errorf("validation error: %w", err)
    }
    
    // Process request
    result, err := r.UseCase.RepoWrapper.PantryEntryRepo.InsertPantryEntry(ctx, &entry)
    if err != nil {
        log.Printf("Failed to insert pantry entry - RequestID: %s, Error: %v", 
            requestInfo["requestID"], err)
        return nil, err
    }
    
    log.Printf("Successfully inserted pantry entry - RequestID: %s, Name: %s", 
        requestInfo["requestID"], result.Name)
    return result, nil
}
```

## 🔄 Testing with Gin

### **1. Test Server Setup**
```go
func TestGinServer(t *testing.T) {
    // Create test server
    logger := zap.NewNop()
    useCase := &usecase.Usecase{}
    server := http.NewServer(logger, useCase)
    
    // Get router for testing
    router := server.GetRouter()
    
    // Test health endpoint
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/health", nil)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

### **2. Test GraphQL with Headers**
```go
func TestGraphQLWithHeaders(t *testing.T) {
    // Create test server
    logger := zap.NewNop()
    useCase := &usecase.Usecase{}
    server := http.NewServer(logger, useCase)
    router := server.GetRouter()
    
    // Test GraphQL query with headers
    query := `query { getRecipe { id name } }`
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/query", strings.NewReader(query))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-User-ID", "test-user")
    req.Header.Set("Authorization", "Bearer test-token")
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Header().Get("X-Request-ID"), "")
}
```

## 🚀 Running the Server

### **1. Start the Server**
```bash
# Build and run
go run cmd/pantry_butler/server.go

# Or use Docker
docker-compose up server
```

### **2. Test Endpoints**
```bash
# Health check
curl http://localhost:8080/health

# API info
curl http://localhost:8080/api/info

# GraphQL query
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-User-ID: test-user" \
  -d '{"query": "{ getRecipe { id name } }"}'
```

### **3. Access GraphQL Playground**
Open your browser and navigate to: `http://localhost:8080/`

## 📊 Monitoring and Logging

### **Request Logs**
```
{"level":"info","msg":"Request started","requestID":"uuid","method":"POST","path":"/query","userID":"test-user","userAgent":"curl/7.68.0","remoteAddr":"127.0.0.1"}
{"level":"info","msg":"Request completed","requestID":"uuid","method":"POST","path":"/query","status":200,"latency":0.001234,"clientIP":"127.0.0.1"}
```

### **GraphQL Logs**
```
{"level":"info","msg":"GraphQL request","requestID":"uuid","method":"POST","path":"/query","userID":"test-user"}
```

### **Error Logs**
```
{"level":"error","msg":"GraphQL panic recovered","requestID":"uuid","panic":"runtime error","method":"POST","path":"/query"}
```

## 🎯 Best Practices

1. **Use Request IDs**: Always include request ID in logs for correlation
2. **Validate Early**: Validate requests in resolvers before processing
3. **Log Appropriately**: Log important events but avoid excessive logging
4. **Handle Errors**: Return meaningful error messages to clients
5. **Use Context**: Pass request context through all layers
6. **Monitor Performance**: Use the built-in latency logging
7. **Secure Headers**: Validate and sanitize all request headers

## 🔧 Configuration

### **Environment Variables**
```bash
# Server configuration
PORT=8080
GIN_MODE=release

# Database configuration
MONGO_URI=mongodb://root:password@mongodb:27017/pantry_butler_dev?authSource=admin
```

### **Middleware Order**
The middleware is applied in this order:
1. Recovery (handles panics)
2. Request Logging
3. CORS
4. Request Context
5. Rate Limiting
6. Health Check
7. GraphQL Context
8. Authentication (optional)

This order ensures proper error handling and request processing. 
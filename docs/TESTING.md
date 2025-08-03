# Testing Guide for Pantry Butler

This document describes the testing structure and how to run tests for the Pantry Butler application.

## Test Structure

The project uses a layered testing approach with both unit tests (using mocks) and integration tests:

```
internal/
├── usecase/
│   ├── test/
│   │   ├── mocks.go          # Generated mock interfaces
│   │   └── pantry_test.go    # Unit tests for usecase layer
│   └── usecase.go
└── adapter/
    └── persistence/
        └── mongo/
            └── test/
                ├── pantry_entry_test.go  # Integration tests for pantry entries
                └── recipe_test.go        # Integration tests for recipes
```

## Test Types

### 1. Unit Tests (Usecase Layer)
- **Location**: `internal/usecase/test/`
- **Purpose**: Test business logic in isolation using mocked dependencies
- **Framework**: Go testing + gomock + testify
- **Dependencies**: Mocked repositories

### 2. Integration Tests (Persistence Layer)
- **Location**: `internal/adapter/persistence/mongo/test/`
- **Purpose**: Test actual MongoDB operations with real database
- **Framework**: Go testing + testify
- **Dependencies**: Real MongoDB container

## Running Tests

### Prerequisites
1. Docker and Docker Compose installed
2. Go 1.21+ installed
3. MongoDB container running

### Quick Start
```bash
# Run all tests
./run_tests.sh
```

### Manual Test Execution

#### 1. Start MongoDB Container
```bash
docker-compose up -d mongodb
```

#### 2. Run Unit Tests (Usecase Layer)
```bash
go test -v ./internal/usecase/test/...
```

#### 3. Run Integration Tests (MongoDB)
```bash
go test -v ./internal/adapter/persistence/mongo/test/...
```

#### 4. Run Specific Test
```bash
# Run specific test function
go test -v -run TestGetAllPantryEntries ./internal/usecase/test/

# Run specific test file
go test -v ./internal/adapter/persistence/mongo/test/ -run TestPantryEntryRepo
```

## Test Coverage

### Usecase Tests
- `TestGetAllPantryEntries`: Tests successful retrieval of pantry entries
- `TestGetAllPantryEntries_Error`: Tests error handling when repository fails

### MongoDB Integration Tests
- `TestPantryEntryRepo_GetPantryEntries`: Tests retrieving pantry entries from MongoDB
- `TestPantryEntryRepo_InsertPantryEntry`: Tests inserting new pantry entries
- `TestPantryEntryRepo_InsertPantryEntry_EmptyName`: Tests validation for empty names
- `TestPantryEntryRepo_GetPantryEntries_Empty`: Tests behavior with empty collection
- `TestRecipeRepo_GetRecipes`: Tests retrieving recipes from MongoDB
- `TestRecipeRepo_GetRecipesByCuisine`: Tests filtering recipes by cuisine
- `TestRecipeRepo_GetRecipes_Empty`: Tests behavior with empty recipe collection
- `TestRecipeRepo_GetRecipesByCuisine_NotFound`: Tests filtering with non-existent cuisine

## Mock Generation

The project uses gomock for generating mock interfaces. To regenerate mocks:

```bash
# Install mockgen if not already installed
go install github.com/golang/mock/mockgen@latest

# Generate mocks (if using go:generate)
go generate ./internal/usecase
```

## Test Database

Integration tests use a separate test database (`pantry_butler_test`) to avoid interfering with development data. The test database is automatically created and cleaned up between tests.

## Adding New Tests

### For Usecase Layer
1. Create test in `internal/usecase/test/`
2. Use mocks for repository dependencies
3. Test both success and error scenarios

### For Persistence Layer
1. Create test in `internal/adapter/persistence/mongo/test/`
2. Use real MongoDB connection
3. Clean up test data in setup/teardown

### Example Test Structure
```go
func TestFunctionName(t *testing.T) {
    // Setup
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    // Create mocks/dependencies
    mockRepo := NewMockRepository(ctrl)
    
    // Set expectations
    mockRepo.EXPECT().Method().Return(result, nil).Times(1)
    
    // Execute
    result, err := functionUnderTest()
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

## Troubleshooting

### MongoDB Connection Issues
- Ensure Docker is running
- Check if MongoDB container is started: `docker ps | grep pantry_mongodb`
- Verify MongoDB is ready: `docker logs pantry_mongodb`

### Test Failures
- Check test database connection
- Ensure all dependencies are installed: `go mod tidy`
- Run tests with verbose output: `go test -v`

### Mock Generation Issues
- Ensure mockgen is installed: `which mockgen`
- Add to PATH if needed: `export PATH=$PATH:$(go env GOPATH)/bin` 
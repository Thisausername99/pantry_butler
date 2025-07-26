#!/bin/bash

echo "🧪 Running Pantry Butler Tests"
echo "==============================="

# Check if Docker containers are running
echo "📋 Checking if MongoDB container is running..."
if ! docker ps | grep -q "pantry_mongodb"; then
    echo "❌ MongoDB container is not running. Starting containers..."
    docker-compose up -d mongodb
    echo "⏳ Waiting for MongoDB to be ready..."
    sleep 10
else
    echo "✅ MongoDB container is already running"
fi

# Check if Go dependencies are available
echo "📦 Checking Go dependencies..."
if ! go mod tidy; then
    echo "❌ Failed to tidy Go modules"
    exit 1
fi

# Run usecase tests (unit tests with mocks)
echo "🧪 Running Usecase Tests (Unit Tests)..."
go test -v ./internal/usecase/test/...

# Run MongoDB integration tests
echo "🧪 Running MongoDB Integration Tests..."
go test -v ./internal/adapter/persistence/mongo/test/...

echo "✅ All tests completed!" 
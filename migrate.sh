#!/bin/bash

# MongoDB Migration Script
# This script runs migrations using the golang-migrate tool

set -e

# Configuration
MONGO_URI="mongodb://root:password@localhost:27017/pantry_butler_dev?authSource=admin"
MIGRATIONS_PATH="./migrations"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🔄 MongoDB Migration Tool"
echo "========================"

# Check if migrate tool is available
if ! command -v migrate &> /dev/null; then
    echo -e "${RED}❌ migrate tool not found. Please install it first:${NC}"
    echo "go install -tags 'mongodb' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
fi

# Check if MongoDB is running
if ! docker ps | grep -q "pantry_mongodb"; then
    echo -e "${YELLOW}⚠️  MongoDB container not running. Starting it...${NC}"
    docker compose up -d mongodb
    echo "⏳ Waiting for MongoDB to be ready..."
    sleep 5
fi

# Function to run migrations
run_migration() {
    local action=$1
    local version=$2
    
    echo -e "${YELLOW}🔄 Running migration ${action}...${NC}"
    
    case $action in
        "up")
            if [ -n "$2" ]; then
                # Start from specific version
                migrate -path=$MIGRATIONS_PATH -database="mongodb://root:password@localhost:27017/pantry_butler_dev?authSource=admin" up $2
            else
                # Apply all pending migrations
                migrate -path=$MIGRATIONS_PATH -database="mongodb://root:password@localhost:27017/pantry_butler_dev?authSource=admin" up
            fi
            ;;
        "down")
            migrate -path=$MIGRATIONS_PATH -database="mongodb://root:password@localhost:27017/pantry_butler_dev?authSource=admin" down
            ;;
        "force")
            migrate -path=$MIGRATIONS_PATH -database="mongodb://root:password@localhost:27017/pantry_butler_dev?authSource=admin" force $version
            ;;
        "version")
            migrate -path=$MIGRATIONS_PATH -database="mongodb://root:password@localhost:27017/pantry_butler_dev?authSource=admin" version
            ;;
        *)
            echo -e "${RED}❌ Unknown action: $action${NC}"
            exit 1
            ;;
    esac
}

# Parse command line arguments
case "${1:-help}" in
    "up")
        if [ -n "$2" ]; then
            echo -e "${YELLOW}🔄 Starting migration from version $2...${NC}"
            run_migration "up" $2
        else
            run_migration "up"
        fi
        echo -e "${GREEN}✅ Migrations applied successfully!${NC}"
        ;;
    "down")
        run_migration "down"
        echo -e "${GREEN}✅ Migrations rolled back successfully!${NC}"
        ;;
    "force")
        if [ -z "$2" ]; then
            echo -e "${RED}❌ Version number required for force command${NC}"
            exit 1
        fi
        run_migration "force" $2
        echo -e "${GREEN}✅ Migration forced to version $2!${NC}"
        ;;
    "version")
        run_migration "version"
        ;;
    "status")
        echo -e "${YELLOW}📊 Migration Status:${NC}"
        migrate -path=$MIGRATIONS_PATH -database="mongodb://root:password@localhost:27017/pantry_butler_dev?authSource=admin" version
        ;;
    "help"|*)
        echo "Usage: $0 [command]"
        echo ""
        echo "Commands:"
        echo "  up [N]  - Apply all pending migrations (or from version N)"
        echo "  down    - Rollback all migrations"
        echo "  force N - Force migration to version N"
        echo "  version - Show current migration version"
        echo "  status  - Show migration status"
        echo "  help    - Show this help message"
        echo ""
        echo "Examples:"
        echo "  $0 up                    # Apply all migrations"
        echo "  $0 up 2                  # Apply migrations starting from version 2"
        echo "  $0 down                  # Rollback all migrations"
        echo "  $0 force 1               # Force migration to version 1"
        echo "  $0 status                # Check migration status"
        ;;
esac 
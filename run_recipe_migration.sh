#!/bin/bash

# Recipe Migration Runner
# This script runs the recipe JSON migration using MongoDB commands

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "🍳 Recipe Migration Runner"
echo "========================="

# Check if MongoDB is running
if ! docker ps | grep -q "pantry_mongodb"; then
    echo -e "${YELLOW}⚠️  MongoDB container not running. Starting it...${NC}"
    docker-compose up -d mongodb
    echo "⏳ Waiting for MongoDB to be ready..."
    sleep 5
fi

# Function to run migration
run_migration() {
    local action=$1
    
    case $action in
        "up")
            echo -e "${BLUE}📄 Running recipe migration (UP)...${NC}"
            
            # Insert recipes
            echo -e "${YELLOW}🔄 Inserting recipes...${NC}"
            docker exec pantry_mongodb mongosh --quiet --eval "
                use pantry_butler_dev;
                db.recipes.insertMany([
                    {
                        id: 'recipe_001',
                        name: 'Apple Pie',
                        cuisine: 'American',
                        ingredients: {
                            apples: 5,
                            flour: '2 cups',
                            sugar: '1 cup',
                            butter: '1/2 cup',
                            cinnamon: '1 tsp',
                            salt: '1/4 tsp'
                        },
                        difficulty: 3,
                        description: 'A delicious homemade apple pie with a flaky crust',
                        prep_time: 30,
                        cook_time: 45,
                        servings: 8,
                        rating: 4,
                        source_url: 'https://example.com/apple-pie',
                        created_at: new Date(),
                        updated_at: new Date()
                    },
                    {
                        id: 'recipe_002',
                        name: 'Pasta Carbonara',
                        cuisine: 'Italian',
                        ingredients: {
                            pasta: '500g',
                            eggs: 4,
                            bacon: '200g',
                            parmesan: '100g',
                            black_pepper: '1 tsp',
                            salt: '1 tsp'
                        },
                        difficulty: 2,
                        description: 'Classic Italian pasta dish with eggs, bacon, and parmesan cheese',
                        prep_time: 15,
                        cook_time: 20,
                        servings: 4,
                        rating: 5,
                        source_url: 'https://example.com/pasta-carbonara',
                        created_at: new Date(),
                        updated_at: new Date()
                    },
                    {
                        id: 'recipe_003',
                        name: 'Grilled Cheese Sandwich',
                        cuisine: 'American',
                        ingredients: {
                            bread: 2,
                            cheese: '4 slices',
                            butter: '2 tbsp'
                        },
                        difficulty: 1,
                        description: 'Simple and delicious grilled cheese sandwich',
                        prep_time: 5,
                        cook_time: 10,
                        servings: 1,
                        rating: 3,
                        source_url: null,
                        created_at: new Date(),
                        updated_at: new Date()
                    },
                    {
                        id: 'recipe_004',
                        name: 'Chicken Curry',
                        cuisine: 'Indian',
                        ingredients: {
                            chicken: '500g',
                            onion: 2,
                            garlic: '4 cloves',
                            ginger: '1 inch',
                            curry_powder: '2 tbsp',
                            coconut_milk: '400ml',
                            rice: '2 cups'
                        },
                        difficulty: 3,
                        description: 'Spicy and aromatic chicken curry with coconut milk',
                        prep_time: 20,
                        cook_time: 40,
                        servings: 6,
                        rating: 4,
                        source_url: 'https://example.com/chicken-curry',
                        created_at: new Date(),
                        updated_at: new Date()
                    },
                    {
                        id: 'recipe_005',
                        name: 'Caesar Salad',
                        cuisine: 'Italian',
                        ingredients: {
                            romaine_lettuce: '1 head',
                            parmesan: '1/2 cup',
                            croutons: '1 cup',
                            lemon: 1,
                            olive_oil: '3 tbsp',
                            garlic: '2 cloves',
                            anchovies: '4 fillets'
                        },
                        difficulty: 1,
                        description: 'Fresh and crisp Caesar salad with homemade dressing',
                        prep_time: 15,
                        cook_time: 0,
                        servings: 4,
                        rating: 4,
                        source_url: 'https://example.com/caesar-salad',
                        created_at: new Date(),
                        updated_at: new Date()
                    }
                ]);
            "
            
            # Create indexes
            echo -e "${YELLOW}🔄 Creating indexes...${NC}"
            docker exec pantry_mongodb mongosh --quiet --eval "
                use pantry_butler_dev;
                db.recipes.createIndex({ id: 1 }, { unique: true, name: 'id_1' });
                db.recipes.createIndex({ name: 1 }, { name: 'name_1' });
                db.recipes.createIndex({ cuisine: 1 }, { name: 'cuisine_1' });
                db.recipes.createIndex({ difficulty: 1 }, { name: 'difficulty_1' });
                db.recipes.createIndex({ rating: -1 }, { name: 'rating_-1' });
                db.recipes.createIndex({ created_at: 1 }, { name: 'created_at_1' });
            "
            
            echo -e "${GREEN}✅ Recipes inserted and indexes created successfully!${NC}"
            ;;
            
        "down")
            echo -e "${BLUE}📄 Running recipe migration (DOWN)...${NC}"
            
            # Delete recipes
            echo -e "${YELLOW}🔄 Deleting recipes...${NC}"
            docker exec pantry_mongodb mongosh --quiet --eval "
                use pantry_butler_dev;
                db.recipes.deleteMany({ id: { \$in: ['recipe_001', 'recipe_002', 'recipe_003', 'recipe_004', 'recipe_005'] } });
            "
            
            echo -e "${GREEN}✅ Recipes deleted successfully!${NC}"
            ;;
            
        "status")
            echo -e "${BLUE}📊 Recipe Collection Status:${NC}"
            docker exec pantry_mongodb mongosh --quiet --eval "
                use pantry_butler_dev;
                print('Total recipes:', db.recipes.countDocuments());
                print('Indexes:');
                db.recipes.getIndexes().forEach(function(index) {
                    print('  -', index.name, ':', JSON.stringify(index.key));
                });
            "
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
        run_migration "up"
        ;;
    "down")
        run_migration "down"
        ;;
    "status")
        run_migration "status"
        ;;
    "help"|*)
        echo "Usage: $0 [command]"
        echo ""
        echo "Commands:"
        echo "  up      - Apply recipe migration (insert recipes and create indexes)"
        echo "  down    - Rollback recipe migration (delete recipes)"
        echo "  status  - Show recipe collection status"
        echo "  help    - Show this help message"
        echo ""
        echo "Examples:"
        echo "  $0 up      # Insert recipes and create indexes"
        echo "  $0 down    # Delete recipes"
        echo "  $0 status  # Check recipe collection status"
        ;;
esac 
-- Migration: 000001_add_recipes_collection.down.sql
-- Description: Rollback recipes collection creation

-- Drop the recipes collection
db.recipes.drop(); 
-- Migration: 000002_add_pantries_collection.up.sql
-- Description: Create pantries collection and add initial data

-- Create pantries collection
db.createCollection("pantries");

-- Create index on id field
db.pantries.createIndex(
    { "id": 1 },
    { 
        "unique": true,
        "name": "id_1"
    }
);

-- Insert initial pantry data
db.pantries.insertMany([
    {
        "id": "1",
        "name": "Flour",
        "expiration": new Date("2024-12-31T00:00:00Z"),
        "quantity": 2,
        "quantityType": "kg",
        "created_at": new Date(),
        "updated_at": new Date()
    },
    {
        "id": "2",
        "name": "Sugar",
        "expiration": new Date("2024-11-30T00:00:00Z"),
        "quantity": 1,
        "quantityType": "kg",
        "created_at": new Date(),
        "updated_at": new Date()
    },
    {
        "id": "3",
        "name": "Salt",
        "expiration": null,
        "quantity": 0.5,
        "quantityType": "kg",
        "created_at": new Date(),
        "updated_at": new Date()
    }
]); 
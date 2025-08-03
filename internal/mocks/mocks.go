//go:generate mockgen -destination=entity_repo_mocks.go -package=mocks github.com/thisausername99/pantry_butler/internal/domain/repository PantryRepository,RecipeRepository,UserRepository
//go:generate mockgen -destination=mongo_mocks.go -package=mocks github.com/thisausername99/pantry_butler/internal/persistence/mongo MongoDB,MongoDatabase,MongoCollection,MongoCursor,MongoSingleResult,MongoInsertOneResult,MongoInsertManyResult,MongoUpdateResult,MongoDeleteResult,MongoIndexView

package mocks

// This file is used to generate all mocks in one centralized location

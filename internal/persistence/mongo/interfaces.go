package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB represents MongoDB client interface
type MongoDB interface {
	// Database operations
	Database(name string) MongoDatabase
	Connect(ctx context.Context) error
	Ping(ctx context.Context) error
	Disconnect(ctx context.Context) error
}

// MongoDatabase represents MongoDB database interface
type MongoDatabase interface {
	Collection(name string) MongoCollection
	RunCommand(ctx context.Context, runCommand interface{}) MongoSingleResult
	Drop(ctx context.Context) error
}

// MongoCollection represents MongoDB collection interface
type MongoCollection interface {
	// Insert operations
	InsertOne(ctx context.Context, document interface{}) (MongoInsertOneResult, error)
	InsertMany(ctx context.Context, documents []interface{}) (MongoInsertManyResult, error)

	// Find operations
	FindOne(ctx context.Context, filter interface{}) MongoSingleResult
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (MongoCursor, error)

	// Update operations
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (MongoUpdateResult, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}) (MongoUpdateResult, error)
	ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}) (MongoUpdateResult, error)

	// Delete operations
	DeleteOne(ctx context.Context, filter interface{}) (MongoDeleteResult, error)
	DeleteMany(ctx context.Context, filter interface{}) (MongoDeleteResult, error)

	// Aggregation
	Aggregate(ctx context.Context, pipeline interface{}) (MongoCursor, error)
	CountDocuments(ctx context.Context, filter interface{}) (int64, error)

	// Index operations
	Indexes() MongoIndexView
	Drop(ctx context.Context) error
}

// MongoCursor represents MongoDB cursor interface
type MongoCursor interface {
	Next(ctx context.Context) bool
	Decode(val interface{}) error
	All(ctx context.Context, results interface{}) error
	Close(ctx context.Context) error
	Err() error
}

// MongoSingleResult represents MongoDB single result interface
type MongoSingleResult interface {
	Decode(v interface{}) error
	Err() error
}

// MongoInsertOneResult represents insert one result
type MongoInsertOneResult interface {
	InsertedID() interface{}
}

// MongoInsertManyResult represents insert many result
type MongoInsertManyResult interface {
	InsertedIDs() []interface{}
}

// MongoUpdateResult represents update result
type MongoUpdateResult interface {
	MatchedCount() int64
	ModifiedCount() int64
	UpsertedCount() int64
	UpsertedID() interface{}
}

// MongoDeleteResult represents delete result
type MongoDeleteResult interface {
	DeletedCount() int64
}

// MongoIndexView represents index operations
type MongoIndexView interface {
	CreateOne(ctx context.Context, model mongo.IndexModel) (string, error)
	CreateMany(ctx context.Context, models []mongo.IndexModel) ([]string, error)
	DropOne(ctx context.Context, name string) error
	List(ctx context.Context) (MongoCursor, error)
}

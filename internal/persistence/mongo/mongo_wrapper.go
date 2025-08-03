package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoClient implements MongoDB interface with real mongo client
type mongoClient struct {
	client *mongo.Client
}

// NewMongoClient creates a new MongoDB client
func NewMongoClient(uri string) (MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &mongoClient{client: client}, nil
}

func (c *mongoClient) Database(name string) MongoDatabase {
	return &mongoDatabase{db: c.client.Database(name)}
}

func (c *mongoClient) Connect(ctx context.Context) error {
	return c.client.Connect(ctx)
}

func (c *mongoClient) Ping(ctx context.Context) error {
	return c.client.Ping(ctx, nil)
}

func (c *mongoClient) Disconnect(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

// mongoDatabase implements MongoDatabase interface
type mongoDatabase struct {
	db *mongo.Database
}

func (d *mongoDatabase) Collection(name string) MongoCollection {
	return &mongoCollection{coll: d.db.Collection(name)}
}

func (d *mongoDatabase) RunCommand(ctx context.Context, runCommand interface{}) MongoSingleResult {
	return &mongoSingleResult{result: d.db.RunCommand(ctx, runCommand)}
}

func (d *mongoDatabase) Drop(ctx context.Context) error {
	return d.db.Drop(ctx)
}

// mongoCollection implements MongoCollection interface
type mongoCollection struct {
	coll *mongo.Collection
}

func (c *mongoCollection) InsertOne(ctx context.Context, document interface{}) (MongoInsertOneResult, error) {
	result, err := c.coll.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return &mongoInsertOneResult{result: result}, nil
}

func (c *mongoCollection) InsertMany(ctx context.Context, documents []interface{}) (MongoInsertManyResult, error) {
	result, err := c.coll.InsertMany(ctx, documents)
	if err != nil {
		return nil, err
	}
	return &mongoInsertManyResult{result: result}, nil
}

func (c *mongoCollection) FindOne(ctx context.Context, filter interface{}) MongoSingleResult {
	return &mongoSingleResult{result: c.coll.FindOne(ctx, filter)}
}

func (c *mongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (MongoCursor, error) {
	cursor, err := c.coll.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	return &mongoCursor{cursor: cursor}, nil
}

func (c *mongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (MongoUpdateResult, error) {
	result, err := c.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &mongoUpdateResult{result: result}, nil
}

func (c *mongoCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{}) (MongoUpdateResult, error) {
	result, err := c.coll.UpdateMany(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &mongoUpdateResult{result: result}, nil
}

func (c *mongoCollection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}) (MongoUpdateResult, error) {
	result, err := c.coll.ReplaceOne(ctx, filter, replacement)
	if err != nil {
		return nil, err
	}
	return &mongoUpdateResult{result: result}, nil
}

func (c *mongoCollection) DeleteOne(ctx context.Context, filter interface{}) (MongoDeleteResult, error) {
	result, err := c.coll.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &mongoDeleteResult{result: result}, nil
}

func (c *mongoCollection) DeleteMany(ctx context.Context, filter interface{}) (MongoDeleteResult, error) {
	result, err := c.coll.DeleteMany(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &mongoDeleteResult{result: result}, nil
}

func (c *mongoCollection) Aggregate(ctx context.Context, pipeline interface{}) (MongoCursor, error) {
	cursor, err := c.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	return &mongoCursor{cursor: cursor}, nil
}

func (c *mongoCollection) CountDocuments(ctx context.Context, filter interface{}) (int64, error) {
	return c.coll.CountDocuments(ctx, filter)
}

func (c *mongoCollection) Indexes() MongoIndexView {
	return &mongoIndexView{view: c.coll.Indexes()}
}

func (c *mongoCollection) Drop(ctx context.Context) error {
	return c.coll.Drop(ctx)
}

// mongoCursor implements MongoCursor interface
type mongoCursor struct {
	cursor *mongo.Cursor
}

func (c *mongoCursor) Next(ctx context.Context) bool {
	return c.cursor.Next(ctx)
}

func (c *mongoCursor) Decode(val interface{}) error {
	return c.cursor.Decode(val)
}

func (c *mongoCursor) All(ctx context.Context, results interface{}) error {
	return c.cursor.All(ctx, results)
}

func (c *mongoCursor) Close(ctx context.Context) error {
	return c.cursor.Close(ctx)
}

func (c *mongoCursor) Err() error {
	return c.cursor.Err()
}

// mongoSingleResult implements MongoSingleResult interface
type mongoSingleResult struct {
	result *mongo.SingleResult
}

func (r *mongoSingleResult) Decode(v interface{}) error {
	return r.result.Decode(v)
}

func (r *mongoSingleResult) Err() error {
	return r.result.Err()
}

// mongoInsertOneResult implements MongoInsertOneResult interface
type mongoInsertOneResult struct {
	result *mongo.InsertOneResult
}

func (r *mongoInsertOneResult) InsertedID() interface{} {
	return r.result.InsertedID
}

// mongoInsertManyResult implements MongoInsertManyResult interface
type mongoInsertManyResult struct {
	result *mongo.InsertManyResult
}

func (r *mongoInsertManyResult) InsertedIDs() []interface{} {
	return r.result.InsertedIDs
}

// mongoUpdateResult implements MongoUpdateResult interface
type mongoUpdateResult struct {
	result *mongo.UpdateResult
}

func (r *mongoUpdateResult) MatchedCount() int64 {
	return r.result.MatchedCount
}

func (r *mongoUpdateResult) ModifiedCount() int64 {
	return r.result.ModifiedCount
}

func (r *mongoUpdateResult) UpsertedCount() int64 {
	return r.result.UpsertedCount
}

func (r *mongoUpdateResult) UpsertedID() interface{} {
	return r.result.UpsertedID
}

// mongoDeleteResult implements MongoDeleteResult interface
type mongoDeleteResult struct {
	result *mongo.DeleteResult
}

func (r *mongoDeleteResult) DeletedCount() int64 {
	return r.result.DeletedCount
}

// mongoIndexView implements MongoIndexView interface
type mongoIndexView struct {
	view mongo.IndexView
}

func (v *mongoIndexView) CreateOne(ctx context.Context, model mongo.IndexModel) (string, error) {
	return v.view.CreateOne(ctx, model)
}

func (v *mongoIndexView) CreateMany(ctx context.Context, models []mongo.IndexModel) ([]string, error) {
	return v.view.CreateMany(ctx, models)
}

func (v *mongoIndexView) DropOne(ctx context.Context, name string) error {
	_, err := v.view.DropOne(ctx, name)
	return err
}

func (v *mongoIndexView) List(ctx context.Context) (MongoCursor, error) {
	cursor, err := v.view.List(ctx)
	if err != nil {
		return nil, err
	}
	return &mongoCursor{cursor: cursor}, nil
}

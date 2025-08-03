package mongo

import (
	"context"

	"github.com/thisausername99/pantry_butler/internal/domain/entity"
	"github.com/thisausername99/pantry_butler/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type UserRepo struct {
	Collection MongoCollection
	Logger     *zap.Logger
}

// Ensure it implements the interface
var _ repository.UserRepository = (*UserRepo)(nil)

func (m *UserRepo) GetUser(ctx context.Context, id string) (*entity.User, error) {
	filter := bson.M{"id": id}
	var user entity.User
	err := m.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserRepo) CreateUser(ctx context.Context, user *entity.User) error {
	_, err := m.Collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	filter := bson.M{"email": email}
	var user entity.User
	err := m.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserRepo) UpdateUserWithPantry(ctx context.Context, userID string, pantryID string) error {
	filter := bson.M{"id": userID}
	_, err := m.Collection.UpdateOne(ctx, filter, bson.M{"$push": bson.M{"pantries": pantryID}})
	if err != nil {
		return err
	}
	return nil
}

func (m *UserRepo) DeletePantryFromUser(ctx context.Context, userID string, pantryID string) error {
	filter := bson.M{"id": userID}
	_, err := m.Collection.UpdateOne(ctx, filter, bson.M{"$pull": bson.M{"pantries": pantryID}})
	if err != nil {
		return err
	}
	return nil
}

func (m *UserRepo) UpdateUser(ctx context.Context, userID string, user *entity.User) error {
	filter := bson.M{"id": userID}
	_, err := m.Collection.UpdateOne(ctx, filter, bson.M{"$set": user})
	if err != nil {
		return err
	}
	return nil
}

package mongo

import (
	"context"

	"github.com/thisausername99/pantry_butler/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type UserRepo struct {
	Collection MongoCollection
	Logger     *zap.Logger
}

func (m *UserRepo) GetUser(ctx context.Context, id string) (*entity.User, error) {
	filter := bson.M{"id": id}
	var user entity.User
	err := m.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserRepo) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	_, err := m.Collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
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

package mongo

import (
	"context"

	"github.com/thisausername99/pantry_butler/internal/domain/entity"
	"github.com/thisausername99/pantry_butler/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (m *UserRepo) DeleteUser(ctx context.Context, userID string) error {
	filter := bson.M{"id": userID}
	_, err := m.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

// User authentication & security
func (m *UserRepo) UpdatePassword(ctx context.Context, userID string, hashedPassword string) error {
	filter := bson.M{"id": userID}
	_, err := m.Collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"password": hashedPassword}})
	if err != nil {
		return err
	}
	return nil
}

func (m *UserRepo) GetUserByCredentials(ctx context.Context, email, hashedPassword string) (*entity.User, error) {
	var user *entity.User
	filter := bson.M{"email": email, "password": hashedPassword}
	err := m.Collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// User profile management
func (m *UserRepo) UpdateProfile(ctx context.Context, userID string, firstName, lastName string) error {
	filter := bson.M{"id": userID}
	_, err := m.Collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"first_name": firstName, "last_name": lastName}})
	if err != nil {
		return err
	}
	return nil
}

func (m *UserRepo) UpdateEmail(ctx context.Context, userID string, email string) error {
	filter := bson.M{"id": userID}
	_, err := m.Collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"email": email}})
	if err != nil {
		return err
	}
	return nil
}

// User pantry management
func (m *UserRepo) GetUserPantries(ctx context.Context, userID string) ([]string, error) {
	filter := bson.M{"id": userID}
	opts := options.FindOne().SetProjection(bson.M{"pantries": 1, "_id": 0})

	var result struct {
		Pantries []string `bson:"pantries"`
	}

	err := m.Collection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Pantries, nil
}

// User search & listing
func (m *UserRepo) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	var users []*entity.User
	cursor, err := m.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// User status & verification
// func (m *UserRepo) UpdateUserStatus(ctx context.Context, userID string, status string) error {
// 	panic("Not implemented!")
// }

// func (m *UserRepo) VerifyUserEmail(ctx context.Context, userID string) error {
// 	panic("Not implemented!")
// }

// func (m *UserRepo) SetEmailVerificationToken(ctx context.Context, userID string, token string) error {
// 	panic("Not implemented!")
// }

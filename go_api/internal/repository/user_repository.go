package repository

import (
	"context"
	"sora_chat/internal/database"
	"sora_chat/internal/model"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var userCollectionName = "user"

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)

	// Thêm
	Insert(user model.User) (*string, error)
}

type userRepo struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	coll := database.GetCollection(db, userCollectionName)
	return &userRepo{
		collection: coll,
	}
}

func (r *userRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Thêm
func (r *userRepo) Insert(user model.User) (*string, error) {
	user.ID = uuid.New().String()
	user.CreatedDate = time.Now()
	user.UpdatedDate = time.Now()
	_, err := r.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	return &user.ID, err
}

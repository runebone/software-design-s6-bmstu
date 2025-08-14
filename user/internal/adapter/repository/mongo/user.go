package repository

import (
	"context"
	"time"
	"user/internal/entity"
	"user/internal/repository"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUserRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		db:         db,
		collection: db.Collection("users"),
	}
}

func (r *MongoUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	repoUser := repository.RepoUser(*user)

	repoUser.Role = "user"
	repoUser.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, repoUser)

	return err
}

func (r *MongoUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var repoUser repository.User

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&repoUser)
	if err != nil {
		return nil, err
	}

	user := repository.UserToEntity(repoUser)

	return &user, nil
}

func (r *MongoUserRepository) GetUsers(ctx context.Context, filter repository.UserFilter) ([]entity.User, error) {
	query := bson.M{}

	if filter.ID != nil {
		query["_id"] = *filter.ID
	}
	if filter.Email != nil {
		query["email"] = *filter.Email
	}
	if filter.Username != nil {
		query["username"] = *filter.Username
	}

	cursor, err := r.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var repoUsers []repository.User
	for cursor.Next(ctx) {
		var repoUser repository.User
		if err := cursor.Decode(&repoUser); err != nil {
			return nil, err
		}
		repoUsers = append(repoUsers, repoUser)
	}

	users := make([]entity.User, len(repoUsers))
	for i, u := range repoUsers {
		users[i] = repository.UserToEntity(u)
	}

	return users, nil
}

func (r *MongoUserRepository) GetUsersBatch(ctx context.Context, limit, offset int) ([]entity.User, error) {
	options := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset)).SetSort(bson.M{"created_at": 1})

	cursor, err := r.collection.Find(ctx, bson.M{}, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var repoUsers []repository.User
	for cursor.Next(ctx) {
		var repoUser repository.User
		if err := cursor.Decode(&repoUser); err != nil {
			return nil, err
		}
		repoUsers = append(repoUsers, repoUser)
	}

	users := make([]entity.User, len(repoUsers))
	for i, u := range repoUsers {
		users[i] = repository.UserToEntity(u)
	}

	return users, nil
}

func (r *MongoUserRepository) GetNewUsers(ctx context.Context, from time.Time, to time.Time) ([]entity.User, error) {
	query := bson.M{
		"created_at": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}

	cursor, err := r.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var repoUsers []repository.User
	for cursor.Next(ctx) {
		var repoUser repository.User
		if err := cursor.Decode(&repoUser); err != nil {
			return nil, err
		}
		repoUsers = append(repoUsers, repoUser)
	}

	users := make([]entity.User, len(repoUsers))
	for i, u := range repoUsers {
		users[i] = repository.UserToEntity(u)
	}

	return users, nil
}

func (r *MongoUserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	repoUser := repository.RepoUser(*user)
	repoUser.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"username":   repoUser.Username,
			"email":      repoUser.Email,
			"updated_at": repoUser.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": repoUser.ID}, update)

	return err
}

func (r *MongoUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

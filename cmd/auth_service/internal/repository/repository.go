package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"uber-popug/cmd/auth_service/internal/types"
)

type Repository struct {
	collection *mongo.Collection
}

func (r *Repository) OnStart() error {
	ctx := context.Background()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	r.collection = client.Database("uber-popug").Collection("users")

	return nil
}

func (r *Repository) CreateUser(ctx context.Context, user *types.User) error {
	userToInsert, err := UserToMongoType(user)
	if err != nil {
		return fmt.Errorf("convert user to mongo type: %s", err)
	}

	_, err = r.collection.InsertOne(ctx, userToInsert)

	return err
}

func (r *Repository) GetAll(ctx context.Context) ([]*types.User, error) {
	// passing bson.D{{}} matches all documents in the collection
	filter := bson.D{{}}

	return r.filterTasks(ctx, filter)
}

func (r *Repository) filterTasks(ctx context.Context, filter interface{}) ([]*types.User, error) {
	// A slice of tasks for storing the decoded documents
	var users []*types.User

	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return users, err
	}

	for cur.Next(ctx) {
		var t types.User
		err := cur.Decode(&t)
		if err != nil {
			return users, err
		}

		users = append(users, &t)
	}

	if err := cur.Err(); err != nil {
		return users, err
	}

	// once exhausted, close the cursor
	_ = cur.Close(ctx)

	if len(users) == 0 {
		return users, mongo.ErrNoDocuments
	}

	return users, nil
}

package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"uber-popug/cmd/auth_service/internal/types"
)

type User struct {
	ID       primitive.ObjectID `bson:"id"`
	Name     string             `bson:"name"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

func UserToMongoType(u *types.User) (*User, error) {
	id, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       id,
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

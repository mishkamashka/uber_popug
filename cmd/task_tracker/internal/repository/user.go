package repository

import (
	"time"
	"uber-popug/pkg/types"
)

type User struct {
	ID        string    `gorm:"primarykey"`
	Username  string    `json:"username" gorm:"unique"`
	Email     string    `json:"email" gorm:"unique"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


func RepoTypeToUser(u *User) *types.User {
	return &types.User{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}

func RepoTypesToUsers(u []*User) []*types.User {
	res := make([]*types.User, 0, len(u))

	for _, user := range u {
		resUser := RepoTypeToUser(user)

		res = append(res, resUser)
	}

	return res
}

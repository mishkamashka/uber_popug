package repository

import (
	"github.com/hashicorp/go-uuid"
	"time"
	"uber-popug/pkg/types"
)

type User struct {
	ID        string    `gorm:"primarykey"`
	Name      string    `json:"name"`
	Username  string    `json:"username" gorm:"unique"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func UserToRepoType(u *types.User) *User {
	id, _ := uuid.GenerateUUID()
	return &User{
		ID:       id,
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		Role:     u.Role,
	}
}

func RepoTypeToUser(u *User) *types.User {
	return &types.User{
		ID:       u.ID,
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		Role:     u.Role,
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

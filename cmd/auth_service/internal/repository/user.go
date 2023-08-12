package repository

import (
	"time"

	"uber-popug/cmd/auth_service/internal/types"
)

type User struct {
	ID        string    `gorm:"primarykey"`
	Name      string    `json:"name"`
	Username  string    `json:"username" gorm:"unique"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func UserToRepoType(u *types.User) (*User, error) {
	return &User{
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

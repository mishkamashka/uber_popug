package v1

import (
	"time"
)

type UserMessageType uint8

const (
	UserRoleUpdated UserMessageType = iota
	UserCreated
	UserDeleted
	UserUpdated
)

type UserMessage struct {
	Type      UserMessageType `json:"type"`
	UserData  UserData        `json:"user_data"`
	CreatedAt time.Time       `json:"created_at"`
}

type UserData struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

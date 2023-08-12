package messages

import "uber-popug/pkg/types"

type UserMessageType uint8

const (
	UserRoleUpdated UserMessageType = iota
	UserCreated
	UserDeleted
	UserUpdated
)

type UserMessage struct {
	Type     UserMessageType `json:"type"`
	UserData *types.User     `json:"user_data"`
}

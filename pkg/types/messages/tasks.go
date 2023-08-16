package messages

import (
	"time"
	"uber-popug/pkg/types"
)

type TaskMessageType uint8

const (
	TaskCreated TaskMessageType = iota
	TaskDeleted
	TaskClosed
	TaskReassigned
)

type TaskMessage struct {
	Type      TaskMessageType `json:"type"`
	Data      *types.Task     `json:"user_data"`
	CreatedAt time.Time       `json:"created_at"`
}

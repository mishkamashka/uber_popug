package messages

import "uber-popug/pkg/types"

type TaskMessageType uint8

const (
	TaskCreated TaskMessageType = iota
	TaskDeleted
	TaskClosed
	TaskReassigned
)

type TaskMessage struct {
	Type TaskMessageType `json:"type"`
	Data *types.Task     `json:"user_data"`
}

package v1

import (
	"time"
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
	Data      TaskData        `json:"user_data"`
	CreatedAt time.Time       `json:"created_at"`
}

type TaskData struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	PriceForAssign  uint8     `json:"price_for_assign"`
	PriceForClosing uint8     `json:"price_for_closing"`
	AssigneeId      string    `json:"assignee_id"`
	CreatorId       string    `json:"creator_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

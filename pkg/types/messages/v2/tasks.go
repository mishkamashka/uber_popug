package v2

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
	ID        string          `json:"id"`
	Type      TaskMessageType `json:"type"`
	Data      TaskData        `json:"user_data"`
	CreatedAt time.Time       `json:"created_at"`
}

type TaskData struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	JiraID          string    `json:"jira_id"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	PriceForAssign  uint8     `json:"price_for_assign"`
	PriceForClosing uint8     `json:"price_for_closing"`
	AssigneeId      string    `json:"assignee_id"`
	CreatorId       string    `json:"creator_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

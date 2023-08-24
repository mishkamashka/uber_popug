package types

import "time"

type AuditLog struct {
	ID        string    `json:"primarykey"`
	UserID    string    `json:"owner_id"`
	Amount    int       `json:"amount"`
	Reason    string    `json:"reason"`
	TaskInfo  *TaskInfo `json:"task_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TaskInfo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	JiraID      string `json:"jira_id"`
	Description string `json:"description"`
}

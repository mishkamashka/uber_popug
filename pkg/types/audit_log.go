package types

import "time"

type AuditLog struct {
	ID        string    `json:"primarykey"`
	UserID    string    `json:"owner_id"`
	Amount    int32     `json:"amount"`
	Reason    string    `json:"reason"`
	TaskID    string    `json:"task_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

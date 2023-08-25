package v1

import "time"

type TransactionMessage struct {
	ID        string          `json:"id"`
	Data      TransactionData `json:"user_data"`
	CreatedAt time.Time       `json:"created_at"`
}

type TransactionData struct {
	AuditLogID  string    `json:"audit_log_id"`
	UserID      string    `json:"owner_id"`
	Amount      int       `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

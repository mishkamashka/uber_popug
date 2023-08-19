package types

import "time"

type Balance struct {
	ID        string    `json:"primarykey"`
	UserID    string    `json:"owner_id"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

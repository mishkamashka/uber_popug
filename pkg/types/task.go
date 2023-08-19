package types

import (
	"math/rand"
	"time"
)

type Task struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	JiraID          string    `json:"jira_id"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	PriceForAssign  uint8     `json:"price_for_assign"`
	PriceForClosing uint8     `json:"price_for_closing"`
	AssigneeId      string    `json:"assignee_id"`
	CreatorId       string    `json:"creator_id"`
	AssignedAt      time.Time `json:"assigned_at"`
	ClosedAt        time.Time `json:"closed_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (t *Task) GeneratePrices() {
	t.PriceForAssign = uint8(rand.Intn(10) + 10)
	t.PriceForClosing = uint8(rand.Intn(20) + 20)
}

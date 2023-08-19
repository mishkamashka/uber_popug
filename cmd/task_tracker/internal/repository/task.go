package repository

import (
	"time"
	"uber-popug/pkg/types"
)

type Task struct {
	ID              string    `gorm:"primarykey"`
	Status          string    `gorm:"status"`
	Title           string    `gorm:"title"`
	JiraID          string    `gorm:"jira_id"`
	Description     string    `gorm:"description"`
	AssigneeId      string    `gorm:"assignee_id"`
	PriceForAssign  uint8     `gorm:"price_for_assign"`
	PriceForClosing uint8     `gorm:"price_for_closing"`
	CreatorId       string    `gorm:"creator_id"`
	AssignedAt      time.Time `gorm:"assigned_at"`
	ClosedAt        time.Time `gorm:"closed_at"`
	CreatedAt       time.Time `gorm:"created_at"`
	UpdatedAt       time.Time `gorm:"updated_at"`
}

func TaskToRepoType(u *types.Task) *Task {
	return &Task{
		ID:              u.ID,
		Status:          u.Status,
		Title:           u.Title,
		JiraID:          u.JiraID,
		Description:     u.Description,
		AssigneeId:      u.AssigneeId,
		PriceForClosing: u.PriceForClosing,
		PriceForAssign:  u.PriceForAssign,
		CreatorId:       u.CreatorId,
		AssignedAt:      u.AssignedAt,
		ClosedAt:        u.ClosedAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

func RepoTypeToTask(u *Task) *types.Task {
	return &types.Task{
		ID:              u.ID,
		Title:           u.Title,
		JiraID:          u.JiraID,
		Description:     u.Description,
		Status:          u.Status,
		PriceForAssign:  u.PriceForAssign,
		PriceForClosing: u.PriceForClosing,
		AssigneeId:      u.AssigneeId,
		CreatorId:       u.CreatorId,
		AssignedAt:      u.AssignedAt,
		ClosedAt:        u.ClosedAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

func RepoTypesToTasks(u []*Task) []*types.Task {
	res := make([]*types.Task, 0, len(u))

	for _, user := range u {
		resUser := RepoTypeToTask(user)

		res = append(res, resUser)
	}

	return res
}

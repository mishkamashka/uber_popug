package repository

import (
	"time"
	"uber-popug/pkg/types"
)

type Task struct {
	ID              string    `gorm:"primarykey"`
	Status          string    `json:"status"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	AssigneeId      string    `json:"assignee_id"`
	PriceForAssign  uint8     `json:"price_for_assign"`
	PriceForClosing uint8     `json:"price_for_closing"`
	CreatorId       string    `json:"creator_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func TaskToRepoType(u *types.Task) *Task {
	return &Task{
		ID:              u.ID,
		Status:          u.Status,
		Name:            u.Name,
		Description:     u.Description,
		AssigneeId:      u.AssigneeId,
		PriceForClosing: u.PriceForClosing,
		PriceForAssign:  u.PriceForAssign,
		CreatorId:       u.CreatorId,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

func RepoTypeToTask(u *Task) *types.Task {
	return &types.Task{
		ID:              u.ID,
		Name:            u.Name,
		Description:     u.Description,
		Status:          u.Status,
		PriceForAssign:  u.PriceForAssign,
		PriceForClosing: u.PriceForClosing,
		AssigneeId:      u.AssigneeId,
		CreatorId:       u.CreatorId,
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

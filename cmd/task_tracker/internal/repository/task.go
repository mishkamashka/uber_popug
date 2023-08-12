package repository

import (
	"time"
	"uber-popug/pkg/types"
)

type Task struct {
	ID            string `gorm:"primarykey"`
	Status        string `json:"status"`
	AssigneeRefer string
	Assignee      User      `gorm:"foreignKey:AssigneeRefer"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func RepoTypeToTask(u *Task) *types.Task {
	return &types.Task{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
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

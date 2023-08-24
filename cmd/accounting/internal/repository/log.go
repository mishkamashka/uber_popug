package repository

import (
	"time"
	"uber-popug/pkg/types"
)

type AuditLog struct {
	ID        string    `gorm:"primarykey"`
	UserID    string    `gorm:"owner_id"`
	Amount    int       `gorm:"amount"`
	Reason    string    `gorm:"reason"`
	TaskInfo  *TaskInfo `gorm:"embedded"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

type TaskInfo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	JiraID      string `json:"jira_id"`
	Description string `json:"description"`
}

func AuditLogToRepoType(u *types.AuditLog) *AuditLog {
	log := &AuditLog{
		ID:        u.ID,
		UserID:    u.UserID,
		Reason:    u.Reason,
		Amount:    u.Amount,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	if u.TaskInfo != nil {
		log.TaskInfo = &TaskInfo{
			ID:          u.TaskInfo.ID,
			Title:       u.TaskInfo.Title,
			JiraID:      u.TaskInfo.JiraID,
			Description: u.TaskInfo.Description,
		}
	}

	return log
}

func RepoTypeToAuditLog(u *AuditLog) *types.AuditLog {
	log := &types.AuditLog{
		ID:        u.ID,
		UserID:    u.UserID,
		Amount:    u.Amount,
		Reason:    u.Reason,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	if u.TaskInfo != nil {
		log.TaskInfo = &types.TaskInfo{
			ID:          u.TaskInfo.ID,
			Title:       u.TaskInfo.Title,
			JiraID:      u.TaskInfo.JiraID,
			Description: u.TaskInfo.Description,
		}
	}
	return log
}

func RepoTypesToAuditLogs(u []*AuditLog) []*types.AuditLog {
	res := make([]*types.AuditLog, 0, len(u))

	for _, user := range u {
		resUser := RepoTypeToAuditLog(user)

		res = append(res, resUser)
	}

	return res
}

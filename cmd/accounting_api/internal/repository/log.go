package repository

import (
	"time"
	"uber-popug/pkg/types"
)

type AuditLog struct {
	ID        string    `gorm:"primarykey"`
	UserID    string    `gorm:"owner_id"`
	Amount    int32     `gorm:"amount"`
	Reason    string    `gorm:"reason"`
	TaskID    string    `gorm:"task_id"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

func AuditLogToRepoType(u *types.AuditLog) *AuditLog {
	return &AuditLog{
		ID:        u.ID,
		UserID:    u.UserID,
		Amount:    u.Amount,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func RepoTypeToAuditLog(u *AuditLog) *types.AuditLog {
	return &types.AuditLog{
		ID:        u.ID,
		UserID:    u.UserID,
		Amount:    u.Amount,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func RepoTypesToAuditLogs(u []*AuditLog) []*types.AuditLog {
	res := make([]*types.AuditLog, 0, len(u))

	for _, user := range u {
		resUser := RepoTypeToAuditLog(user)

		res = append(res, resUser)
	}

	return res
}

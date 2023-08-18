package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
	"uber-popug/pkg/types"
)

type Repository struct {
	client *gorm.DB
}

func NewRepository() (*Repository, error) {
	r := &Repository{}

	err := r.OnStart()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Repository) OnStart() error {
	dbURL := "postgres://postgres:postgres@localhost:5432/popug"

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return err
	}

	if err = db.AutoMigrate(&AuditLog{}); err != nil {
		return err
	}

	if err = db.AutoMigrate(&Balance{}); err != nil {
		return err
	}

	r.client = db

	return nil
}

func (r *Repository) CreateAuditLog(log *types.AuditLog) error {
	taskToInsert := AuditLogToRepoType(log)

	res := r.client.FirstOrCreate(taskToInsert)

	return res.Error
}

func (r *Repository) GetUserAuditLogs(userID string) ([]*types.AuditLog, error) {
	var logs []*AuditLog

	tx := r.client.Where("user_id = ?", userID).Find(&logs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return RepoTypesToAuditLogs(logs), nil
}

func (r *Repository) GetUserAuditLogsForPeriod(userID string, daysAmount int) ([]*types.AuditLog, error) {
	var logs []*AuditLog

	// TODO
	//tx := r.client.Where("user_id = ?", userID).Find(&logs)
	//if tx.Error != nil {
	//	return nil, tx.Error
	//}

	return RepoTypesToAuditLogs(logs), nil
}

func (r *Repository) GetUserAuditLogsForDay(userID string, day time.Time) ([]*types.AuditLog, error) {
	var logs []*AuditLog

	// TODO
	//tx := r.client.Where("user_id = ?", userID).Find(&logs)
	//if tx.Error != nil {
	//	return nil, tx.Error
	//}

	return RepoTypesToAuditLogs(logs), nil
}

func (r *Repository) GetTodayBalance(userID string) (*types.Balance, error) {

}

func (r *Repository) GetBalanceForDay(userID string, day time.Time) (*types.Balance, error) {

}

func (r *Repository) UpdateBalance(userID string, amount int32) (*types.Balance, error) {

}

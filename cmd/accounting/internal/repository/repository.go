package repository

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

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

func (r *Repository) GetUserAuditLogsForPeriod(userID string, from, to time.Time) ([]*types.AuditLog, error) {
	var logs []*AuditLog

	err := r.client.Where("user_id = ? and created_at > ? and created_at < ?", userID, from, to).Order("created_at DESC").Find(logs).Error
	if err != nil {
		return nil, err
	}

	return RepoTypesToAuditLogs(logs), nil
}

func (r *Repository) GetPopugBalance(userID string) (*types.Balance, error) {
	var balance *Balance

	err := r.client.Where("user_id = ? ", userID).First(balance).Error
	if err != nil {
		return nil, err
	}

	return RepoTypeToBalance(balance), nil
}

func (r *Repository) UpdatePopugBalanceByValue(userID string, amount int) error {
	balance := &Balance{
		UserID: userID,
	}

	err := r.client.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(&balance).FirstOrCreate(&balance).Error; err != nil {
			return err
		}

		if err := tx.Model(&balance).Where("user_id = ?", userID).Update("amount", balance.Amount+amount).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *Repository) GetAllNegativePopugsBalances() ([]*types.Balance, error) {
	var balances []*Balance

	if err := r.client.Where("amount < ?", 0).Find(balances).Error; err != nil {
		return nil, err
	}

	return RepoTypesToBalances(balances), nil
}

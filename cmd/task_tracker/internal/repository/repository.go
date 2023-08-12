package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	if err = db.AutoMigrate(&User{}); err != nil {
		return err
	}

	r.client = db

	return nil
}

func (r *Repository) CreateTask(task *types.Task) error {
	taskToInsert := TaskToRepoType(task)

	res := r.client.FirstOrCreate(taskToInsert)

	return res.Error
}

func (r *Repository) GetUserTasks(userID string) ([]*types.Task, error) {
	var tasks []*Task

	tx := r.client.Where("assignee_id = ?", userID).Find(&tasks)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return RepoTypesToTasks(tasks), nil
}

func (r *Repository) DeleteUserTasks(userID string) error {
	tx := r.client.Where("assignee_id = ?", userID).Delete(&Task{})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *Repository) UpdateTaskStatus(taskID, status string) (*types.Task, error) {
	task := &Task{}

	if err := r.client.Model(&task).Clauses(clause.Returning{}).Where("id = ?", taskID).Update("status", status).Error; err != nil {
		return nil, err
	}

	return RepoTypeToTask(task), nil
}

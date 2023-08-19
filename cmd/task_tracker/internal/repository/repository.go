package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	if err = db.AutoMigrate(&Task{}); err != nil {
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

func (r *Repository) DeleteTask(taskID string) error {
	tx := r.client.Where("id = ?", taskID).Delete(&Task{})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *Repository) CloseTask(taskID string) (*types.Task, error) {
	task := &Task{
		ID:       taskID,
		Status:   "closed",
		ClosedAt: time.Now(),
	}

	if err := r.client.Save(task).Clauses(clause.Returning{}).First(task).Error; err != nil {
		return nil, err
	}

	return RepoTypeToTask(task), nil
}

func (r *Repository) UpdateTask(task *types.Task) error {
	taskToSave := TaskToRepoType(task)

	return r.client.Save(taskToSave).Error
}

func (r *Repository) GetAllOpenTasks() ([]*types.Task, error) {
	var tasks []*Task

	tx := r.client.Where("status = ?", "open").Find(&tasks)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return RepoTypesToTasks(tasks), nil
}

func (r *Repository) TopTask(from time.Time) (*types.Task, error) {
	var task *Task

	tx := r.client.Where("status = ? and closed_at >= ?", "closed", from).Order("price_for_closing DESC").First(task)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return RepoTypeToTask(task), nil
}

func (r *Repository) GetAssignedTasksFromTime(from time.Time) ([]*types.Task, error) {
	var tasks []*Task

	tx := r.client.Where("assigned_at >= ?", from).Find(tasks)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return RepoTypesToTasks(tasks), nil
}

func (r *Repository) GetClosedTasksFromTime(from time.Time) ([]*types.Task, error) {
	var tasks []*Task

	tx := r.client.Where("closed_at >= ?", from).Find(tasks)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return RepoTypesToTasks(tasks), nil
}

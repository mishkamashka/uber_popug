package repository

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"uber-popug/cmd/auth_service/internal/types"
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

func (r *Repository) CreateUser(user *types.User) error {
	userToInsert, err := UserToRepoType(user)
	if err != nil {
		return fmt.Errorf("convert user to repo type: %s", err)
	}

	res := r.client.FirstOrCreate(userToInsert)

	return res.Error
}

func (r *Repository) GetUserByEmail(email string) (*types.User, error) {
	var user types.User

	tx := r.client.Where("email = ?", email).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &user, nil
}

func (r *Repository) GetUsersByRole(role string) ([]*types.User, error) {
	var users []*types.User

	tx := r.client.Where("role = ?", role).Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return users, nil
}

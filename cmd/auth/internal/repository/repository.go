package repository

import (
	"errors"
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

func (r *Repository) CreateUser(user *types.User) error {
	userToInsert := UserToRepoType(user)

	res := r.client.FirstOrCreate(userToInsert)

	return res.Error
}

func (r *Repository) GetUserByEmail(email string) (*types.User, error) {
	var user *User

	tx := r.client.Where("email = ?", email).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if user == nil {
		return nil, errors.New("user with email " + email + "not found")
	}

	return RepoTypeToUser(user), nil
}

func (r *Repository) GetUsersByRole(role string) ([]*types.User, error) {
	var users []*User

	err := r.client.Where("role = ?", role).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return RepoTypesToUsers(users), nil
}

func (r *Repository) UpdateUserRole(email, role string) (*types.User, error) {
	user := &User{}

	if err := r.client.Model(&user).Clauses(clause.Returning{}).Where("email = ?", email).Update("role", role).Error; err != nil {
		return nil, err
	}

	return RepoTypeToUser(user), nil
}

func (r *Repository) GetAllPopugsIDs() ([]string, error) {
	var popugs []*types.User

	err := r.client.Select("id").Where("role = ?", "popug").Find(&popugs).Error
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(popugs))
	for _, p := range popugs {
		res = append(res, p.ID)
	}

	return res, nil
}

func (r *Repository) DeleteUser(email string) (*types.User, error) {
	var user *User

	err := r.client.Clauses(clause.Returning{}).Where("email = ?", email).Delete(&user).Error
	if err != nil {
		return nil, err
	}

	return RepoTypeToUser(user), nil
}

func (r *Repository) GetUsersEmail(userID string) (string, error) {
	var user *User

	err := r.client.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return "", err
	}

	return user.Email, nil
}

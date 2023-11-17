package repository

import (
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	FindByID(id int32) (*model.User, error)
	SaveUser(user *model.User) (*model.User, error)
	ExistsByFirstNameAndLastName(user *model.User) bool
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) FindByID(id int32) (*model.User, error) {
	var user model.User
	result := repo.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (repo *UserRepository) SaveUser(user *model.User) (*model.User, error) {
	result := repo.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *UserRepository) ExistsByFirstNameAndLastName(user *model.User) bool {
	var count int64
	tx := repo.db.Model(&model.User{}).Where("first_name = ? AND last_name = ?", user.FirstName, user.LastName).Count(&count)
	if tx.Error != nil {
		panic(tx.Error)
	}
	return count > 0
}

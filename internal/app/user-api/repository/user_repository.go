package repository

import (
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	FindByID(id *uuid.UUID) (*model.User, error)
	SaveUser(user *model.User) (*model.User, error)
	ExistsByFirstNameAndLastName(user *model.User) bool
}

type DB interface {
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Create(value interface{}) (tx *gorm.DB)
	Model(value interface{}) (tx *gorm.DB)
	Where(query interface{}, args ...interface{}) (tx *gorm.DB)
}

type UserRepository struct {
	Db DB
}

func (repo *UserRepository) FindByID(id *uuid.UUID) (*model.User, error) {
	var user model.User
	result := repo.Db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (repo *UserRepository) SaveUser(user *model.User) (*model.User, error) {
	result := repo.Db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *UserRepository) ExistsByFirstNameAndLastName(user *model.User) bool {
	var count int64
	tx := repo.Db.Model(&model.User{}).Where("first_name = ? AND last_name = ?", user.FirstName, user.LastName).Count(&count)
	if tx.Error != nil {
		panic(tx.Error)
	}
	return count > 0
}

package repository_test

import (
	"errors"
	"testing"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/repository"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUserRepository_FindByID(t *testing.T) {
	mockDB := new(mocks.DB)
	repo := repository.UserRepository{Db: mockDB}

	t.Run("Return error when no user in request", func(t *testing.T) {
		user := new(model.User)

		gormDB := new(gorm.DB)
		gormDB.Error = errors.New("record not found")

		mockDB.On("First", user, &user.ID).Return(gormDB)

		_, err := repo.FindByID(&user.ID)
		assert.EqualError(t, err, "record not found")
	})

	t.Run("Return error with valid UUID but no user", func(t *testing.T) {
		user := new(model.User)
		userID := uuid.New()

		gormDB := new(gorm.DB)
		gormDB.Error = errors.New("record not found")

		mockDB.On("First", user, &userID).Return(gormDB)

		_, err := repo.FindByID(&userID)
		assert.EqualError(t, err, "record not found")
	})

	t.Run("Return user that exists in the repo", func(t *testing.T) {
		user := new(model.User)
		userID := uuid.New()

		firstName := "John"
		lastName := "Doe"
		email := "john.doe@example.com"
		age := int8(18)

		gormDB := new(gorm.DB)
		gormDB.Error = nil

		mockDB.On("First", user, &userID).Run(func(args mock.Arguments) {
			userArg := args.Get(0).(*model.User)
			userArg.FirstName = firstName
			userArg.LastName = lastName
			userArg.Email = email
			userArg.Age = age
		}).Return(gormDB)

		repoUser, err := repo.FindByID(&userID)
		assert.NoError(t, err)
		assert.Equal(t, firstName, repoUser.FirstName)
		assert.Equal(t, lastName, repoUser.LastName)
		assert.Equal(t, email, repoUser.Email)
		assert.Equal(t, age, repoUser.Age)
	})
}

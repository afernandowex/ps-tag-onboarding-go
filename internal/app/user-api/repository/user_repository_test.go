package repository_test

import (
	"errors"
	"testing"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Create(value interface{}) (tx *gorm.DB) {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Model(value interface{}) *gorm.DB {
	m.Called(value)
	return &gorm.DB{}
}

func (m *MockDB) Where(query interface{}, arguments ...interface{}) *gorm.DB {
	m.Called(query, arguments)
	return &gorm.DB{}
}

func (m *MockDB) Count(count *int64) *gorm.DB {
	m.Called(count)
	*count = 1
	return &gorm.DB{}
}

func TestUserRepository_FindByID(t *testing.T) {
	mockDB := new(MockDB)
	repo := repository.UserRepository{Db: mockDB}

	t.Run("Return error when no user in request", func(t *testing.T) {
		user := new(model.User)

		gormDB := new(gorm.DB)
		gormDB.Error = errors.New("record not found")

		mockDB.On("First", user, []interface{}{&user.ID}).Return(gormDB)

		_, err := repo.FindByID(&user.ID)
		assert.EqualError(t, err, "record not found")
	})

	t.Run("Return error with valid UUID but no user", func(t *testing.T) {
		user := new(model.User)
		userID := uuid.New()

		gormDB := new(gorm.DB)
		gormDB.Error = errors.New("record not found")

		mockDB.On("First", user, []interface{}{&userID}).Return(gormDB)

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

		mockDB.On("First", user, []interface{}{&userID}).Run(func(args mock.Arguments) {
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

// func TestUserRepository_ExistsByFirstAndLastName(t *testing.T) {
// 	mockDB := new(MockDB)
// 	repo := repository.UserRepository{Db: mockDB}

// 	user := &model.User{
// 		FirstName: "John",
// 		LastName:  "Doe",
// 		Email:     "john.doe@gmail.com",
// 		Age:       18,
// 	}

// 	mockDB.On("Model", &model.User{}).Return(mockDB)
// 	mockDB.On("Where", "first_name = ? AND last_name = ?", user.FirstName, user.LastName).Return(mockDB)
// 	mockDB.On("Count", mock.Anything).Run(func(args mock.Arguments) {
// 		countArg := args.Get(0).(*int64)
// 		*countArg = 1
// 	}).Return(mockDB)
// 	repo.ExistsByFirstNameAndLastName(user)
// }

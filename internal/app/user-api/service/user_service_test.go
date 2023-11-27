package service

import (
	"errors"
	"strconv"
	"testing"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/errormessage"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/model"
	"github.com/stretchr/testify/assert"
)

type MockUserRepository struct {
	Users map[int32]*model.User
}

func (repo *MockUserRepository) FindByID(id int32) (*model.User, error) {
	user, ok := repo.Users[id]
	if !ok {
		return nil, errors.New("User not found")
	}
	return user, nil
}

func (repo *MockUserRepository) SaveUser(user *model.User) (*model.User, error) {
	id := int32(len(repo.Users) + 1)
	user.ID = id
	repo.Users[id] = user
	return user, nil
}

func (repo *MockUserRepository) ExistsByFirstNameAndLastName(user *model.User) bool {
	return false
}

type MockUserValidationService struct {
	ValidationErrors map[*model.User][]string
}

func (validator *MockUserValidationService) ValidateUser(user *model.User) []string {
	return validator.ValidationErrors[user]
}

func (validator *MockUserValidationService) ValidateUserID(ID *string) []string {
	return nil
}

func TestFindUserByID(t *testing.T) {
	repo := &MockUserRepository{
		Users: make(map[int32]*model.User),
	}
	validator := &MockUserValidationService{}
	service := &UserService{
		Repository: repo,
		Validator:  validator,
	}

	t.Run("Return user when valid ID", func(t *testing.T) {
		user := &model.User{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       25,
		}
		repo.Users[user.ID] = user

		ID := strconv.Itoa(int(user.ID))
		result, err := service.FindUserByID(&ID)
		assert.Nil(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("Return error when invalid ID", func(t *testing.T) {
		ID := "invalid"
		result, err := service.FindUserByID(&ID)
		assert.Nil(t, result)
		assert.Equal(t, &errormessage.ErrorMessage{ErrorMessageText: "User not found", ErrorStatus: 404}, err)
	})
}

func TestSaveUser(t *testing.T) {
	repo := &MockUserRepository{
		Users: make(map[int32]*model.User),
	}
	validator := &MockUserValidationService{}
	service := &UserService{
		Repository: repo,
		Validator:  validator,
	}

	t.Run("Save user when valid user", func(t *testing.T) {
		user := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       25,
		}

		result, err := service.SaveUser(user)
		assert.Nil(t, err)
		assert.Equal(t, user, result)
		assert.Equal(t, int32(1), user.ID)
	})

	t.Run("Return error when invalid user", func(t *testing.T) {
		user := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "invalidemail",
			Age:       16,
		}
		validator.ValidationErrors = map[*model.User][]string{
			user: {"Email is not valid", "Age must be at least 18"},
		}

		result, err := service.SaveUser(user)
		assert.Nil(t, result)
		assert.Equal(t, &errormessage.ErrorMessage{ErrorMessageText: "Email is not valid, Age must be at least 18", ErrorStatus: 400}, err)
	})
}

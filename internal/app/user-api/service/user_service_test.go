package service

import (
	"errors"
	"testing"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/errormessage"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockUserRepository struct {
	Users map[uuid.UUID]*model.User
}

func (repo *MockUserRepository) FindByID(id uuid.UUID) (*model.User, error) {
	user, ok := repo.Users[id]
	if !ok {
		return nil, errors.New("User not found")
	}
	return user, nil
}

func (repo *MockUserRepository) SaveUser(user *model.User) (*model.User, error) {
	user.ID = uuid.New()
	repo.Users[user.ID] = user
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
		Users: make(map[uuid.UUID]*model.User),
	}
	validator := &MockUserValidationService{}
	service := &UserService{
		Repository: repo,
		Validator:  validator,
	}

	t.Run("Return user when valid ID", func(t *testing.T) {
		user := &model.User{
			ID:        uuid.New(),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       25,
		}
		repo.Users[user.ID] = user

		userID := user.ID.String()

		result, err := service.FindUserByID(&userID)
		assert.Nil(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("Return error when invalid UUID", func(t *testing.T) {
		ID := "invalid"
		result, err := service.FindUserByID(&ID)
		assert.Nil(t, result)
		assert.Equal(t, &errormessage.ErrorMessage{ErrorMessageText: "Invalid user ID.", ErrorStatus: 400}, err)
	})

	t.Run("Return error when not found UUID", func(t *testing.T) {
		ID := uuid.New().String()
		result, err := service.FindUserByID(&ID)
		assert.Nil(t, result)
		assert.Equal(t, &errormessage.ErrorMessage{ErrorMessageText: "User not found", ErrorStatus: 404}, err)
	})
}

func TestSaveUser(t *testing.T) {
	repo := &MockUserRepository{
		Users: make(map[uuid.UUID]*model.User),
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
		id, guidErr := uuid.Parse(result.ID.String())
		if guidErr != nil {
			assert.Fail(t, "Invalid UUID %s", id)
		}
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

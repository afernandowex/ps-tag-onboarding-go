package service

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/repository/mocks"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/constant"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/errormessage"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/validation/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func initialiseService() (*mocks.IUserRepository, *mock.IUserValidationService, *UserService) {
	repo := new(mocks.IUserRepository)
	validator := new(mock.IUserValidationService)
	service := &UserService{
		Repository: repo,
		Validator:  validator,
	}
	return repo, validator, service
}

func TestFindUserByID(t *testing.T) {

	t.Run("Return user when valid ID", func(t *testing.T) {
		repo, validator, service := initialiseService()
		user := &model.User{
			ID:        uuid.New(),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       25,
		}

		userID := user.ID.String()

		repo.On("FindByID", &user.ID).Return(user, nil)
		validator.On("ValidateUserID", &userID).Return([]string{})

		result, err := service.FindUserByID(&userID)
		assert.Nil(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("Validator test - Return error when invalid UUID", func(t *testing.T) {
		_, validator, service := initialiseService()
		ID := "invalid"
		validator.On("ValidateUserID", &ID).Return([]string{"Invalid user ID."})
		result, err := service.FindUserByID(&ID)
		assert.Nil(t, result)
		assert.Equal(t, &errormessage.ErrorMessage{ErrorMessageText: "Invalid user ID.", ErrorStatus: 400}, err)
	})

	t.Run("Repository test - Return error when not found UUID", func(t *testing.T) {
		repo, validator, service := initialiseService()
		userID := uuid.New()
		userIDStr := userID.String()
		validator.On("ValidateUserID", &userIDStr).Return([]string{})
		repo.On("FindByID", &userID).Return(nil, errors.New(constant.ErrorUserNotFound))

		result, err := service.FindUserByID(&userIDStr)

		assert.Nil(t, result)
		expectedError := errormessage.ErrorMessage{ErrorMessageText: constant.ErrorUserNotFound, ErrorStatus: http.StatusNotFound}
		assert.Equal(t, &expectedError, err)
	})
}

func TestSaveUser(t *testing.T) {

	t.Run("Save user when valid user", func(t *testing.T) {
		repo, validator, service := initialiseService()
		user := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       25,
		}

		repo.On("SaveUser", user).Return(user, nil)
		validator.On("ValidateUser", user).Return([]string{})
		result, err := service.SaveUser(user)
		assert.Nil(t, err)
		assert.Equal(t, user, result)
		id, guidErr := uuid.Parse(result.ID.String())
		if guidErr != nil {
			assert.Fail(t, "Invalid UUID %s", id)
		}
	})

	t.Run("Return error when invalid user", func(t *testing.T) {
		_, validator, service := initialiseService()
		user := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "invalidemail",
			Age:       16,
		}
		validator.On("ValidateUser", user).Return([]string{constant.ErrorEmailFormat, constant.ErrorAgeMinimum})

		result, err := service.SaveUser(user)
		assert.Nil(t, result)
		assert.Equal(t, &errormessage.ErrorMessage{ErrorMessageText: fmt.Sprintf("%s, %s", constant.ErrorEmailFormat, constant.ErrorAgeMinimum), ErrorStatus: 400}, err)
	})

	t.Run("Return error when name already exists", func(t *testing.T) {
		repo, validator, service := initialiseService()
		user := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       25,
		}

		repo.On("SaveUser", user).Return(nil, errors.New(constant.ErrorNameAlreadyExists))
		validator.On("ValidateUser", user).Return([]string{})

		result, err := service.SaveUser(user)
		assert.Nil(t, result)
		assert.Equal(t, &errormessage.ErrorMessage{ErrorMessageText: constant.ErrorNameAlreadyExists, ErrorStatus: 400}, err)
	})

	t.Run("Return error when repository returns an error", func(t *testing.T) {
		repo, validator, service := initialiseService()
		user := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       25,
		}

		repo.On("SaveUser", user).Return(nil, errors.New("repository connection error"))
		validator.On("ValidateUser", user).Return([]string{})

		result, err := service.SaveUser(user)
		assert.Nil(t, result)
		assert.Equal(t, &errormessage.ErrorMessage{ErrorMessageText: "repository connection error", ErrorStatus: 500}, err)
	})
}

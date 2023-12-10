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

func TestFindUserByID(t *testing.T) {
	repo := new(mocks.IUserRepository)
	validator := new(mock.IUserValidationService)
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

		userID := user.ID.String()

		repo.On("FindByID", &user.ID).Return(user, nil)
		validator.On("ValidateUserID", &userID).Return([]string{})

		result, err := service.FindUserByID(&userID)
		assert.Nil(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("Validator test - Return error when invalid UUID", func(t *testing.T) {
		ID := "invalid"
		validator.On("ValidateUserID", &ID).Return([]string{"Invalid user ID."})
		result, err := service.FindUserByID(&ID)
		assert.Nil(t, result)
		assert.Equal(t, &errormessage.ErrorMessage{ErrorMessageText: "Invalid user ID.", ErrorStatus: 400}, err)
	})

	t.Run("Repository test - Return error when not found UUID", func(t *testing.T) {
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
	repo := new(mocks.IUserRepository)
	validator := new(mock.IUserValidationService)

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
}

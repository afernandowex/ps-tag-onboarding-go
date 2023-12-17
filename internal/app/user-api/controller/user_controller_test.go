package controller_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/constant"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/controller"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/controller/mocks"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/errormessage"
	serviceMocks "github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/service/mocks"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserController_FindUser(t *testing.T) {

	t.Run("Find user success", func(t *testing.T) {
		mockService := new(serviceMocks.IUserService)
		contextMock := new(mocks.Context)
		responseWriterMock := new(mocks.ResponseWriter)
		e := echo.New()
		response := echo.NewResponse(responseWriterMock, e)

		controller := controller.UserController{Service: mockService}

		savedUser := &model.User{
			ID:        uuid.New(),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       18,
		}

		var nilErr *errormessage.ErrorMessage = nil
		savedUserStr := savedUser.ID.String()
		mockService.On("FindUserByID", &savedUserStr).Return(savedUser, nilErr)
		contextMock.On("Param", "id").Return(savedUserStr)
		responseWriterMock.On("Header").Return(http.Header{})
		responseWriterMock.On("WriteHeader", http.StatusOK).Return(nil)
		contextMock.On("Response").Return(response)
		contextMock.On("JSON", http.StatusOK, savedUser).Return(nil)
		err := controller.FindUser(contextMock)

		assert.True(t, contextMock.AssertCalled(t, "JSON", http.StatusOK, savedUser))
		assert.NoError(t, err)
	})

	t.Run("Find user send back error message if generated in service", func(t *testing.T) {
		mockService := new(serviceMocks.IUserService)
		contextMock := new(mocks.Context)
		responseWriterMock := new(mocks.ResponseWriter)
		e := echo.New()
		response := echo.NewResponse(responseWriterMock, e)

		controller := controller.UserController{Service: mockService}

		errorMessage := errormessage.NewErrorMessage(constant.ErrorInvalidUserID, http.StatusBadRequest)
		mockService.On("FindUserByID", mock.Anything).Return(nil, &errorMessage)
		contextMock.On("Param", "id").Return("invalid")
		responseWriterMock.On("Header").Return(http.Header{})
		contextMock.On("Response").Return(response)
		contextMock.On("JSON", http.StatusBadRequest, &errorMessage).Return(nil)

		err := controller.FindUser(contextMock)
		assert.NoError(t, err)
		assert.True(t, contextMock.AssertCalled(t, "JSON", http.StatusBadRequest, &errorMessage))
	})
}

func TestUserController_SaveUser(t *testing.T) {

	t.Run("Save user success", func(t *testing.T) {
		mockService := new(serviceMocks.IUserService)
		contextMock := new(mocks.Context)
		responseWriterMock := new(mocks.ResponseWriter)
		e := echo.New()
		response := echo.NewResponse(responseWriterMock, e)

		controller := controller.UserController{Service: mockService}

		savedUser := &model.User{
			ID:        uuid.New(),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@gmail.com",
			Age:       18,
		}

		contextMock.On("Response").Return(response)
		responseWriterMock.On("Header").Return(http.Header{})
		responseWriterMock.On("WriteHeader", http.StatusOK).Return(nil)
		contextMock.On("Bind", mock.Anything).Return(nil)
		mockService.On("SaveUser", mock.Anything).Return(savedUser, nil)
		contextMock.On("JSON", http.StatusOK, savedUser).Return(nil)

		err := controller.SaveUser(contextMock)
		assert.NoError(t, err)
		assert.True(t, contextMock.AssertCalled(t, "JSON", http.StatusOK, savedUser))
	})

	t.Run("Save user error while sending malformed JSON", func(t *testing.T) {
		mockService := new(serviceMocks.IUserService)
		contextMock := new(mocks.Context)
		responseWriterMock := new(mocks.ResponseWriter)
		e := echo.New()
		response := echo.NewResponse(responseWriterMock, e)

		controller := controller.UserController{Service: mockService}

		contextMock.On("Response").Return(response)
		responseWriterMock.On("Header").Return(http.Header{})
		contextMock.On("Bind", mock.Anything).Return(errors.New("error while binding JSON"))
		contextMock.On("JSON", http.StatusBadRequest, errormessage.NewErrorMessage(constant.ErrorInvalidUserObject, http.StatusBadRequest)).Return(nil)

		err := controller.SaveUser(contextMock)

		assert.NoError(t, err)
		assert.True(t, contextMock.AssertCalled(t, "JSON", http.StatusBadRequest, errormessage.NewErrorMessage(constant.ErrorInvalidUserObject, http.StatusBadRequest)))
	})

	t.Run("Save user return error while saving user", func(t *testing.T) {
		mockService := new(serviceMocks.IUserService)
		contextMock := new(mocks.Context)
		responseWriterMock := new(mocks.ResponseWriter)
		e := echo.New()
		response := echo.NewResponse(responseWriterMock, e)
		controller := controller.UserController{Service: mockService}

		contextMock.On("Response").Return(response)
		responseWriterMock.On("Header").Return(http.Header{})
		contextMock.On("Bind", mock.Anything).Return(nil)
		errorMessage := errormessage.NewErrorMessage(constant.ErrorAgeMinimum, http.StatusBadRequest)
		mockService.On("SaveUser", mock.Anything).Return(nil, &errorMessage)
		contextMock.On("JSON", http.StatusBadRequest, &errorMessage).Return(nil)

		err := controller.SaveUser(contextMock)

		assert.NoError(t, err)
		assert.True(t, contextMock.AssertCalled(t, "JSON", http.StatusBadRequest, &errorMessage))
	})
}

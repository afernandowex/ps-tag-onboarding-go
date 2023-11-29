package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/constant"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/controller"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/errormessage"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) FindUserByID(ID *string) (*model.User, *errormessage.ErrorMessage) {
	args := m.Called(ID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errormessage.ErrorMessage)
	}
	return args.Get(0).(*model.User), args.Get(1).(*errormessage.ErrorMessage)
}

func (m *MockUserService) SaveUser(user *model.User) (*model.User, *errormessage.ErrorMessage) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errormessage.ErrorMessage)
	}
	return args.Get(0).(*model.User), args.Get(1).(*errormessage.ErrorMessage)
}

func TestUserController_FindUser(t *testing.T) {

	t.Run("Find user success", func(t *testing.T) {
		mockService := &MockUserService{}
		controller := controller.UserController{Service: mockService}

		savedUser := &model.User{
			ID:        uuid.New(),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       18,
		}

		var nilErr *errormessage.ErrorMessage = nil
		mockService.On("FindUserByID", mock.Anything).Return(savedUser, nilErr)
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/users/"+savedUser.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := controller.FindUser(c)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, err)
		assert.Equal(t, http.Header{echo.HeaderContentType: []string{echo.MIMEApplicationJSON}}, rec.Header())

		var responseUser model.User
		err = json.Unmarshal(rec.Body.Bytes(), &responseUser)
		assert.NoError(t, err)
		assert.Equal(t, savedUser, &responseUser)
	})

	t.Run("Find user send back error message if generated in service", func(t *testing.T) {
		mockService := &MockUserService{}
		controller := controller.UserController{Service: mockService}

		errorMessage := errormessage.NewErrorMessage(constant.ErrorInvalidUserID, http.StatusBadRequest)
		mockService.On("FindUserByID", mock.Anything).Return(nil, &errorMessage)
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/users/invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := controller.FindUser(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NoError(t, err)

		assert.Equal(t, http.Header{echo.HeaderContentType: []string{echo.MIMEApplicationJSON}}, rec.Header())
		assert.Equal(t, "{\"message\":\"Invalid user ID.\",\"status\":400}\n", string(rec.Body.Bytes()))
	})
}

func TestUserController_SaveUser(t *testing.T) {

	t.Run("Save user success", func(t *testing.T) {
		mockService := &MockUserService{}
		controller := controller.UserController{Service: mockService}

		user := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       18,
		}

		savedUser := &model.User{
			ID:        uuid.New(),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       18,
		}

		var nilErr *errormessage.ErrorMessage = nil
		mockService.On("SaveUser", user).Return(savedUser, nilErr)

		e := echo.New()
		reqBody, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := controller.SaveUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, http.Header{echo.HeaderContentType: []string{echo.MIMEApplicationJSON}}, rec.Header())

		var responseUser model.User
		err = json.Unmarshal(rec.Body.Bytes(), &responseUser)
		assert.NoError(t, err)
		assert.Equal(t, savedUser, &responseUser)
	})

	t.Run("Save user error while sending malformed JSON", func(t *testing.T) {
		mockService := &MockUserService{}
		controller := controller.UserController{Service: mockService}

		e := echo.New()

		malformedJSON := `{"name": "John Doe", "email": "johndoe@example.com", "age": 30,` // missing closing bracket
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(malformedJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := controller.SaveUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, http.Header{echo.HeaderContentType: []string{echo.MIMEApplicationJSON}}, rec.Header())
		assert.Equal(t, "{\"message\":\"User object cannot be unmarshalled\",\"status\":400}\n", string(rec.Body.Bytes()))
	})

	t.Run("Save user return error while saving user", func(t *testing.T) {
		mockService := &MockUserService{}
		controller := controller.UserController{Service: mockService}

		var saveErr errormessage.ErrorMessage = errormessage.NewErrorMessage(constant.ErrorInvalidUserID, http.StatusBadRequest)
		user := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       18,
		}
		mockService.On("SaveUser", user).Return(user, &saveErr)
		e := echo.New()

		reqBody, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := controller.SaveUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, http.Header{echo.HeaderContentType: []string{echo.MIMEApplicationJSON}}, rec.Header())
		assert.Equal(t, "{\"message\":\"Invalid user ID.\",\"status\":400}\n", string(rec.Body.Bytes()))
	})
}

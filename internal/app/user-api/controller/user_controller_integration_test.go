package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/errors"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/repository"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/service"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/validation"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestListUsers(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	var repo repository.IUserRepository = &repository.UserRepository{Db: db}
	var validator validation.IUserValidationService = &validation.UserValidationService{UserRepository: repo}
	var service service.IUserService = &service.UserService{Repository: repo, Validator: validator}
	var controller IUserController = &UserController{Service: service}

	t.Run("Return user not found when invalid user", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/find/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		controller.FindUser(c)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var errorMessage errors.ErrorMessage
		err := json.Unmarshal(rec.Body.Bytes(), &errorMessage)
		if err != nil {
			t.Fatalf("failed to unmarshal error response: %v", err)
		}
		assert.Equal(t, "User not found", errorMessage.ErrorMessageText)
	})

	t.Run("Add test user 1", func(t *testing.T) {
		e := echo.New()
		userOne := model.User{
			FirstName: "WexFirstName",
			LastName:  "WexLastName",
			Email:     "wexfirstname.wexlastname@wexinc.com",
			Age:       18,
		}
		requestBody, err := json.Marshal(userOne)
		if err != nil {
			t.Fatalf("Error marshaling JSON: %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/save", bytes.NewReader(requestBody))
		headers := make(http.Header)
		headers.Set("Content-Type", "application/json")
		req.Header = headers

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		controller.SaveUser(c)
		assert.Equal(t, http.StatusOK, rec.Code)

		var responseUser model.User
		error := json.Unmarshal(rec.Body.Bytes(), &responseUser)
		if error != nil {
			t.Fatalf("failed to unmarshal error response: %v", err)
		}
		assert.Equal(t, int32(1), responseUser.ID)
	})

	t.Run("Add invalid test user 2", func(t *testing.T) {
		e := echo.New()
		userOne := model.User{
			FirstName: "WexFirst2Name",
			LastName:  "WexLast2Name",
			Email:     "wexfirstname.wexlastnamewexinc.com", // No @ symbol in email
			Age:       16,                                   // Below minimum age
		}
		requestBody, err := json.Marshal(userOne)
		if err != nil {
			t.Fatalf("Error marshaling JSON: %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/save", bytes.NewReader(requestBody))
		headers := make(http.Header)
		headers.Set("Content-Type", "application/json")
		req.Header = headers

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		controller.SaveUser(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var errorMessage errors.ErrorMessage
		error := json.Unmarshal(rec.Body.Bytes(), &errorMessage)
		if error != nil {
			t.Fatalf("failed to unmarshal error response: %v", err)
		}
		assert.Equal(t, "Email is not valid, Age must be at least 18", errorMessage.ErrorMessageText)
		assert.Equal(t, http.StatusBadRequest, errorMessage.ErrorStatus)
	})
}

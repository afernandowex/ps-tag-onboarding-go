package controller_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/repository"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/controller"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/errormessage"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/routing"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/service"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/validation"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testCleanUp(db *gorm.DB) {
	// Reset DB between tests.
	err := db.Migrator().DropTable(&model.User{})
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
}

func setupServerAndDB(t *testing.T) (*echo.Echo, *gorm.DB, controller.IUserController, repository.IUserRepository) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	var repo repository.IUserRepository = &repository.UserRepository{Db: db}
	var validator validation.IUserValidationService = &validation.UserValidationService{}
	var service service.IUserService = &service.UserService{Repository: repo, Validator: validator}
	var controller controller.IUserController = &controller.UserController{Service: service}
	var routes = routing.Routes{UserController: controller}

	e := echo.New()
	routes.InitializeRoutes(e)

	return e, db, controller, repo
}

func TestListUsers(t *testing.T) {

	t.Run("Return user not found when invalid user", func(t *testing.T) {
		e, db, _, _ := setupServerAndDB(t)
		req := httptest.NewRequest(http.MethodGet, "/find/1", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var errorMessage errormessage.ErrorMessage
		err := json.Unmarshal(rec.Body.Bytes(), &errorMessage)
		if err != nil {
			t.Fatalf("failed to unmarshal error response: %v", err)
		}
		assert.Equal(t, "Invalid user ID.", errorMessage.ErrorMessageText)
		t.Cleanup(func() {
			testCleanUp(db)
		})
	})

	t.Run("Return user when valid user", func(t *testing.T) {

		e, db, _, repo := setupServerAndDB(t)

		user := model.User{FirstName: "WexFirst", LastName: "WexLast", Email: "wexfirst.wexlast@wexinc.com", Age: 18}
		repo.SaveUser(&user)

		url := fmt.Sprintf("/find/%s", user.ID.String())
		req := httptest.NewRequest(http.MethodGet, url, nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var responseUser model.User
		err := json.Unmarshal(rec.Body.Bytes(), &responseUser)
		if err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		assert.Equal(t, user.FirstName, responseUser.FirstName)
		assert.Equal(t, user.LastName, responseUser.LastName)
		assert.Equal(t, user.Email, responseUser.Email)
		assert.Equal(t, user.Age, responseUser.Age)
		t.Cleanup(func() {
			testCleanUp(db)
		})
	})

	t.Run("Add test user 1", func(t *testing.T) {
		e, db, controller, _ := setupServerAndDB(t)

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
		assert.Equal(t, userOne.FirstName, responseUser.FirstName)
		t.Cleanup(func() {
			testCleanUp(db)
		})
	})

	t.Run("Add invalid test user 2", func(t *testing.T) {
		e, db, controller, _ := setupServerAndDB(t)

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

		var errorMessage errormessage.ErrorMessage
		error := json.Unmarshal(rec.Body.Bytes(), &errorMessage)
		if error != nil {
			t.Fatalf("failed to unmarshal error response: %v", err)
		}
		assert.Equal(t, "Age must be at least 18, Email is not valid", errorMessage.ErrorMessageText)
		assert.Equal(t, http.StatusBadRequest, errorMessage.ErrorStatus)
		t.Cleanup(func() {
			testCleanUp(db)
		})
	})
}

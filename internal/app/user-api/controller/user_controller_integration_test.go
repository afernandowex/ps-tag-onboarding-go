package controller_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/controller"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/errormessage"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/repository"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/routing"
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
	var validator validation.IUserValidationService = &validation.UserValidationService{}
	var service service.IUserService = &service.UserService{Repository: repo, Validator: validator}
	var controller controller.IUserController = &controller.UserController{Service: service}
	var routes = routing.Routes{Controller: controller}

	t.Cleanup(func() {
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
	})

	t.Run("Return user not found when invalid user", func(t *testing.T) {
		e := echo.New()
		routes.InitializeRoutes(e)
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
	})

	t.Run("Return user when valid user", func(t *testing.T) {
		user := model.User{FirstName: "WexFirst", LastName: "WexLast", Email: "wexfirst.wexlast@wexinc.com", Age: 18}
		repo.SaveUser(&user)

		userNew, erro := repo.FindByID(&user.ID)
		fmt.Println("Here")
		fmt.Println(userNew)
		fmt.Println(erro)

		e := echo.New()
		routes.InitializeRoutes(e)
		url := fmt.Sprintf("/find/%s", user.ID.String())
		fmt.Println(url)
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
	})

	t.Run("Add test user 1", func(t *testing.T) {
		e := echo.New()
		routes.InitializeRoutes(e)

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
	})

	t.Run("Add invalid test user 2", func(t *testing.T) {
		e := echo.New()
		routes.InitializeRoutes(e)

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
	})
}

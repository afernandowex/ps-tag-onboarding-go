package service_test

import (
	"net/http"
	"testing"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/repository"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/service"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/validation"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserService(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	repo := repository.NewUserRepository(db)

	t.Run("SaveUser", func(t *testing.T) {
		validator := validation.NewUserValidationService(repo)
		service := service.NewUserService(repo, validator)

		// Test with a valid user
		user := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@example.com",
			Age:       18,
		}
		savedUser, err := service.SaveUser(user)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if savedUser == nil {
			t.Errorf("Expected saved user to not be nil")
		}
		assert.Equal(t, int32(1), savedUser.ID)

		// Test with an invalid user
		user = &model.User{
			FirstName: "",
			LastName:  "",
			Email:     "",
		}
		savedUser, err = service.SaveUser(user)
		if err == nil {
			t.Errorf("Expected error to not be nil")
		}
		if savedUser != nil {
			t.Errorf("Expected saved user to be nil")
		}
	})

	t.Run("FindUserByID", func(t *testing.T) {
		validator := validation.NewUserValidationService(repo)
		service := service.NewUserService(repo, validator)

		userIdOne := "1"
		// Test with a valid user ID
		user, err := service.FindUserByID(&userIdOne)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if user == nil {
			t.Errorf("Expected user to not be nil")
		}

		userIdZero := "0"
		// Test with an invalid user ID
		user, err = service.FindUserByID(&userIdZero)
		if err == nil {
			t.Errorf("Expected error to not be nil")
		} else {
			assert.Equal(t, "User not found", err.Message())
			assert.Equal(t, http.StatusNotFound, err.HttpStatus())
		}
		if user != nil {
			t.Errorf("Expected user to be nil")
		}
	})
}

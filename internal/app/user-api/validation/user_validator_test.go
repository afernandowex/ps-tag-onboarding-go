package validation

import (
	"testing"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/repository"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/constant"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserValidationService_ValidateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	var userRepository repository.IUserRepository = &repository.UserRepository{Db: db}
	var validator IUserValidationService = &UserValidationService{}

	testUserOne := model.User{FirstName: "Wex", LastName: "User", Email: "wexuser@wexinc.com", Age: 18}

	userRepository.SaveUser(&testUserOne)

	tests := []struct {
		name     string
		user     *model.User
		expected []string
	}{
		{
			name: "valid user",
			user: &model.User{
				FirstName: "John",
				LastName:  "Doe",
				Age:       30,
				Email:     "johndoe@example.com",
			},
			expected: nil,
		},
		{
			name: "invalid first name",
			user: &model.User{
				FirstName: "",
				LastName:  "Doe",
				Age:       30,
				Email:     "johndoe@example.com",
			},
			expected: []string{constant.ErrorNameRequired},
		},
		{
			name: "invalid last name",
			user: &model.User{
				FirstName: "John",
				LastName:  "",
				Age:       30,
				Email:     "johndoe@example.com",
			},
			expected: []string{constant.ErrorNameRequired},
		},
		{
			name: "invalid age",
			user: &model.User{
				FirstName: "John",
				LastName:  "Doe",
				Age:       0,
				Email:     "johndoe@example.com",
			},
			expected: []string{constant.ErrorAgeMinimum},
		},
		{
			name: "invalid email",
			user: &model.User{
				FirstName: "John",
				LastName:  "Doe",
				Age:       0,
				Email:     "johndoe@example.com",
			},
			expected: []string{constant.ErrorAgeMinimum},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := validator.ValidateUser(tt.user)
			if len(actual) != len(tt.expected) {
				t.Errorf("Expected %d errors, but got %d", len(tt.expected), len(actual))
			}
			for i, err := range actual {
				if err != tt.expected[i] {
					t.Errorf("Expected error '%s', but got '%s'", tt.expected[i], err)
				}
			}
		})
	}
}

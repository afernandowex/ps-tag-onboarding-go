package repository_test

import (
	"testing"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	var repo repository.IUserRepository = &repository.UserRepository{Db: db}

	t.Run("SaveUser", func(t *testing.T) {
		user := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       18,
		}
		savedUser, err := repo.SaveUser(user)
		assert.NoError(t, err)
		assert.Equal(t, user, savedUser)
	})

	t.Run("FindByID", func(t *testing.T) {
		user := &model.User{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@example.com",
			Age:       21,
		}
		savedUser, err := repo.SaveUser(user)
		assert.NoError(t, err)

		foundUser, err := repo.FindByID(&savedUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, savedUser, foundUser)
	})

	t.Run("ExistsByFirstNameAndLastName", func(t *testing.T) {
		user := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		}
		_, err = repo.SaveUser(user)

		exists := repo.ExistsByFirstNameAndLastName(user)
		assert.True(t, exists)

		user.LastName = "Smith"
		exists = repo.ExistsByFirstNameAndLastName(user)
		assert.False(t, exists)
	})
}

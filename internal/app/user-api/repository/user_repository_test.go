package repository

import (
	"testing"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/model"
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

	var repo IUserRepository = &UserRepository{Db: db}

	t.Run("SaveUser", func(t *testing.T) {
		user := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		}
		savedUser, err := repo.SaveUser(user)
		assert.NoError(t, err)
		assert.Equal(t, user.FirstName, savedUser.FirstName)
		assert.Equal(t, user.LastName, savedUser.LastName)
		assert.Equal(t, user.Email, savedUser.Email)
	})

	t.Run("FindByID", func(t *testing.T) {
		user := &model.User{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@example.com",
		}
		savedUser, err := repo.SaveUser(user)
		assert.NoError(t, err)

		foundUser, err := repo.FindByID(savedUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, savedUser.ID, foundUser.ID)
		assert.Equal(t, savedUser.FirstName, foundUser.FirstName)
		assert.Equal(t, savedUser.LastName, foundUser.LastName)
		assert.Equal(t, savedUser.Email, foundUser.Email)
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

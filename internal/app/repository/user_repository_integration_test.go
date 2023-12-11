package repository_test

import (
	"testing"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/repository"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/constant"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initAndGetRepo(t *testing.T) repository.IUserRepository {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return &repository.UserRepository{Db: db}
}

func TestUserRepository(t *testing.T) {
	t.Run("SaveUser Success", func(t *testing.T) {
		repo := initAndGetRepo(t)
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

	t.Run("SaveUser Failure due to duplicate entry", func(t *testing.T) {
		repo := initAndGetRepo(t)
		user := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       18,
		}
		savedUser, err := repo.SaveUser(user)
		assert.NoError(t, err)
		assert.Equal(t, user, savedUser)
		_, err2 := repo.SaveUser(user)
		assert.Error(t, err2)
		assert.Equal(t, constant.ErrorNameAlreadyExists, err2.Error())
	})

	t.Run("FindByID Success", func(t *testing.T) {
		repo := initAndGetRepo(t)
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

	t.Run("ExistsByFirstNameAndLastName Success", func(t *testing.T) {
		repo := initAndGetRepo(t)
		user := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		}
		repo.SaveUser(user)

		exists := repo.ExistsByFirstNameAndLastName(user)
		assert.True(t, exists)

		user.LastName = "Smith"
		exists = repo.ExistsByFirstNameAndLastName(user)
		assert.False(t, exists)
	})
}

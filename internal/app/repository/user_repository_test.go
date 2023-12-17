package repository_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/repository"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/constant"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening gorm database", err)
	}

	return gormDB, mock
}

func TestUserRepository_FindByID(t *testing.T) {
	gormDB, mock := NewMockDB()
	repo := repository.UserRepository{Db: gormDB}

	t.Run("Return no rows error when no user exists", func(t *testing.T) {
		user := &model.User{
			ID:        uuid.New(),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       18,
		}

		mock.ExpectQuery("SELECT").WithArgs(user.ID).WillReturnError(sql.ErrNoRows)

		_, err := repo.FindByID(&user.ID)
		assert.EqualError(t, err, "sql: no rows in result set")
	})

	t.Run("Return user that exists in the repo", func(t *testing.T) {
		user := &model.User{
			ID:        uuid.New(),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       18,
		}

		rows := sqlmock.NewRows([]string{"first_name", "last_name", "email", "age"}).
			AddRow(user.FirstName, user.LastName, user.Email, user.Age)

		mock.ExpectQuery("SELECT").WithArgs(user.ID).WillReturnRows(rows)

		repoUser, err := repo.FindByID(&user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.FirstName, repoUser.FirstName)
		assert.Equal(t, user.LastName, repoUser.LastName)
		assert.Equal(t, user.Email, repoUser.Email)
		assert.Equal(t, user.Age, repoUser.Age)
	})
}
func TestUserRepository_SaveUser(t *testing.T) {
	gormDB, mock := NewMockDB()
	repo := repository.UserRepository{Db: gormDB}

	t.Run("Return error when user with the same name already exists", func(t *testing.T) {
		user := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       18,
		}

		mock.ExpectQuery("SELECT").WithArgs(user.FirstName, user.LastName).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		_, err := repo.SaveUser(user)
		assert.EqualError(t, err, constant.ErrorNameAlreadyExists)
	})

	t.Run("Save user when no user with the same name exists", func(t *testing.T) {
		user := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       18,
		}

		mock.ExpectQuery("SELECT").WithArgs(user.FirstName, user.LastName).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `users`").WithArgs(sqlmock.AnyArg(), user.FirstName, user.LastName, user.Email, user.Age).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		savedUser, err := repo.SaveUser(user)
		assert.NoError(t, err)
		assert.Equal(t, user, savedUser)
	})
}
func TestUserRepository_ExistsByFirstNameAndLastName(t *testing.T) {
	gormDB, mock := NewMockDB()
	repo := repository.UserRepository{Db: gormDB}

	t.Run("Return true when user with the same first name and last name exists", func(t *testing.T) {
		user := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       18,
		}

		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery("SELECT").WithArgs(user.FirstName, user.LastName).WillReturnRows(rows)

		exists := repo.ExistsByFirstNameAndLastName(user)
		assert.True(t, exists)
	})

	t.Run("Return false when no user with the same first name and last name exists", func(t *testing.T) {
		user := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Age:       18,
		}

		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery("SELECT").WithArgs(user.FirstName, user.LastName).WillReturnRows(rows)

		exists := repo.ExistsByFirstNameAndLastName(user)
		assert.False(t, exists)
	})
}

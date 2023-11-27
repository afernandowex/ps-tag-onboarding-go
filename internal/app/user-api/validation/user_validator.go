package validation

import (
	"fmt"
	"log"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/constant"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/repository"
	"github.com/google/uuid"
)

type IUserValidationService interface {
	ValidateUser(user *model.User) []string
	ValidateUserID(ID *string) []string
}

type UserValidationService struct {
	UserRepository repository.IUserRepository
}

func (s *UserValidationService) ValidateUserID(ID *string) []string {
	var validations = []func(ID *string) string{
		s.validateID,
	}

	var errors []string
	for _, validation := range validations {
		err := validation(ID)
		if err != "" {
			errors = append(errors, err)
		}
	}

	return errors
}

func (s *UserValidationService) validateID(ID *string) string {
	if ID == nil {
		return constant.ErrorInvalidUserID
	}
	if !isUUID(*ID) {
		return constant.ErrorInvalidUserID
	}
	return ""
}

func isUUID(str string) bool {
	id, err := uuid.Parse(str)

	if id != uuid.Nil {
		return true
	} else {
		log.Println(fmt.Sprintf("Invalid UUID %s Error:=%s", str, err))
		return false
	}
}

func (s *UserValidationService) ValidateUser(user *model.User) []string {
	var validations = []func(user *model.User) string{
		s.validateName,
		s.validateEmail,
		s.validateAge,
	}

	var errors []string
	for _, validation := range validations {
		err := validation(user)
		if err != "" {
			errors = append(errors, err)
		}
	}

	return errors
}

func (s *UserValidationService) validateName(user *model.User) string {
	if user.FirstName == "" || user.LastName == "" {
		return constant.ErrorNameRequired
	}
	if s.UserRepository.ExistsByFirstNameAndLastName(user) {
		return constant.ErrorNameAlreadyExists
	}
	return ""
}

func (s *UserValidationService) validateAge(user *model.User) string {
	if user.Age < 18 {
		return constant.ErrorAgeMinimum
	}
	return ""
}

func (s *UserValidationService) validateEmail(user *model.User) string {
	if user.Email == "" {
		return constant.ErrorEmailRequired
	}
	if !containsAtSymbol(user.Email) {
		return constant.ErrorEmailFormat
	}
	return ""
}

func containsAtSymbol(email string) bool {
	for _, char := range email {
		if char == '@' {
			return true
		}
	}
	return false
}

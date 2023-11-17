package validation

import (
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/constant"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/repository"
)

type IUserValidationService interface {
	ValidateUser(user *model.User) []string
	ValidateUserID(ID *string) []string
}

type userValidationServiceImpl struct {
	userRepository repository.IUserRepository
}

func NewUserValidationService(userRepository repository.IUserRepository) IUserValidationService {
	return &userValidationServiceImpl{userRepository: userRepository}
}

func (s *userValidationServiceImpl) ValidateUserID(ID *string) []string {
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

func (s *userValidationServiceImpl) validateID(ID *string) string {
	if ID == nil {
		return constant.ErrorInvalidUserID
	}
	if !containsOnlyDigits(*ID) {
		return constant.ErrorInvalidUserID
	}
	return ""
}

func containsOnlyDigits(str string) bool {
	for _, char := range str {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func (s *userValidationServiceImpl) ValidateUser(user *model.User) []string {
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

func (s *userValidationServiceImpl) validateName(user *model.User) string {
	if user.FirstName == "" || user.LastName == "" {
		return constant.ErrorNameRequired
	}
	if s.userRepository.ExistsByFirstNameAndLastName(user) {
		return constant.ErrorNameAlreadyExists
	}
	return ""
}

func (s *userValidationServiceImpl) validateAge(user *model.User) string {
	if user.Age < 18 {
		return constant.ErrorAgeMinimum
	}
	return ""
}

func (s *userValidationServiceImpl) validateEmail(user *model.User) string {
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

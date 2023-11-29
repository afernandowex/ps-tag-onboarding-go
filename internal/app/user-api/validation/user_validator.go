package validation

import (
	"strings"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/constant"
	"github.com/google/uuid"
)

type IUserValidationService interface {
	ValidateUser(user *model.User) []string
	ValidateUserID(ID *string) []string
}

type UserValidationService struct {
}

func (s *UserValidationService) ValidateUserID(ID *string) []string {
	var errors []string
	if ID == nil || !isUUID(*ID) {
		errors = append(errors, constant.ErrorInvalidUserID)
	}

	return errors
}

func isUUID(str string) bool {
	id, _ := uuid.Parse(str)

	if id != uuid.Nil {
		return true
	} else {
		return false
	}
}

func (s *UserValidationService) ValidateUser(user *model.User) []string {
	var errors []string
	var err string = s.validateName(user)
	if err != "" {
		errors = append(errors, err)
	}
	err = s.validateAge(user)
	if err != "" {
		errors = append(errors, err)
	}
	err = s.validateEmail(user)
	if err != "" {
		errors = append(errors, err)
	}

	return errors
}

func (s *UserValidationService) validateName(user *model.User) string {
	if user.FirstName == "" || user.LastName == "" {
		return constant.ErrorNameRequired
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
	if !strings.Contains(user.Email, "@") {
		return constant.ErrorEmailFormat
	}
	return ""
}

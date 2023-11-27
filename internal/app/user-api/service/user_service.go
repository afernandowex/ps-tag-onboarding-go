package service

import (
	"log"
	"net/http"
	"strings"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/constant"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/errormessage"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/repository"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/validation"
	"github.com/google/uuid"
)

type UserValidationService interface {
	ValidateUser(user *model.User) []string
	ValidateUserID(ID *string) []string
}

type IUserService interface {
	FindUserByID(ID *string) (*model.User, *errormessage.ErrorMessage)
	SaveUser(user *model.User) (*model.User, *errormessage.ErrorMessage)
}

type UserService struct {
	Repository repository.IUserRepository
	Validator  validation.IUserValidationService
}

func (service *UserService) FindUserByID(ID *string) (*model.User, *errormessage.ErrorMessage) {
	log.Printf("Fetching user with id:= %s", *ID)
	valErrors := service.Validator.ValidateUserID(ID)
	if len(valErrors) > 0 {
		log.Println(valErrors)
		error := errormessage.NewErrorMessage(strings.Join(valErrors, ", "), http.StatusBadRequest)
		return nil, &error
	}

	id, err := uuid.Parse(*ID)

	if err != nil {
		error := errormessage.NewErrorMessage(constant.ErrorInvalidUserID, http.StatusBadRequest)
		return nil, &error
	}

	user, err := service.Repository.FindByID(id)
	if err != nil {
		error := errormessage.NewErrorMessage(constant.ErrorUserNotFound, http.StatusNotFound)
		return nil, &error
	}

	return user, nil
}

func (service *UserService) SaveUser(user *model.User) (*model.User, *errormessage.ErrorMessage) {
	log.Printf("Saving user %s %s", user.FirstName, user.LastName)
	valErrors := service.Validator.ValidateUser(user)
	if len(valErrors) > 0 {
		log.Println(valErrors)
		error := errormessage.NewErrorMessage(strings.Join(valErrors, ", "), http.StatusBadRequest)
		return nil, &error
	}
	savedUser, err := service.Repository.SaveUser(user)
	if err != nil {
		log.Printf("Error while saving user %s %s %s", user.FirstName, user.LastName, err.Error())
		error := errormessage.NewErrorMessage(err.Error(), http.StatusInternalServerError)
		return nil, &error
	}
	return savedUser, nil
}

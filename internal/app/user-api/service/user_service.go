package service

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/constant"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/errors"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/repository"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/validation"
)

type UserValidationService interface {
	ValidateUser(user *model.User) []string
	ValidateUserID(ID *string) []string
}

type IUserService interface {
	FindUserByID(ID *string) (*model.User, errors.IErrorMessage)
	SaveUser(user *model.User) (*model.User, errors.IErrorMessage)
}

type UserService struct {
	Repository repository.IUserRepository
	Validator  validation.IUserValidationService
}

func (service *UserService) FindUserByID(ID *string) (*model.User, errors.IErrorMessage) {
	log.Printf("Fetching user with id:= %s", *ID)
	valErrors := service.Validator.ValidateUserID(ID)
	if len(valErrors) > 0 {
		log.Println(valErrors)
		return nil, errors.NewErrorMessage(strings.Join(valErrors, ", "), http.StatusBadRequest)
	}
	id, err := strconv.Atoi(*ID)
	user, err := service.Repository.FindByID(int32(id))
	if err != nil {
		return nil, errors.NewErrorMessage(constant.ErrorUserNotFound, http.StatusNotFound)
	}

	return user, nil
}

func (service *UserService) SaveUser(user *model.User) (*model.User, errors.IErrorMessage) {
	log.Printf("Saving user %s %s", user.FirstName, user.LastName)
	valErrors := service.Validator.ValidateUser(user)
	if len(valErrors) > 0 {
		log.Println(valErrors)
		return nil, errors.NewErrorMessage(strings.Join(valErrors, ", "), http.StatusBadRequest)
	}
	savedUser, err := service.Repository.SaveUser(user)
	if err != nil {
		log.Printf("Error while saving user %s %s %s", user.FirstName, user.LastName, err.Error())
		return nil, errors.NewErrorMessage(err.Error(), http.StatusInternalServerError)
	}
	return savedUser, nil
}

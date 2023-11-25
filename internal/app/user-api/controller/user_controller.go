package controller

import (
	"net/http"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/constant"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/errors"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/model"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/service"
	"github.com/labstack/echo/v4"
)

type IUserController interface {
	FindUser(c echo.Context) error
	SaveUser(c echo.Context) error
}

type UserController struct {
	service service.IUserService
}

func NewUserController(service service.IUserService) IUserController {
	controller := UserController{service: service}
	return &controller
}

func (controller *UserController) FindUser(c echo.Context) error {
	id := c.Param("id")
	user, errorMessage := controller.service.FindUserByID(&id)
	if errorMessage != nil {
		return c.JSON(errorMessage.HttpStatus(), errorMessage)
	}
	return c.JSON(http.StatusOK, user)
}

func (controller *UserController) SaveUser(c echo.Context) error {
	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, errors.NewErrorMessage(constant.ErrorInvalidUserObject, http.StatusBadRequest))
	}
	savedUser, errorMessage := controller.service.SaveUser(&user)
	if errorMessage != nil {
		return c.JSON(errorMessage.HttpStatus(), errorMessage)
	}
	return c.JSON(http.StatusOK, savedUser)
}

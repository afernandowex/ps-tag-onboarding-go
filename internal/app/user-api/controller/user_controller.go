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
	Service service.IUserService
}

func (controller *UserController) FindUser(c echo.Context) error {
	id := c.Param("id")
	user, errorMessage := controller.Service.FindUserByID(&id)
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if errorMessage != nil {
		return c.JSON(errorMessage.HttpStatus(), errorMessage)
	}
	c.Response().WriteHeader(http.StatusOK)
	return c.JSON(http.StatusOK, user)
}

func (controller *UserController) SaveUser(c echo.Context) error {
	var user model.User
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, errors.NewErrorMessage(constant.ErrorInvalidUserObject, http.StatusBadRequest))
	}
	savedUser, errorMessage := controller.Service.SaveUser(&user)
	if errorMessage != nil {
		return c.JSON(errorMessage.HttpStatus(), errorMessage)
	}
	c.Response().WriteHeader(http.StatusOK)
	return c.JSON(http.StatusOK, savedUser)
}

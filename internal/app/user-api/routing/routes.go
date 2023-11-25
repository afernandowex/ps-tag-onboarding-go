package routing

import (
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/controller"
	"github.com/labstack/echo/v4"
)

type IRoutes interface {
	InitializeRoutes(e *echo.Echo)
}

type Routes struct {
	controller controller.IUserController
}

func NewUserRoutes(controller controller.IUserController) IRoutes {
	routes := Routes{controller: controller}
	return &routes
}

func (routes *Routes) InitializeRoutes(e *echo.Echo) {
	e.GET("/find/:id", routes.controller.FindUser)
	e.POST("/save", routes.controller.SaveUser)
}

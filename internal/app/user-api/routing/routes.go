package routing

import (
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/controller"
	"github.com/labstack/echo/v4"
)

type IRoutes interface {
	InitializeRoutes(e *echo.Echo)
}

type Routes struct {
	UserController controller.IUserController
}

func (routes *Routes) InitializeRoutes(e *echo.Echo) {
	e.GET("/find/:id", routes.UserController.FindUser)
	e.POST("/save", routes.UserController.SaveUser)
}

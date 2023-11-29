package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/mysql"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/repository"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/controller"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/routing"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/service"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/validation"
	"github.com/labstack/echo/v4"
)

func main() {
	db := mysql.InitialiseMySQL()
	var repo repository.IUserRepository = &repository.UserRepository{Db: db}
	var validator validation.IUserValidationService = &validation.UserValidationService{}
	var service service.IUserService = &service.UserService{Repository: repo, Validator: validator}
	var controller controller.IUserController = &controller.UserController{Service: service}

	e := echo.New()
	routes := routing.Routes{Controller: controller}
	routes.InitializeRoutes(e)

	// Fetch HttpPort from environment variable
	httpPortStr := os.Getenv("HTTP_PORT")
	httpPort, err := strconv.Atoi(httpPortStr)
	if err != nil {
		log.Fatalf("Invalid HTTP_PORT value: %s", httpPortStr)
	}

	log.Printf("Starting server on port %d", httpPort)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", httpPort)))
}

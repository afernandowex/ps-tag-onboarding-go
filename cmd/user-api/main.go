package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/controller"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/mysql"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/repository"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/routing"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/service"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/validation"
	"github.com/labstack/echo/v4"
)

func main() {
	db := mysql.InitialiseMySQL()
	repo := repository.NewUserRepository(db)
	validator := validation.NewUserValidationService(repo)
	service := service.NewUserService(repo, validator)
	controller := controller.NewUserController(service)

	e := echo.New()
	routes := routing.NewUserRoutes(controller)
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

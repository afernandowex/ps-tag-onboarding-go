package main

import (
	"fmt"
	"log"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/constant"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/controller"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/mysql"
	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/repository"
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
	e.GET("/find/:id", controller.FindUser)
	e.POST("/save", controller.SaveUser)
	log.Printf("Starting server on port %d", constant.HttpPort)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", constant.HttpPort)))
}

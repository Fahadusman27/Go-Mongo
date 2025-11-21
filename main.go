// @title API
// @version 1.0
// @description Dokumentasi API Mongo Project

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	. "Mongo/domain/config"
	"Mongo/domain/repository"
	"Mongo/domain/routes"
	"Mongo/domain/service"
	"log"

	_ "Mongo/docs"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	LoadEnv()

	client := ConnectDB()

	defer CloseDB(client)

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
	})

	api := app.Group("/api")

	api.Get("/swagger/*", fiberSwagger.WrapHandler)

	api.Use(cors.New())
	api.Use(logger.New())

	api.Static("/uploads", "./uploads")

	userRepo := repository.NewUserRepository(client)
	authService := service.NewAuthService(userRepo)

	mongoDatabase := DB.Database("alumni_management_db")
	fileRepo := repository.NewUploadsRepository(mongoDatabase)
	UploadsService := service.NewUploadsService(fileRepo, "./uploads")

	routes.SetupFileRoutes(api, UploadsService)
	routes.AuthRoutes(api, authService)
	routes.Alumni(api, &userRepo, &service.AlumniService{})
	routes.PekerjaanAlumni(api, &userRepo)
	routes.UserRoutes(api)

	port := "3000"
	log.Printf("🚀 Server running on http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}
package main

import (
	"log"
	. "tugas/domain/config"
	"tugas/domain/repository"
	"tugas/domain/routes"
	"tugas/domain/service"

	"github.com/gofiber/fiber/v2"
	"github.com/swaggo/fiber-swagger"

	_ "github.com/swaggo/fiber-swagger/example/docs"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	LoadEnv()

	client := ConnectDB()
	
	defer CloseDB(client) 

	ConnectDB()

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
	})

	api := app.Group("/api/v1")

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
	routes.Alumni(api, &userRepo)
	routes.PekerjaanAlumni(api, &userRepo)
	routes.UserRoutes(api)
	
	port := "3000"
	log.Printf("🚀 Server running on http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}
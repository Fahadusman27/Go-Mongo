package main

import (
	"log"
	"tugas/domain/config"
	. "tugas/domain/config"
	"tugas/domain/repository"
	"tugas/domain/routes"
	"tugas/domain/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	LoadEnv()

	client := ConnectDB()
	
	defer CloseDB(client) 

	config.ConnectDB()

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10MB
	})

	app.Use(cors.New())
	app.Use(logger.New())

	app.Static("/uploads", "./uploads")

	userRepo := repository.NewUserRepository(client)

	mongoDatabase := config.DB.Database("uploads_db")
	fileRepo := repository.NewUploadsRepository(mongoDatabase)
	UploadsService := service.NewUploadsService(fileRepo, "./uploads")

	app = routes.NewApp(client)
	
	routes.AuthRoutes(app, userRepo) 
	routes.Alumni(app, &userRepo)
	routes.PekerjaanAlumni(app, &userRepo)
	routes.UserRoutes(app)
	routes.SetupFileRoutes(app, UploadsService)
	
	port := "3000"
	log.Printf("🚀 Server running on http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}
package routes

import (
"tugas/domain/service"
"github.com/gofiber/fiber/v2"
)
func SetupFileRoutes(app *fiber.App, service service.UploadsService) {
	app.Post("/upload", service.UploadFile)
	app.Get("/", service.GetAllFiles)
	app.Get("/:id", service.GetFileByID)
	app.Delete("/:id", service.DeleteFile)
}
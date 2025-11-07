package routes

import (
"tugas/domain/service"
"github.com/gofiber/fiber/v2"
)

func SetupFileRoutes(api fiber.Router, service service.UploadsService) {
	api.Post("/upload", service.UploadFile)
	api.Get("/Files", service.GetAllFiles)
	api.Get("/Files/:id", service.GetFileByID)
	api.Delete("deleted/:id", service.DeleteFile)
}
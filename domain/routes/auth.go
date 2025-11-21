package routes

import (
	"Mongo/domain/service"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(api fiber.Router, authService service.AuthService) {
	api.Post("/register", authService.RegisterHandler())
	api.Post("/login", authService.LoginHandler())
}

func UserRoutes(api fiber.Router) {
	api.Get("/users", service.GetUsersService)
}


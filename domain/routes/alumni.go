package routes

import (
	. "tugas/domain/middleware"
	"tugas/domain/model"
	"tugas/domain/service"

	"github.com/gofiber/fiber/v2"
)

func Alumni(api fiber.Router, userRepo *model.UserRepository) {
    api.Get("/alumni", JWTAuth(userRepo), RequireRole("admin", "user"), service.GetAllAlumniService)
    api.Get("/alumni/:nim", JWTAuth(userRepo), RequireRole("admin", "user"), service.CheckAlumniService)
    api.Post("/alumni", JWTAuth(userRepo), RequireRole("admin"), service.CreateAlumniService)
    api.Put("/alumni/:nim", JWTAuth(userRepo), RequireRole("admin"), service.UpdateAlumniService)
    api.Delete("/alumni/:nim", JWTAuth(userRepo), RequireRole("admin"), service.DeleteAlumniService)
}

package routes

import (
	. "Mongo/domain/middleware"
	"Mongo/domain/model"
	"Mongo/domain/service"

	"github.com/gofiber/fiber/v2"
)

func Alumni(api fiber.Router, userRepo *model.UserRepository, alumniService *service.AlumniService) {
    api.Get("/alumni", JWTAuth(userRepo), RequireRole("admin", "user"), alumniService.GetAllAlumniService)
    api.Get("/alumni/:nim", JWTAuth(userRepo), RequireRole("admin", "user"), alumniService.CheckAlumniService)
    api.Post("/alumni", JWTAuth(userRepo), RequireRole("admin"), alumniService.CreateAlumniService)
    api.Put("/alumni/:nim", JWTAuth(userRepo), RequireRole("admin"), alumniService.UpdateAlumniService)
    api.Delete("/alumni/:nim", JWTAuth(userRepo), RequireRole("admin"), alumniService.DeleteAlumniService)
}


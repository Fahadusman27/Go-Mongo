package routes

import (
	"tugas/domain/middleware"
	. "tugas/domain/middleware"
	"tugas/domain/model"
	"tugas/domain/service"

	"github.com/gofiber/fiber/v2"
)

func PekerjaanAlumni(api fiber.Router, userRepo *model.UserRepository) {
	api.Get("/pekerjaan", JWTAuth(userRepo), RequireRole("admin", "user"), service.GetAllpekerjaanAlumniService)
	api.Get("/pekerjaan/:id", JWTAuth(userRepo), RequireRole("admin", "user"), service.CheckpekerjaanAlumniService)
	api.Get("/pekerjaan/alumni/:nim_alumni", JWTAuth(userRepo), RequireRole("admin"), service.CheckpekerjaanAlumniService)
	api.Post("/pekerjaan", JWTAuth(userRepo), RequireRole("admin"), service.CreatepekerjaanAlumniService)
	api.Put("/softdeleted/:id", middleware.JWTAuth(userRepo), service.SoftDeleteBynimService)
	api.Get("/trash", middleware.JWTAuth(userRepo), service.GetAllTrashService)
	api.Put("/restore/:id", middleware.JWTAuth(userRepo), service.RestoreBynimService)
	api.Delete("/deleted/:id",middleware.JWTAuth(userRepo), service.DeletePekerjaanAlumniService)
}

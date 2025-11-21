package test

import (
	"Mongo/domain/model"
	"Mongo/domain/service"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)
func setupApp() *fiber.App {
    app := fiber.New()
    return app
}

func TestCreatepekerjaanAlumniService_ValidationError(t *testing.T) {
    app := setupApp()
    app.Post("/api/pekerjaan", service.CreatepekerjaanAlumniService)

    t.Run("Body Kosong/Invalid", func(t *testing.T) {
        req := httptest.NewRequest("POST", "/api/pekerjaan", nil)
        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
    })

    t.Run("Field Wajib Kosong", func(t *testing.T) {
        payload := model.PekerjaanAlumni{
            NimAlumni: "",
            StatusKerja: "",
        }
        body, _ := json.Marshal(payload)
        req := httptest.NewRequest("POST", "/api/pekerjaan", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        
        resp, _ := app.Test(req)
        
        assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
    })
}

func TestCheckpekerjaanAlumniService_ParamValidation(t *testing.T) {
    app := setupApp()
    app.Get("/api/pekerjaan/:id", service.CheckpekerjaanAlumniService)

    req := httptest.NewRequest("GET", "/api/pekerjaan/12345", nil)
    resp, _ := app.Test(req)
    
    assert.NotEqual(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestSoftDeleteBynimService_Authorization(t *testing.T) {
    
    t.Run("Gagal - Unauthorized (Data Locals Kosong)", func(t *testing.T) {
        app := fiber.New()
        app.Put("/api/softdeleted/:id", service.SoftDeleteBynimService)
        
        req := httptest.NewRequest("PUT", "/api/softdeleted/123", nil)
        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
    })

    t.Run("Gagal - User Mencoba Hapus Punya Orang Lain", func(t *testing.T) {
        app := fiber.New()
        
        app.Use(func(c *fiber.Ctx) error {
            c.Locals("role", "user")
            c.Locals("id", 1001)
            return c.Next()
        })
        app.Put("/api/softdeleted/:id", service.SoftDeleteBynimService)

        req := httptest.NewRequest("PUT", "/api/softdeleted/2002", nil)
        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
    })

    t.Run("Sukses Logic Auth - User Hapus Punya Sendiri", func(t *testing.T) {
        app := fiber.New()
        
        app.Use(func(c *fiber.Ctx) error {
            c.Locals("role", "user")
            c.Locals("id", 1001) 
            return c.Next()
        })
        app.Put("/api/softdeleted/:id", service.SoftDeleteBynimService)

        req := httptest.NewRequest("PUT", "/api/softdeleted/1001", nil)
        resp, _ := app.Test(req)

        assert.NotEqual(t, fiber.StatusForbidden, resp.StatusCode)
        assert.NotEqual(t, fiber.StatusUnauthorized, resp.StatusCode)
    })
    
    t.Run("Sukses Logic Auth - Admin Bebas Hapus", func(t *testing.T) {
        app := fiber.New()
        
        app.Use(func(c *fiber.Ctx) error {
            c.Locals("role", "admin")
            c.Locals("id", 999) 
            return c.Next()
        })
        app.Put("/api/softdeleted/:id", service.SoftDeleteBynimService)

        req := httptest.NewRequest("PUT", "/api/softdeleted/randomID", nil)
        resp, _ := app.Test(req)

        assert.NotEqual(t, fiber.StatusForbidden, resp.StatusCode)
    })
}
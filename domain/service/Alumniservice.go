package service

import (
	"Mongo/domain/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlumniService struct {
	repo model.AlumniRepository
}

func NewAlumniService(repo model.AlumniRepository) *AlumniService {
    return &AlumniService{
        repo: repo,
    }
}

func (s *AlumniService) GetAllAlumniService(c *fiber.Ctx) error {
    // Panggil method dari interface repo
    alumniList, err := s.repo.GetAllAlumni()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal mendapatkan daftar alumni karena " + err.Error(),
            "success": false,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Berhasil mendapatkan daftar alumni",
        "success": true,
        "alumni":  alumniList,
    })
}

// HandleGetAllUsers godoc
// @Summary Dapatkan semua Alumni
// @Description Mengambil daftar semua Alumni dari database
// @Tags Alumni
// @Accept json
// @Produce json
// @Failure 400 {object} model.ErrorResponse
// @Param credentials body model.Alumni true "Data Alumni"
// @Success 200 {array} model.Alumni
// @Router /api/alumni [get]
// @Router /api/alumni/:nim [get]
// @Router /api/alumni [post]
// @Router /api/alumni/:nim [put]
// @Router /api/alumni/:nim [delete]
func (s *AlumniService) CheckAlumniService(c *fiber.Ctx) error {
    nim := c.Params("nim")
    if nim == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Gagal di ambil (NIM is missing)",
            "success": false,
        })
    }

    // Panggil method dari interface repo, bukan fungsi global
    alumni, err := s.repo.CheckAlumniByNim(nim)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return c.Status(fiber.StatusOK).JSON(fiber.Map{
                "message":  "alumni_management_db bukan alumni",
                "success":  true,
                "isAlumni": false,
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal cek alumni karena " + err.Error(),
            "success": false,
        })
    }
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message":  "Berhasil mendapatkan data alumni",
        "success":  true,
        "isAlumni": true,
        "alumni":   alumni,
    })
}

func (s *AlumniService) CreateAlumniService(c *fiber.Ctx) error {
    var alumni model.Alumni
    if err := c.BodyParser(&alumni); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid request body",
            "success": false,
        })
    }

    if alumni.NIM == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "NIM wajib diisi",
            "success": false,
        })
    }

    // Panggil method dari interface repo
    if err := s.repo.CreateAlumni(&alumni); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal membuat alumni karena " + err.Error(),
            "success": false,
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Berhasil membuat data alumni",
        "success": true,
        "alumni":  alumni,
    })
}

func (s *AlumniService) UpdateAlumniService(c *fiber.Ctx) error {
    nim := c.Params("nim")
    if nim == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "NIM wajib diisi",
            "success": false,
        })
    }

    var alumni model.Alumni
    if err := c.BodyParser(&alumni); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid request body",
            "success": false,
        })
    }

    // Panggil method dari interface repo
    if err := s.repo.UpdateAlumni(nim, &alumni); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal update alumni karena " + err.Error(),
            "success": false,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Berhasil update data alumni",
        "success": true,
        "alumni":  alumni,
    })
}

func (s *AlumniService) DeleteAlumniService(c *fiber.Ctx) error {
    nim := c.Params("nim")
    if nim == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "NIM wajib diisi",
            "success": false,
        })
    }

    // Panggil method dari interface repo
    if err := s.repo.DeleteAlumni(nim); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal menghapus alumni karena " + err.Error(),
            "success": false,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Berhasil menghapus data alumni",
        "success": true,
    })
}


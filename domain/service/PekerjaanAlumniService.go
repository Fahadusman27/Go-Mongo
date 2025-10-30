package service

import (
	"database/sql"
	"fmt"
	"strconv"
	"tugas/domain/config"
	"tugas/domain/model"
	. "tugas/domain/repository"

	"github.com/gofiber/fiber/v2"
)

func CheckpekerjaanAlumniService(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak ditemukan",
			"success": false,
		})
	}

	pekerjaan, err := CheckpekerjaanAlumniByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "Data pekerjaan alumni tidak ditemukan",
				"success": true,
				"exists":  false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal cek pekerjaan alumni karena " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Berhasil mendapatkan data pekerjaan alumni",
		"success":   true,
		"exists":    true,
		"pekerjaan": pekerjaan,
	})
}

func CreatepekerjaanAlumniService(c *fiber.Ctx) error {
	var pekerjaan model.PekerjaanAlumni
	if err := c.BodyParser(&pekerjaan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"success": false,
		})
	}

	if pekerjaan.NimAlumni == "" || pekerjaan.StatusKerja == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "IDAlumni dan StatusKerja wajib diisi",
			"success": false,
		})
	}

	if err := CreatepekerjaanAlumni(&pekerjaan); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat pekerjaan alumni karena " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "Berhasil membuat data pekerjaan alumni",
		"success":   true,
		"pekerjaan": pekerjaan,
	})
}

func UpdatepekerjaanAlumniService(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID wajib diisi",
			"success": false,
		})
	}

	var pekerjaan model.PekerjaanAlumni
	if err := c.BodyParser(&pekerjaan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"success": false,
		})
	}

	if pekerjaan.ID == 10 || pekerjaan.StatusKerja == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID dan StatusKerja wajib diisi",
			"success": false,
		})
	}

	if err := UpdatepekerjaanAlumni(id, &pekerjaan); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal update pekerjaan alumni karena " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Berhasil update data pekerjaan alumni",
		"success":   true,
		"pekerjaan": pekerjaan,
	})
}

func GetAllpekerjaanAlumniService(c *fiber.Ctx) error {
	pekerjaanList, err := GetAllpekerjaanAlumni()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mendapatkan daftar pekerjaan alumni karena " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Berhasil mendapatkan daftar pekerjaan alumni",
		"success":   true,
		"pekerjaan": pekerjaanList,
	})
}

func SoftDeleteBynimService(c *fiber.Ctx) error {
    nim := c.Params("id")
    if nim == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "NIM wajib diisi",
            "success": false,
        })
    }

    userRole, okRole := c.Locals("role").(string)
    loggedInUserID, okUser := c.Locals("id").(int)
	loggedInIDString := strconv.Itoa(loggedInUserID)

    if !okRole || !okUser || userRole == "" || loggedInUserID == 0 {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Data otentikasi tidak lengkap atau tidak valid",
            "success": false,
        })
    }

    if userRole == "admin" {
        fmt.Println("HASIL: Akses diberikan (ADMIN)")
    } else if userRole == "user" {
        if loggedInIDString != nim {
            fmt.Println("HASIL: Akses DITOLAK (NIM tidak cocok)")
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "message": "Akses ditolak: Anda hanya dapat menghapus data Anda sendiri",
                "success": false,
            })
        }
         fmt.Println("HASIL: Akses diberikan (USER, ID/NIM cocok)")
    } else {
        fmt.Println("HASIL: Akses DITOLAK (Role tidak dikenal)")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "message": "Akses ditolak: Role tidak diizinkan",
            "success": false,
        })
    }

    if err := SoftDeleteBynim(nim); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal menghapus pekerjaan alumni karena " + err.Error(),
            "success": false,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Berhasil menghapus data pekerjaan alumni",
        "success": true,
    })
}

func GetAllTrashService(c *fiber.Ctx) error {
    userID, ok := c.Locals("id").(int)
	if !ok || userID == 0 {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "user ID tidak ditemukan di token",
        })
    }
	
	role, _ := c.Locals("role").(string)

	var nimAlumni string
    if role != "admin" {
        row := config.DB.QueryRow("SELECT nim FROM alumni WHERE user_id = $1", userID)
        if err := row.Scan(&nimAlumni); err != nil {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "data alumni tidak ditemukan untuk user ini",
            })
        }
    }

    trashes, err := GetAllTrash(nimAlumni)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(trashes)
}

func RestoreBynimService(c *fiber.Ctx) error {
    nim := c.Params("id")
    if nim == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "NIM wajib diisi",
            "success": false,
        })
    }

    userRole, okRole := c.Locals("role").(string)
    loggedInUserID, okUser := c.Locals("id").(int)
	loggedInIDString := strconv.Itoa(loggedInUserID)

    if !okRole || !okUser || userRole == "" || loggedInUserID == 0 {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Data otentikasi tidak lengkap atau tidak valid",
            "success": false,
        })
    }

    if userRole == "admin" {
        fmt.Println("HASIL: Akses diberikan (ADMIN)")
    } else if userRole == "user" {
        if loggedInIDString != nim {
            fmt.Println("HASIL: Akses DITOLAK (NIM tidak cocok)")
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "message": "Akses ditolak: Anda hanya dapat mengembalikan data Anda sendiri",
                "success": false,
            })
        }
         fmt.Println("HASIL: Akses diberikan (USER, ID/NIM cocok)")
    } else {
        fmt.Println("HASIL: Akses DITOLAK (Role tidak dikenal)")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "message": "Akses ditolak: Role tidak diizinkan",
            "success": false,
        })
    }

    if err := RestoreTrashBynim(nim); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal mengembalikan pekerjaan alumni karena " + err.Error(),
            "success": false,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Berhasil mengembalikan data pekerjaan alumni",
        "success": true,
    })
}

func DeletePekerjaanAlumniService(c *fiber.Ctx) error {
    nim := c.Params("id")
    if nim == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "NIM wajib diisi",
            "success": false,
        })
    }

    userRole, okRole := c.Locals("role").(string)
    loggedInUserID, okUser := c.Locals("id").(int)
	loggedInIDString := strconv.Itoa(loggedInUserID)

    if !okRole || !okUser || userRole == "" || loggedInUserID == 0 {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Data otentikasi tidak lengkap atau tidak valid",
            "success": false,
        })
    }

    if userRole == "admin" {
        fmt.Println("HASIL: Akses diberikan (ADMIN)")
    } else if userRole == "user" {
        if loggedInIDString != nim {
            fmt.Println("HASIL: Akses DITOLAK (NIM tidak cocok)")
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "message": "Akses ditolak: Anda hanya dapat menghapus data Anda sendiri",
                "success": false,
            })
        }
         fmt.Println("HASIL: Akses diberikan (USER, ID/NIM cocok)")
    } else {
        fmt.Println("HASIL: Akses DITOLAK (Role tidak dikenal)")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "message": "Akses ditolak: Role tidak diizinkan",
            "success": false,
        })
    }

    if err := DeletePekerjaanByid(nim); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal menghapus pekerjaan alumni karena " + err.Error(),
            "success": false,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Berhasil menghapus data pekerjaan alumni",
        "success": true,
    })
}
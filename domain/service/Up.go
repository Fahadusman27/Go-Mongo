package service

import (
	"Mongo/domain/model"
	"Mongo/domain/repository"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UploadsService interface {
	UploadFile(c *fiber.Ctx) error
	GetAllFiles(c *fiber.Ctx) error
	GetFileByID(c *fiber.Ctx) error
	DeleteFile(c *fiber.Ctx) error
}
type upService struct {
	repo       repository.UploadsRepository
	uploadPath string
}

// @Param credentials body model.Uploads true "Data Files"
// @Summary Dapatkan semua Files
// @Description Mengambil daftar semua FIles dari database
// @Tags Files
// @Accept json
// @Produce json
// @Failure 400 {object} model.ErrorResponse
// @Success 200 {array} model.Uploads
// @Router /api/Files [get]
func NewUploadsService(repo repository.UploadsRepository, uploadPath string) UploadsService {
	return &upService{
		repo:       repo,
		uploadPath: uploadPath,
	}
}
func (s *upService) UploadFile(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "No file uploaded",
			"error":   err.Error(),
		})
	}

	if fileHeader.Size > 10*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File size exceeds 10MB",
		})
	}

	allowedTypes := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"image/jpg":       true,
		"application/pdf": true,
		"text/html":       true,
		"text/plain":      true,
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File type not allowed",
		})
	}

	ext := filepath.Ext(fileHeader.Filename)
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(s.uploadPath, newFileName)

	if err := os.MkdirAll(s.uploadPath, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create upload directory",
			"error":   err.Error(),
		})
	}
	// Simpan file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to open file",
			"error":   err.Error(),
		})
	}
	defer file.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to save file",
			"error":   err.Error(),
		})
	}
	defer out.Close()

	if _, err := out.ReadFrom(file); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to write file",
			"error":   err.Error(),
		})
	}
	// Simpan metadata ke database
	UploadsModel := &model.Uploads{
		UploadsName:  newFileName,
		OriginalName: fileHeader.Filename,
		UploadsPath:  filePath,
		UploadsSize:  fileHeader.Size,
		UploadsType:  contentType,
	}
	if err := s.repo.Create(UploadsModel); err != nil {
		// Hapus file jika gagal simpan ke database
		os.Remove(filePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to save file metadata",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "File uploaded successfully",
		"data":    s.toFileResponse(UploadsModel),
	})
}

func (s *upService) GetAllFiles(c *fiber.Ctx) error {
	files, err := s.repo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get files",
			"error":   err.Error(),
		})
	}
	var responses []model.UploadsResponse
	for _, file := range files {
		responses = append(responses, *s.toFileResponse(&file))
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Files retrieved successfully",
		"data":    responses,
	})
}

func (s *upService) GetFileByID(c *fiber.Ctx) error {
	id := c.Params("id")
	file, err := s.repo.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "File not found",
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "File retrieved successfully",
		"data":    s.toFileResponse(file),
	})
}

func (s *upService) DeleteFile(c *fiber.Ctx) error {
	id := c.Params("id")
	file, err := s.repo.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "File not found",
			"error":   err.Error(),
		})
	}
	// Hapus file dari storage
	if err := os.Remove(file.UploadsPath); err != nil {
		fmt.Println("Warning: Failed to delete file from storage:", err)
	}
	// Hapus dari database
	if err := s.repo.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete file",
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "File deleted successfully",
	})
}

func (s *upService) toFileResponse(file *model.Uploads) *model.UploadsResponse {
	return &model.UploadsResponse{
		ID:           file.ID.Hex(),
		UploadsName:  file.UploadsName,
		OriginalName: file.OriginalName,
		UploadsPath:  file.UploadsPath,
		UploadsSize:  file.UploadsSize,
		UploadsType:  file.UploadsType,
		UploadedAt:   file.UploadedAt,
	}
}

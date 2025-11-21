package test

import (
	"Mongo/domain/model"
	"Mongo/domain/service"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockUploadsRepository struct {
	mock.Mock
}

func (m *MockUploadsRepository) Create(upload *model.Uploads) error {
	args := m.Called(upload)
	return args.Error(0)
}

func (m *MockUploadsRepository) FindAll() ([]model.Uploads, error) {
	args := m.Called()
	return args.Get(0).([]model.Uploads), args.Error(1)
}

func (m *MockUploadsRepository) FindByID(id string) (*model.Uploads, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Uploads), args.Error(1)
}

func (m *MockUploadsRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func createTestApp() *fiber.App {
	return fiber.New()
}

func TestUploadFile_Success(t *testing.T) {
	mockRepo := new(MockUploadsRepository)
	tempDir := t.TempDir()
	service := service.NewUploadsService(mockRepo, tempDir)
	app := createTestApp()

	app.Post("/upload", service.UploadFile)

	mockRepo.On("Create", mock.AnythingOfType("*model.Uploads")).Return(nil)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.jpg")
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, "test.jpg"))
	h.Set("Content-Type", "image/jpeg")
	part.Write([]byte("dummy image content"))
	writer.WriteField("other_field", "some_value")
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockRepo.AssertExpectations(t)

	files, _ := os.ReadDir(tempDir)
	assert.Len(t, files, 1)
}

func TestUploadFile_InvalidFileType(t *testing.T) {
	mockRepo := new(MockUploadsRepository)
	tempDir := t.TempDir()
	service := service.NewUploadsService(mockRepo, tempDir)
	app := createTestApp()

	app.Post("/upload", service.UploadFile)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "malware.exe")
	part.Write([]byte("malicious content"))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUploadFile_NoFile(t *testing.T) {
	mockRepo := new(MockUploadsRepository)
	tempDir := t.TempDir()
	service := service.NewUploadsService(mockRepo, tempDir)
	app := createTestApp()

	app.Post("/upload", service.UploadFile)

	req := httptest.NewRequest("POST", "/upload", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGetAllFiles_Success(t *testing.T) {
	mockRepo := new(MockUploadsRepository)
	service := service.NewUploadsService(mockRepo, t.TempDir())
	app := createTestApp()

	app.Get("/files", service.GetAllFiles)

	dummyFiles := []model.Uploads{
		{
			ID:           primitive.NewObjectID(),
			UploadsName:  "uuid.jpg",
			OriginalName: "test.jpg",
		},
	}

	mockRepo.On("FindAll").Return(dummyFiles, nil)

	req := httptest.NewRequest("GET", "/files", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "test.jpg")

	mockRepo.AssertExpectations(t)
}

func TestGetFileByID_Success(t *testing.T) {
	mockRepo := new(MockUploadsRepository)
	service := service.NewUploadsService(mockRepo, t.TempDir())
	app := createTestApp()

	app.Get("/files/:id", service.GetFileByID)

	objID := primitive.NewObjectID()
	dummyFile := &model.Uploads{
		ID:           objID,
		OriginalName: "found.jpg",
	}

	mockRepo.On("FindByID", objID.Hex()).Return(dummyFile, nil)

	req := httptest.NewRequest("GET", "/files/"+objID.Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestGetFileByID_NotFound(t *testing.T) {
	mockRepo := new(MockUploadsRepository)
	service := service.NewUploadsService(mockRepo, t.TempDir())
	app := createTestApp()

	app.Get("/files/:id", service.GetFileByID)

	mockRepo.On("FindByID", "invalid-id").Return(nil, errors.New("not found"))

	req := httptest.NewRequest("GET", "/files/invalid-id", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestDeleteFile_Success(t *testing.T) {
	mockRepo := new(MockUploadsRepository)
	tempDir := t.TempDir()
	service := service.NewUploadsService(mockRepo, tempDir)
	app := createTestApp()

	app.Delete("/files/:id", service.DeleteFile)

	fileName := "todelete.jpg"
	filePath := filepath.Join(tempDir, fileName)
	os.WriteFile(filePath, []byte("data"), 0644)

	objID := primitive.NewObjectID()
	dummyFile := &model.Uploads{
		ID:          objID,
		UploadsPath: filePath,
	}

	mockRepo.On("FindByID", objID.Hex()).Return(dummyFile, nil)
	mockRepo.On("Delete", objID.Hex()).Return(nil)

	req := httptest.NewRequest("DELETE", "/files/"+objID.Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	_, err := os.Stat(filePath)
	assert.True(t, os.IsNotExist(err))

	mockRepo.AssertExpectations(t)
}

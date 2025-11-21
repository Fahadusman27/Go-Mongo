package test

import (
	"Mongo/domain/model"
	"Mongo/domain/service"
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockAlumniRepository struct {
	mock.Mock
}

func (m *MockAlumniRepository) GetAllAlumni() ([]model.Alumni, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Alumni), args.Error(1)
}

func (m *MockAlumniRepository) CheckAlumniByNim(nim string) (*model.Alumni, error) {
	args := m.Called(nim)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Alumni), args.Error(1)
}

func (m *MockAlumniRepository) CreateAlumni(alumni *model.Alumni) error {
	args := m.Called(alumni)
	return args.Error(0)
}

func (m *MockAlumniRepository) UpdateAlumni(nim string, alumni *model.Alumni) error {
	args := m.Called(nim, alumni)
	return args.Error(0)
}

func (m *MockAlumniRepository) DeleteAlumni(nim string) error {
	args := m.Called(nim)
	return args.Error(0)
}

func TestGetAllAlumniService(t *testing.T) {
	// Setup
	mockRepo := new(MockAlumniRepository)
	svc := service.NewAlumniService(mockRepo)
	app := fiber.New()
	app.Get("/alumni", svc.GetAllAlumniService)

	// Data Dummy
	dummyAlumni := []model.Alumni{
		{NIM: "123", Nama: "Budi"},
		{NIM: "456", Nama: "Siti"},
	}
	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetAllAlumni").Return(dummyAlumni, nil).Once()

		req := httptest.NewRequest("GET", "/alumni", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockRepo.On("GetAllAlumni").Return(nil, errors.New("db error")).Once()

		req := httptest.NewRequest("GET", "/alumni", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 500, resp.StatusCode)
	})
}

func TestCheckAlumniService(t *testing.T) {
	mockRepo := new(MockAlumniRepository)
	svc := service.NewAlumniService(mockRepo)
	app := fiber.New()
	app.Get("/alumni/:nim", svc.CheckAlumniService)

	t.Run("Found (Is Alumni)", func(t *testing.T) {
		nim := "123"
		dummyAlumni := &model.Alumni{NIM: nim, Nama: "Budi"}
		
		mockRepo.On("CheckAlumniByNim", nim).Return(dummyAlumni, nil).Once()

		req := httptest.NewRequest("GET", "/alumni/123", nil)
		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Not Found (Not Alumni / No Documents)", func(t *testing.T) {
		nim := "999"
		mockRepo.On("CheckAlumniByNim", nim).Return(nil, mongo.ErrNoDocuments).Once()

		req := httptest.NewRequest("GET", "/alumni/999", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("DB Error", func(t *testing.T) {
		nim := "error_case"
		mockRepo.On("CheckAlumniByNim", nim).Return(nil, errors.New("connection lost")).Once()

		req := httptest.NewRequest("GET", "/alumni/error_case", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 500, resp.StatusCode)
	})
}
func TestCreateAlumniService(t *testing.T) {
	mockRepo := new(MockAlumniRepository)
	svc := service.NewAlumniService(mockRepo)
	app := fiber.New()
	app.Post("/alumni", svc.CreateAlumniService)

	t.Run("Success Create", func(t *testing.T) {
		input := model.Alumni{NIM: "123", Nama: "New User"}
		body, _ := json.Marshal(input)

		mockRepo.On("CreateAlumni", mock.AnythingOfType("*model.Alumni")).Return(nil).Once()

		req := httptest.NewRequest("POST", "/alumni", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, 201, resp.StatusCode)
	})

	t.Run("Bad Request - Empty NIM", func(t *testing.T) {
		input := model.Alumni{NIM: "", Nama: "No Nim"}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/alumni", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("Bad Request - Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/alumni", bytes.NewReader([]byte(`{invalid-json`)))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestUpdateAlumniService(t *testing.T) {
	mockRepo := new(MockAlumniRepository)
	svc := service.NewAlumniService(mockRepo)
	app := fiber.New()
	app.Put("/alumni/:nim", svc.UpdateAlumniService)

	t.Run("Success Update", func(t *testing.T) {
		nim := "123"
		input := model.Alumni{NIM: "123", Nama: "Updated Name"}
		body, _ := json.Marshal(input)

		mockRepo.On("UpdateAlumni", nim, mock.AnythingOfType("*model.Alumni")).Return(nil).Once()

		req := httptest.NewRequest("PUT", "/alumni/123", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Failed Update - DB Error", func(t *testing.T) {
		nim := "123"
		input := model.Alumni{NIM: "123", Nama: "Updated Name"}
		body, _ := json.Marshal(input)

		mockRepo.On("UpdateAlumni", nim, mock.AnythingOfType("*model.Alumni")).Return(errors.New("failed update")).Once()

		req := httptest.NewRequest("PUT", "/alumni/123", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, 500, resp.StatusCode)
	})
}

func TestDeleteAlumniService(t *testing.T) {
	mockRepo := new(MockAlumniRepository)
	svc := service.NewAlumniService(mockRepo)
	app := fiber.New()
	app.Delete("/alumni/:nim", svc.DeleteAlumniService)

	t.Run("Success Delete", func(t *testing.T) {
		nim := "123"
		mockRepo.On("DeleteAlumni", nim).Return(nil).Once()

		req := httptest.NewRequest("DELETE", "/alumni/123", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Failed Delete", func(t *testing.T) {
		nim := "123"
		mockRepo.On("DeleteAlumni", nim).Return(errors.New("delete error")).Once()

		req := httptest.NewRequest("DELETE", "/alumni/123", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 500, resp.StatusCode)
	})
}
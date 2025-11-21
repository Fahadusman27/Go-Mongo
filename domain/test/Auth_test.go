package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"Mongo/domain/model"
	"Mongo/domain/service"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(email string) (*model.Users, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Users), args.Error(1)
}

func (m *MockUserRepository) Create(user *model.Users) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Count(role string) (int, error) {
	args := m.Called(role)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockUserRepository) Delete(id primitive.ObjectID) error {
    args := m.Called(id)
    return args.Error(0)
}

func (m *MockUserRepository) FindByID(id primitive.ObjectID) (*model.Users, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Users), args.Error(1)
}

func (m *MockUserRepository) FindAll() ([]model.Users, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Users), args.Error(1)
}

func (m *MockUserRepository) Update(user *model.Users) error {
	args := m.Called(user)
	return args.Error(0)
}

func setupTestApp() (*fiber.App, *MockUserRepository) {
	mockRepo := new(MockUserRepository)
	
	authService := service.NewAuthService(mockRepo)

	app := fiber.New()
	app.Post("/register", authService.RegisterHandler())
	app.Post("/login", authService.LoginHandler())

	return app, mockRepo
}

func TestRegisterHandler(t *testing.T) {
	t.Run("Success Register", func(t *testing.T) {
		app, mockRepo := setupTestApp()

		payload := model.Register{
			Email:    "new@example.com",
			Username: "newuser",
			Password: "password123",
			Role:     "user",
		}
		body, _ := json.Marshal(payload)

		mockRepo.On("FindByEmail", payload.Email).Return(nil, nil).Once()
		mockRepo.On("Create", mock.AnythingOfType("*model.Users")).Return(nil).Once()

		req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail - Email Already Exists", func(t *testing.T) {
		app, mockRepo := setupTestApp()

		payload := model.Register{
			Email:    "exist@example.com",
			Password: "password123",
			Role:     "user",
		}
		body, _ := json.Marshal(payload)

		existingUser := &model.Users{Email: payload.Email}
		mockRepo.On("FindByEmail", payload.Email).Return(existingUser, nil).Once()

		req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Fail - Invalid Role", func(t *testing.T) {
		app, _ := setupTestApp()

		payload := model.Register{
			Email:    "test@example.com",
			Password: "password123",
			Role:     "superadmin",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestLoginHandler(t *testing.T) {
	// Setup dummy password hash
	rawPassword := "secret123"
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	hashedPassword := string(hashedBytes)

	mockUserDB := &model.Users{
		ID:       primitive.NewObjectID(),
		Email:    "user@example.com",
		Password: hashedPassword,
		Role:     "user",
		Username: "testuser",
	}

	t.Run("Success Login", func(t *testing.T) {
		app, mockRepo := setupTestApp()

		payload := model.Login{
			Email:    "user@example.com",
			Password: "secret123",
		}
		body, _ := json.Marshal(payload)

		mockRepo.On("FindByEmail", payload.Email).Return(mockUserDB, nil).Once()

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Fail - Wrong Password", func(t *testing.T) {
		app, mockRepo := setupTestApp()

		payload := model.Login{
			Email:    "user@example.com",
			Password: "wrongpassword",
		}
		body, _ := json.Marshal(payload)

		mockRepo.On("FindByEmail", payload.Email).Return(mockUserDB, nil).Once()

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Fail - Email Not Found", func(t *testing.T) {
		app, mockRepo := setupTestApp()

		payload := model.Login{
			Email:    "ghost@example.com",
			Password: "any",
		}
		body, _ := json.Marshal(payload)

		mockRepo.On("FindByEmail", payload.Email).Return(nil, nil).Once()

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Fail - Database Error", func(t *testing.T) {
		app, mockRepo := setupTestApp()

		payload := model.Login{
			Email:    "error@example.com",
			Password: "any",
		}
		body, _ := json.Marshal(payload)

		mockRepo.On("FindByEmail", payload.Email).Return(nil, errors.New("db connection failed")).Once()

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}
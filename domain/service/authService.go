package service

import (
	"time"
	"tugas/domain/config"
	"tugas/domain/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterHandler() fiber.Handler
	LoginHandler() fiber.Handler
	MeHandler() fiber.Handler
}

type authService struct {
	userRepo model.UserRepository
	// bisa tambah jwtSecret, expiry, dll
}

func NewAuthService(userRepo model.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// ---------------- HANDLER REGISTER ----------------
func (s *authService) RegisterHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return s.registerLogic(c)
	}
}

func (s *authService) registerLogic(c *fiber.Ctx) error {
	// 1. Parse body JSON
	var body struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	// 2. Validasi role
	if body.Role != "admin" && body.Role != "user" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "role tidak valid, harus 'admin' atau 'user'"})
	}

	// 3. Cek email sudah ada
	existing, err := s.userRepo.FindByEmail(body.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "gagal mengecek email"})
	}
	if existing != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email sudah terdaftar"})
	}

	// 4. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "gagal mengenkripsi password"})
	}

	// 5. Buat user
	user := &model.Users{
		Email:    body.Email,
		Username: body.Username,
		Password: string(hashedPassword),
		Role:     body.Role,
	}

	// 6. Simpan ke DB
	if err := s.userRepo.Create(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "gagal membuat user"})
	}

	// 7. Buat JWT
	token, err := s.generateJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "gagal membuat token"})
	}

	// 8. Hapus password sebelum return
	user.Password = ""

	// 9. Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user berhasil didaftarkan",
		"token":   token,
		"user":    user,
	})
}

// ---------------- HANDLER LOGIN ----------------
func (s *authService) LoginHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return s.loginLogic(c)
	}
}

func (s *authService) loginLogic(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	user, err := s.userRepo.FindByEmail(body.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "gagal mencari user"})
	}
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "email atau password salah"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "email atau password salah"})
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "gagal membuat token"})
	}

	user.Password = ""

	return c.JSON(fiber.Map{
		"message": "login berhasil",
		"token":   token,
		"user":    user,
	})
}

func (s *authService) generateJWT(user *model.Users) (string, error) {
	secret := config.GetJWTSecret()
	expiry := config.GetJWTExpiry()
	userIDHex := user.ID.Hex()

	claims := jwt.MapClaims{
		"sub":      userIDHex,
		"role":     user.Role,
		"username": user.Username,
		"exp":      time.Now().Add(expiry).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *authService) MeHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return s.meLogic(c)
	}
}

func (s *authService) meLogic(c *fiber.Ctx) error {
	return nil
}

package service

import (
	"errors"
	"time"
	"tugas/domain/config"
	"tugas/domain/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo model.UserRepository
}

func NewUserService(userRepo model.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

type AuthService interface {
	Login(email, password string) (string, *model.Users, error)
	ParseToken(tokenString string) (*jwt.RegisteredClaims, error)
}

type authService struct {
	userRepo model.UserRepository
	secret   string
	expiry   time.Duration
}

func NewAuthService(userRepo model.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
		secret:   config.GetJWTSecret(),
		expiry:   config.GetJWTExpiry(),
	}
}

func (s *authService) Login(email, password string) (string, *model.Users, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("gagal mencari email") 
	}

	if user == nil {
		return "", nil, errors.New("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("email atau password salah")
	}
	
	userIDHex := user.ID.Hex() 

	claims := jwt.RegisteredClaims{
		Subject:   userIDHex, 
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      userIDHex, 
		"role":     user.Role,
		"username": user.Username,
		"exp":      claims.ExpiresAt.Unix(),
		"iat":      claims.IssuedAt.Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", nil, errors.New("gagal membuat token JWT")
	}

	user.Password = ""
	return tokenString, user, nil
}

func (s *authService) ParseToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token claims")
}
package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email     string             `bson:"email" json:"email"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"password"`
	Role      string             `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type Login struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

type Register struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type ErrorResponse struct {
	Code  int    `json:"code" example:"500"`
	Error string `json:"error" example:"Failed to fetch users"`
}

type UserRepository interface {
	FindByID(id primitive.ObjectID) (*Users, error)
	FindByEmail(email string) (*Users, error)
	FindAll() ([]Users, error)
	Create(user *Users) error
	Update(user *Users) error
	Delete(id primitive.ObjectID) error
	Count(search string) (int, error)
}


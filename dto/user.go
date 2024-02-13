package dto

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CreateUserRequest struct {
	UserStrID string `json:"user_str_id" validate:"required,alphanum"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type JwtCustomClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

type CreateUserResponse struct {
	ID        uint      `json:"id"`
	UserStrID string    `json:"user_str_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

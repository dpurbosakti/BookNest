package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone" binding:"phone"`
	Address  string `json:"address" binding:"required"`
}

type UserResponse struct {
	Id               uuid.UUID      `json:"id"`
	Name             string         `json:"name"`
	Email            string         `json:"email"`
	Phone            string         `json:"phone"`
	Address          string         `json:"address"`
	Role             string         `json:"role"`
	VerificationCode string         `json:"verificationCode"`
	IsVerified       bool           `json:"is_verified"`
	CreatedAt        time.Time      `json:"created_At"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at"`
}

type UserVerifyRequest struct {
	Email            string `json:"email" binding:"required,email"`
	VerificationCode string `json:"verification_code" binding:"required"`
}

type UserVerificationCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

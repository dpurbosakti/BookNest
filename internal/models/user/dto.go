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
	IsVerified       bool           `json:"isVerified"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `json:"deletedAt"`
}

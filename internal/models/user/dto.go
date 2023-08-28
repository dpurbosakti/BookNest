package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserCreateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type UserResponse struct {
	Id               uuid.UUID      `json:"id"`
	Name             string         `json:"name"`
	Email            string         `json:"email"`
	Password         string         `json:"password"`
	Phone            string         `json:"phone"`
	Address          string         `json:"address"`
	Role             string         `json:"role"`
	VerificationCode string         `json:"verificationCode"` //verification code
	IsVerified       bool           `json:"isVerified"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `json:"deletedAt"`
}

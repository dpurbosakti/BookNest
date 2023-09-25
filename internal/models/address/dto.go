package address

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AddressCreateRequest struct {
	UserId     uuid.UUID `json:"user_id"`
	Address    string    `json:"address" binding:"required"`
	Notes      string    `json:"notes" binding:"required"`
	Latitude   float64   `json:"latitude" binding:"required"`
	Longitude  float64   `json:"longitude" binding:"required"`
	PostalCode string    `json:"postal_code" binding:"required"`
}

type AddressResponse struct {
	Id         uint64         `json:"id"`
	UserId     uuid.UUID      `json:"user_id"`
	Address    string         `json:"address"`
	Notes      string         `json:"notes"`
	Latitude   float64        `json:"latitude"`
	Longitude  float64        `json:"longitude"`
	PostalCode string         `json:"postal_code"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}

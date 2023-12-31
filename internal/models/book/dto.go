package book

import (
	"time"

	"gorm.io/gorm"
)

type BookCreateRequest struct {
	Title         string  `json:"title" binding:"required"`
	Author        string  `json:"author" binding:"required"`
	RentFeePerDay float64 `json:"rent_fee_per_day" binding:"required"`
	Length        uint    `json:"length" binding:"required"`
	Width         uint    `json:"width" binding:"required"`
	Height        uint    `json:"height" binding:"required"`
}

type BookResponse struct {
	Id            uint           `json:"id"`
	Title         string         `json:"title"`
	Author        string         `json:"author"`
	RentFeePerDay float64        `json:"rent_fee_per_day"`
	Length        uint           `json:"length"`
	Width         uint           `json:"width"`
	Height        uint           `json:"height"`
	IsAvailable   bool           `json:"is_available"`
	AvailableAt   *time.Time     `json:"available_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at"`
}

type BookUpdateRequest struct {
	Title         *string  `json:"title"`
	Author        *string  `json:"author"`
	RentFeePerDay *float64 `json:"rent_fee_per_day"`
	Length        *uint    `json:"length"`
	Width         *uint    `json:"width"`
	Height        *uint    `json:"height"`
}

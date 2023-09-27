package book

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	Id            uint           `json:"id" gorm:"primaryKey,autoIncrement"`
	Title         string         `json:"title" gorm:"type:varchar(250)"`
	Author        string         `json:"author" gorm:"type:varchar(250)"`
	RentFeePerDay float64        `json:"rent_fee_per_day"`
	Length        uint           `json:"length"`
	Width         uint           `json:"width"`
	Height        uint           `json:"height"`
	IsAvailable   bool           `json:"is_available" gorm:"default:true"`
	AvailableAt   *time.Time     `json:"available_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (b *Book) Copier(input *BookUpdateRequest) {
	if input.Title != nil {
		b.Title = *input.Title
	}

	if input.Author != nil {
		b.Author = *input.Author
	}

	if input.RentFeePerDay != nil {
		b.RentFeePerDay = *input.RentFeePerDay
	}

	if input.Length != nil {
		b.Length = *input.Length
	}

	if input.Width != nil {
		b.Width = *input.Width
	}

	if input.Height != nil {
		b.Height = *input.Height
	}
}

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
	IsAvailable   bool           `json:"is_available"`
	AvailableAt   *time.Time     `json:"available_at" gorm:"default:true"`
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
}

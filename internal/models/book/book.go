package book

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	Id          uint    `json:"id" gorm:"primaryKey,autoIncrement"`
	Title       string  `json:"title" gorm:"type:varchar(250)"`
	Author      string  `json:"author" gorm:"type:varchar(250)"`
	RentFee     float64 `json:"rent_fee"`
	IsAvailable bool    `json:"is_available"`
	AvailableAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

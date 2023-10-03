package order

import (
	mb "book-nest/internal/models/book"
	mu "book-nest/internal/models/user"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	Id            uint           `json:"id" gorm:"primaryKey,autoIncrement"`
	ReferenceId   string         `json:"reference_id"`
	UserId        uuid.UUID      `json:"user_id"`
	User          *mu.User       `json:"user,omitempty"`
	BookId        uint           `json:"book_id"`
	Book          *mb.Book       `json:"book,omitempty"`
	BorrowingDate time.Time      `json:"borrowing_date"`
	ReturnedDate  time.Time      `json:"returned_date"`
	Fee           float64        `json:"fee"`
	PaymentMethod string         `json:"payment_method"`
	PaymentStatus string         `json:"payment_status" gorm:"default:initiate"`
	Status        string         `json:"status" gorm:"default:initiate"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (o *Order) GetDaysBetween() int {
	duration := o.ReturnedDate.Sub(o.BorrowingDate)
	return int(duration.Hours() / 24)
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ReturnedDate.Before(o.BorrowingDate) {
		return errors.New("returned_date cannot be earlier than borrowing_date")
	}
	return nil
}

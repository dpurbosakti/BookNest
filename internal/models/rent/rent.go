package rent

import (
	mb "book-nest/internal/models/book"
	mu "book-nest/internal/models/user"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Rent struct {
	Id            uint           `json:"id" gorm:"primaryKey,autoIncrement"`
	ReferenceId   string         `json:"reference_id"`
	UserId        uuid.UUID      `json:"user_id"`
	User          *mu.User       `json:"user,omitempty"`
	BookId        uint           `json:"book_id"`
	Book          *mb.Book       `json:"book,omitempty"`
	BorrowingDate time.Time      `json:"borrowing_date"`
	ReturnedDate  time.Time      `json:"returned_date"`
	Fee           float64        `json:"fee"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (r *Rent) GetDaysBetween() int {
	duration := r.ReturnedDate.Sub(r.BorrowingDate)
	return int(duration.Hours() / 24)
}

func (r *Rent) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ReturnedDate.Before(r.BorrowingDate) {
		return errors.New("returned_date cannot be earlier than borrowing_date")
	}
	return nil
}

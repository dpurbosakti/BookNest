package order

import (
	mb "book-nest/internal/models/book"
	mu "book-nest/internal/models/user"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderCreateRequest struct {
	BookId        uint      `json:"book_id" binding:"required"`
	BorrowingDate time.Time `json:"borrowing_date" binding:"required"`
	ReturnedDate  time.Time `json:"returned_date" binding:"required"`
}

type OrderResponse struct {
	Id            uint           `json:"id"`
	ReferenceId   string         `json:"reference_id"`
	UserId        uuid.UUID      `json:"user_id"`
	User          *mu.User       `json:"user,omitempty"`
	BookId        uint           `json:"book_id"`
	Book          *mb.Book       `json:"book,omitempty"`
	BorrowingDate time.Time      `json:"borrowing_date"`
	ReturnedDate  time.Time      `json:"returned_date"`
	Fee           float64        `json:"fee"`
	PaymentMethod string         `json:"payment_method"`
	PaymentStatus string         `json:"payment_status"`
	Status        string         `json:"status"`
	Token         *string        `json:"token,omitempty"`
	RedirectURL   *string        `json:"redirect_url,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at"`
}

type OrderUpdateRequest struct {
	ReferenceId     string
	PaymentMethod   string
	PaymentStatus   string
	TransactionTime string
	PaymentType     string
	GrossAmount     string
}

func (or *OrderResponse) GetDaysBetween() int {
	duration := or.ReturnedDate.Sub(or.BorrowingDate)
	return int(duration.Hours() / 24)
}

func (our *OrderUpdateRequest) Copier(paymentStatus, referenceId, transactionTime, paymentType string) {
	if paymentStatus != "" {
		our.PaymentStatus = paymentStatus
	}

	if referenceId != "" {
		our.ReferenceId = referenceId
	}

	if transactionTime != "" {
		our.TransactionTime = transactionTime
	}

	if paymentType != "" {
		our.PaymentType = paymentType
	}
}

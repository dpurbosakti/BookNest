package rent

import (
	mb "book-nest/internal/models/book"
	mu "book-nest/internal/models/user"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RentCreateRequest struct {
	BookId        uint      `json:"book_id" binding:"required"`
	BorrowingDate time.Time `json:"borrowing_date" binding:"required"`
	ReturnedDate  time.Time `json:"returned_date" binding:"required"`
}

type RentResponse struct {
	Id            uint           `json:"id"`
	ReferenceId   string         `json:"reference_id"`
	UserId        uuid.UUID      `json:"user_id"`
	User          *mu.User       `json:"user,omitempty"`
	BookId        uint           `json:"book_id"`
	Book          *mb.Book       `json:"book,omitempty"`
	BorrowingDate time.Time      `json:"borrowing_date"`
	ReturnedDate  time.Time      `json:"returned_date"`
	Fee           float64        `json:"fee"`
	PaymentStatus string         `json:"payment_status"`
	Token         *string        `json:"token,omitempty"`
	RedirectURL   *string        `json:"redirect_url,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at"`
}

type RentUpdateRequest struct {
	ReferenceId     string
	PaymentStatus   string
	TransactionTime string
	PaymentType     string
}

func (rr *RentResponse) GetDaysBetween() int {
	duration := rr.ReturnedDate.Sub(rr.BorrowingDate)
	return int(duration.Hours() / 24)
}

func (rur *RentUpdateRequest) Copier(paymentStatus, referenceId, transactionTime, paymentType string) {
	if paymentStatus != "" {
		rur.PaymentStatus = paymentStatus
	}

	if referenceId != "" {
		rur.ReferenceId = referenceId
	}

	if transactionTime != "" {
		rur.TransactionTime = transactionTime
	}

	if paymentType != "" {
		rur.PaymentType = paymentType
	}
}

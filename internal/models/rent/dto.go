package rent

import (
	mb "book-nest/internal/models/book"
	mu "book-nest/internal/models/user"
	"time"

	"github.com/google/uuid"
)

type RentCreateRequest struct {
	BookId        uint      `json:"book_id"`
	BorrowingDate time.Time `json:"borrowing_date"`
	ReturnedDate  time.Time `json:"returned_date"`
}

type RentResponse struct {
	Id            uint      `json:"id"`
	UserId        uuid.UUID `json:"user_id"`
	User          *mu.User  `json:"user,omitempty"`
	BookId        uint      `json:"book_id"`
	Book          *mb.Book  `json:"book,omitempty"`
	BorrowingDate time.Time `json:"borrowing_date"`
	ReturnedDate  time.Time `json:"returned_date"`
	Fee           float64   `json:"fee"`
}

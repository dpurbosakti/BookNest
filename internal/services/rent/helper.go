package rent

import (
	mr "book-nest/internal/models/rent"
	"math/rand"
	"time"
)

var (
	src = rand.NewSource(time.Now().UnixNano())
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// mappers
func requestToModel(input *mr.RentCreateRequest) *mr.Rent {
	return &mr.Rent{
		BookId:        input.BookId,
		BorrowingDate: input.BorrowingDate,
		ReturnedDate:  input.ReturnedDate,
	}
}

func modelToResponse(input *mr.Rent) *mr.RentResponse {
	return &mr.RentResponse{
		Id:            input.Id,
		ReferenceId:   input.ReferenceId,
		UserId:        input.UserId,
		User:          input.User,
		BookId:        input.BookId,
		Book:          input.Book,
		BorrowingDate: input.BorrowingDate,
		ReturnedDate:  input.ReturnedDate,
		Fee:           input.Fee,
		CreatedAt:     input.CreatedAt,
		UpdatedAt:     input.UpdatedAt,
		DeletedAt:     input.DeletedAt,
	}
}

func calculateRentFee(dayRent int, feePerDay float64) float64 {
	return float64(dayRent) * feePerDay
}

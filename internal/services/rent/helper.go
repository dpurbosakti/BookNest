package rent

import (
	mr "book-nest/internal/models/rent"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const refLength = 6
const PaymentSettlement = "settlement"

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
		PaymentStatus: input.PaymentStatus,
		Status:        input.Status,
		CreatedAt:     input.CreatedAt,
		UpdatedAt:     input.UpdatedAt,
		DeletedAt:     input.DeletedAt,
	}
}

func calculateRentFee(dayRent int, feePerDay float64) float64 {
	return float64(dayRent) * feePerDay
}

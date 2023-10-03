package order

import (
	mo "book-nest/internal/models/order"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const refLength = 6

// mappers
func requestToModel(input *mo.OrderCreateRequest) *mo.Order {
	return &mo.Order{
		BookId:        input.BookId,
		BorrowingDate: input.BorrowingDate,
		ReturnedDate:  input.ReturnedDate,
	}
}

func modelToResponse(input *mo.Order) *mo.OrderResponse {
	return &mo.OrderResponse{
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

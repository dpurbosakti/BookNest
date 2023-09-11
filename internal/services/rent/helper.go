package rent

import (
	mr "book-nest/internal/models/rent"
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

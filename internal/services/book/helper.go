package book

import (
	mb "book-nest/internal/models/book"
)

// mappers
func requestToModel(input *mb.BookCreateRequest) *mb.Book {
	return &mb.Book{
		Title:         input.Title,
		Author:        input.Author,
		RentFeePerDay: input.RentFeePerDay,
	}
}

func modelToResponse(input *mb.Book) *mb.BookResponse {
	return &mb.BookResponse{
		Id:            input.Id,
		Title:         input.Title,
		Author:        input.Author,
		RentFeePerDay: input.RentFeePerDay,
		IsAvailable:   input.IsAvailable,
		AvailableAt:   input.AvailableAt,
		CreatedAt:     input.CreatedAt,
		UpdatedAt:     input.UpdatedAt,
		DeletedAt:     input.DeletedAt,
	}
}

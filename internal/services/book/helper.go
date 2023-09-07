package book

import (
	mb "book-nest/internal/models/book"
)

// mappers
func requestToModel(input *mb.BookCreateRequest) *mb.Book {
	return &mb.Book{
		Title:   input.Title,
		Author:  input.Author,
		RentFee: input.RentFee,
	}
}

func modelToResponse(input *mb.Book) *mb.BookResponse {
	return &mb.BookResponse{
		Id:          input.Id,
		Title:       input.Title,
		Author:      input.Author,
		RentFee:     input.RentFee,
		IsAvailable: input.IsAvailable,
		AvailableAt: input.AvailableAt,
		CreatedAt:   input.CreatedAt,
		UpdatedAt:   input.UpdatedAt,
		DeletedAt:   input.DeletedAt,
	}
}

func copier(originalData *mb.Book, newData *mb.BookUpdateRequest) {
	if newData.Title != nil {
		originalData.Title = *newData.Title
	}

	if newData.Author != nil {
		originalData.Author = *newData.Author
	}

	if newData.RentFee != nil {
		originalData.RentFee = *newData.RentFee
	}
}

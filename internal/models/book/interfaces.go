package book

import (
	"book-nest/utils/pagination"

	"gorm.io/gorm"
)

type BookRepository interface {
	Create(tx *gorm.DB, input *Book) (*Book, error)
	CheckDuplicate(tx *gorm.DB, input *Book) error
	Update(tx *gorm.DB, input *Book) (*Book, error)
	GetDetail(tx *gorm.DB, bookId uint) (*Book, error)
	GetList(tx *gorm.DB, page pagination.Pagination) (pagination.Pagination, error)
	Delete(tx *gorm.DB, bookId uint) error
}

type BookService interface {
	Create(input *BookCreateRequest) (*BookResponse, error)
	Update(input *BookUpdateRequest, bookId uint) (*BookResponse, error)
	Delete(bookId uint) error
	GetDetail(userId uint) (*BookResponse, error)
	GetList(page pagination.Pagination) (pagination.Pagination, error)
}

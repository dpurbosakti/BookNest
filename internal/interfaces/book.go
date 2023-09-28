package interfaces

import (
	mb "book-nest/internal/models/book"
	"book-nest/utils/pagination"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookHandler interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetDetail(c *gin.Context)
	GetList(c *gin.Context)
}

type BookService interface {
	Create(input *mb.BookCreateRequest) (*mb.BookResponse, error)
	Update(input *mb.BookUpdateRequest, bookId uint) (*mb.BookResponse, error)
	Delete(bookId uint) error
	GetDetail(userId uint) (*mb.BookResponse, error)
	GetList(page pagination.Pagination) (pagination.Pagination, error)
}

type BookRepository interface {
	Create(tx *gorm.DB, input *mb.Book) (*mb.Book, error)
	CheckDuplicate(tx *gorm.DB, input *mb.Book) error
	Update(tx *gorm.DB, input *mb.Book) (*mb.Book, error)
	GetDetail(tx *gorm.DB, bookId uint) (*mb.Book, error)
	GetList(tx *gorm.DB, page pagination.Pagination) (pagination.Pagination, error)
	Delete(tx *gorm.DB, bookId uint) error
}

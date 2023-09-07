package book

import (
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
	Create(input *BookCreateRequest) (*BookResponse, error)
	Update(input *BookUpdateRequest, bookId uint) (*BookResponse, error)
	Delete(bookId uint) error
	GetDetail(userId uint) (*BookResponse, error)
	GetList(page pagination.Pagination) (pagination.Pagination, error)
}

type BookRepository interface {
	Create(tx *gorm.DB, input *Book) (*Book, error)
	CheckDuplicate(tx *gorm.DB, input *Book) error
	Update(tx *gorm.DB, input *Book) (*Book, error)
	GetDetail(tx *gorm.DB, bookId uint) (*Book, error)
	GetList(tx *gorm.DB, page pagination.Pagination) (pagination.Pagination, error)
	Delete(tx *gorm.DB, bookId uint) error
}

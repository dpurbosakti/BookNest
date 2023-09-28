package book

import (
	i "book-nest/internal/interfaces"
	mb "book-nest/internal/models/book"
	"book-nest/utils/pagination"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type BookRepository struct {
}

func NewBookRepository() i.BookRepository {
	return &BookRepository{}
}

func (repo *BookRepository) Create(tx *gorm.DB, input *mb.Book) (*mb.Book, error) {
	resultcreate := tx.Create(&input)
	if resultcreate.Error != nil {
		return nil, resultcreate.Error
	}

	return input, nil
}

func (repo *BookRepository) CheckDuplicate(tx *gorm.DB, input *mb.Book) error {
	var count int64
	if resultTitle := tx.Model(&mb.Book{}).Where("title = $1 ", input.Title).Count(&count); resultTitle.Error != nil {
		return errors.New("error checking title")
	}
	if count > 0 {
		return errors.New("title already exists in database")
	}

	return nil
}

func (repo *BookRepository) Update(tx *gorm.DB, input *mb.Book) (*mb.Book, error) {
	result := tx.Save(&input)
	if result.Error != nil {
		return nil, errors.New("error updating your data")
	}

	return input, nil
}

func (repo *BookRepository) GetDetail(tx *gorm.DB, bookId uint) (*mb.Book, error) {
	book := new(mb.Book)
	result := tx.First(&book, bookId)
	if result.Error != nil {
		return nil, fmt.Errorf("user id %d not found", bookId)
	}

	return book, nil
}

func (repo *BookRepository) GetList(tx *gorm.DB, page pagination.Pagination) (pagination.Pagination, error) {
	var books []mb.Book

	tx.Scopes(pagination.Paginate(books, &page, tx)).Find(&books)
	page.Rows = books

	return page, nil
}

func (repo *BookRepository) Delete(tx *gorm.DB, bookId uint) error {
	var book mb.Book
	result := tx.Delete(&book, bookId)
	if result.Error != nil {
		return fmt.Errorf("user id %d not found", bookId)
	}

	return nil
}

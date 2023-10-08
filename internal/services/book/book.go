package book

import (
	"book-nest/internal/constant"
	i "book-nest/internal/interfaces"
	mb "book-nest/internal/models/book"
	eh "book-nest/utils/errorhelper"
	"book-nest/utils/pagination"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BookService struct {
	BookRepository  i.BookRepository
	OrderRepository i.OrderRepository
	DB              *gorm.DB
}

func NewBookService(bookRepository i.BookRepository, orderRepository i.OrderRepository, db *gorm.DB) i.BookService {
	return &BookService{
		BookRepository:  bookRepository,
		OrderRepository: orderRepository,
		DB:              db,
	}
}

func (srv *BookService) Create(input *mb.BookCreateRequest) (*mb.BookResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "create",
		"scope": "book service",
	})

	if input == nil {
		return nil, errors.New("input is nil")
	}

	logger.WithField("data", input)
	result := new(mb.BookResponse)
	data := requestToModel(input)
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.WithField("data", data).Info("db transaction begin")
		err := srv.BookRepository.CheckDuplicate(tx, data)
		if err != nil {
			logger.WithError(err).Error("failed to create book, there is duplicated data")
			return err
		}
		resultRepo, err := srv.BookRepository.Create(tx, data)
		if err != nil {
			logger.WithError(err).Error("failed to create book")
			return err
		}
		result = modelToResponse(resultRepo)
		logger.WithField("data", data).Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.WithError(err).Error("failed to create book")
		return nil, err
	}

	return result, nil
}

func (srv *BookService) GetDetail(bookId uint) (*mb.BookResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_detail",
		"scope": "book service",
	})
	logger.WithField("id", bookId)
	result := new(mb.BookResponse)
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultRepo, err := srv.BookRepository.GetDetail(tx, bookId)
		if err != nil {
			eh.FailedGetDetail(logger, err, "book")
			return err
		}
		result = modelToResponse(resultRepo)
		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		eh.FailedGetDetail(logger, err, "book")
		return result, err
	}

	return result, nil
}

func (srv *BookService) GetList(page pagination.Pagination) (pagination.Pagination, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_list",
		"scope": "book service",
	})
	var result pagination.Pagination
	logger.Info("data page", page)
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultRepo, err := srv.BookRepository.GetList(tx, page)
		if err != nil {
			logger.WithError(err).Error("failed to get list")
			return err
		}
		result = resultRepo
		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.Error("failed to get list")
		return result, err
	}

	return result, nil
}

func (srv *BookService) Delete(bookId uint) error {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "delete",
		"scope": "book service",
	})
	logger.WithField("id", bookId)
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		err := srv.BookRepository.Delete(tx, bookId)
		if err != nil {
			logger.WithError(err).Error("failed to delete data")
			return err
		}
		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.Error("failed to delete data")
		return err
	}

	return nil
}

func (srv *BookService) Update(input *mb.BookUpdateRequest, bookId uint) (*mb.BookResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "update",
		"scope": "book service",
	})
	logger.WithField("data", input)

	result := new(mb.BookResponse)

	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultGet, err := srv.BookRepository.GetDetail(tx, bookId)
		if err != nil {
			eh.FailedGetDetail(logger, err, "book")
			return err
		}
		resultGet.Copier(input)
		resultUpdate, err := srv.BookRepository.Update(tx, resultGet)
		if err != nil {
			eh.FailedUpdate(logger, err, "book")
			return err
		}
		result = modelToResponse(resultUpdate)
		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		eh.FailedUpdate(logger, err, "book")
		return result, err
	}

	return result, nil
}

func (srv *BookService) Return(bookId uint) error {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "update",
		"scope": "book service",
	})
	logger.WithField("book_id", bookId).Info()
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		book, err := srv.BookRepository.GetDetail(tx, bookId)
		if err != nil {
			eh.FailedGetDetail(logger, err, "book")
			return err
		}

		book.AvailableAt = nil
		book.IsAvailable = true

		_, err = srv.BookRepository.Update(tx, book)
		if err != nil {
			eh.FailedUpdate(logger, err, "book")
			return err
		}

		order, err := srv.OrderRepository.GetByBook(tx, bookId)
		if err != nil {
			eh.FailedGetDetail(logger, err, "order")
			return err
		}

		order.Status = constant.StatusCompleted

		_, err = srv.OrderRepository.Update(tx, order)
		if err != nil {
			eh.FailedUpdate(logger, err, "order")
			return err
		}

		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		eh.FailedUpdate(logger, err, "book")
		return err
	}

	return nil
}

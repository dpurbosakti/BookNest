package rent

import (
	mb "book-nest/internal/models/book"
	mr "book-nest/internal/models/rent"
	mu "book-nest/internal/models/user"
	eh "book-nest/utils/emailhelper"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RentService struct {
	RentRepository mr.RentRepository
	BookRepository mb.BookRepository
	UserRepository mu.UserRepository
	DB             *gorm.DB
	EmailHelper    eh.Emailer
}

func NewRentService(rentRepository mr.RentRepository, bookRepository mb.BookRepository, userRepository mu.UserRepository, db *gorm.DB, emailer eh.Emailer) mr.RentService {
	return &RentService{
		RentRepository: rentRepository,
		BookRepository: bookRepository,
		UserRepository: userRepository,
		DB:             db,
		EmailHelper:    emailer,
	}
}

func (srv *RentService) Create(input *mr.RentCreateRequest, userId uuid.UUID) (*mr.RentResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "create",
		"scope": "rent service",
	})

	if input == nil {
		return nil, errors.New("input is nil")
	}

	logger.WithField("data", input)
	result := new(mr.RentResponse)
	data := requestToModel(input)
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.WithField("data", data).Info("db transaction begin")
		resultBook, err := srv.BookRepository.GetDetail(tx, data.BookId)
		if err != nil {
			logger.WithError(err).Error("failed to get book")
			return err
		}
		data.Book = resultBook

		resultUser, err := srv.UserRepository.GetDetail(tx, data.UserId)
		if err != nil {
			logger.WithError(err).Error("failed to get user")
			return err
		}
		data.User = resultUser

		fee := calculateRentFee(data.GetDaysBetween(), data.Book.RentFeePerDay)
		data.Fee = fee

		resultRepo, err := srv.RentRepository.Create(tx, data)
		if err != nil {
			logger.WithError(err).Error("failed to create rent")
			return err
		}
		resultBook.AvailableAt = &data.ReturnedDate
		resultBook.IsAvailable = false
		_, err = srv.BookRepository.Update(tx, resultBook)
		if err != nil {
			logger.WithError(err).Error("failed to update book")
			return err
		}
		result = modelToResponse(resultRepo)
		logger.WithField("data", data).Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.WithError(err).Error("failed to create rent")
		return nil, err
	}

	return result, nil
}

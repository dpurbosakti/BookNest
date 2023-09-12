package rent

import (
	"book-nest/clients/gomail"
	"book-nest/clients/midtrans"
	mb "book-nest/internal/models/book"
	mr "book-nest/internal/models/rent"
	mu "book-nest/internal/models/user"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RentService struct {
	RentRepository mr.RentRepository
	BookRepository mb.BookRepository
	UserRepository mu.UserRepository
	DB             *gorm.DB
	Gomail         *gomail.Gomail
	Midtrans       *midtrans.Midtrans
}

func NewRentService(rentRepository mr.RentRepository, bookRepository mb.BookRepository, userRepository mu.UserRepository, db *gorm.DB, gomail *gomail.Gomail, midtrans *midtrans.Midtrans) mr.RentService {
	return &RentService{
		RentRepository: rentRepository,
		BookRepository: bookRepository,
		UserRepository: userRepository,
		DB:             db,
		Gomail:         gomail,
		Midtrans:       midtrans,
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
	token := new(string)
	redirect_url := new(string)
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

		refId := srv.GenerateReferenceId(tx)
		data.ReferenceId = refId

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
		token, redirect_url, err = srv.Midtrans.CreatePayment(data)
		if err != nil {
			logger.WithError(err).Error("failed to create payment")
			return err
		}
		result = modelToResponse(resultRepo)
		result.Token = token
		result.RedirectURL = redirect_url
		logger.WithField("data", data).Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.WithError(err).Error("failed to create rent")
		return nil, err
	}

	return result, nil
}

func (srv *RentService) GenerateReferenceId(tx *gorm.DB) string {
	n := 7
	sb := strings.Builder{}
	sb.Grow(n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	refId := fmt.Sprintf("%s%s", time.Now().Format("20060102"), sb.String())
	rent, _ := srv.RentRepository.GetDetail(tx, refId)
	// check rent with reference id exist, if exist do GenerateReferenceId again
	if rent != nil {
		return srv.GenerateReferenceId(tx)
	}

	return refId
}

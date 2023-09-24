package rent

import (
	"book-nest/clients/gomail"
	"book-nest/clients/midtrans"
	mb "book-nest/internal/models/book"
	mr "book-nest/internal/models/rent"
	mu "book-nest/internal/models/user"
	ch "book-nest/utils/calendarhelper"
	"book-nest/utils/pagination"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	jwt5 "github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type RentService struct {
	RentRepository mr.RentRepository
	BookRepository mb.BookRepository
	UserRepository mu.UserRepository
	DB             *gorm.DB
	Gomail         *gomail.Gomail
	Midtrans       *midtrans.Midtrans
	mu             sync.Mutex
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
	data.UserId = userId
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

		token, redirect_url, err = srv.Midtrans.CreatePayment(data)
		if err != nil {
			logger.WithError(err).Error("failed to create payment")
			return err
		}
		result = modelToResponse(resultRepo)
		result.Token = token
		result.RedirectURL = redirect_url
		err = srv.Gomail.SendInvoice(result)
		if err != nil {
			logger.WithError(err).Error("failed to send invoice email")
			return err
		}
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
	var sb strings.Builder

	for {
		// Generate a random index based on the length of the charset
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err) // Handle the error appropriately in your application
		}

		// Use the random index to select a character from the charset
		randomChar := charset[idx.Int64()]

		// Append the random character to the result string
		sb.WriteByte(randomChar)

		if sb.Len() == refLength {
			break
		}
	}

	refId := fmt.Sprintf("%s%s", time.Now().Format("20060102"), sb.String())
	rent, _ := srv.RentRepository.GetDetail(tx, refId)
	// check rent with reference id exist, if exist do GenerateReferenceId again
	if rent != nil {
		return srv.GenerateReferenceId(tx)
	}

	return refId
}

func (srv *RentService) Update(input *mr.RentUpdateRequest) (*mr.RentResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "update",
		"scope": "rent service",
	})

	if input == nil {
		return nil, errors.New("input is nil")
	}

	logger.WithField("data", input).Info()
	result := new(mr.RentResponse)
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultRent, err := srv.RentRepository.GetDetail(tx, input.ReferenceId)
		if err != nil {
			logger.WithError(err).Error("failed to get rent data")
			return err
		}
		resultRent.PaymentStatus = input.PaymentStatus
		if input.PaymentStatus == "refund" {
			resultRent.Status = "rejected"
		}

		resultRepo, err := srv.RentRepository.Update(tx, resultRent)
		if err != nil {
			logger.WithError(err).Error("failed to update rent")
			return err
		}

		result = modelToResponse(resultRepo)

		if input.PaymentStatus == "settlement" {
			err := srv.Gomail.SendSuccessPayment(input, resultRepo)
			if err != nil {
				logger.WithError(err).Error("failed to send payment success email")
				return err
			}
		}

		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.WithError(err).Error("failed to create rent")
		return nil, err
	}

	return result, nil
}

func (srv *RentService) Accept(ctx *gin.Context, referenceId string) error {
	logger := logrus.WithFields(logrus.Fields{
		"func":         "accept",
		"scope":        "rent service",
		"reference_id": referenceId,
	})
	logger.Info()

	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultRent, err := srv.RentRepository.GetDetail(tx, referenceId)
		if err != nil {
			logger.WithError(err).Error("failed to get rent data")
			return err
		}

		if resultRent.PaymentStatus != PaymentSettlement {
			logger.Error("cannot accpet, payment status is not settlement")
			return errors.New("cannot accpet, payment status is not settlement")
		}

		if resultRent.Status == "rejected" {
			logger.Error("cannot accept, rent already rejected")
			return errors.New("cannot accept, rent already rejected")
		}

		resultRent.Book.AvailableAt = &resultRent.ReturnedDate
		resultRent.Book.IsAvailable = false

		_, err = srv.BookRepository.Update(tx, resultRent.Book)
		if err != nil {
			logger.WithError(err).Error("failed to update book")
			return err
		}

		resultRent.Status = "accepted"
		_, err = srv.RentRepository.Update(tx, resultRent)
		if err != nil {
			logger.WithError(err).Error("failed to update rent")
			return err
		}

		oauthToken := oauth2.Token{
			AccessToken: resultRent.User.OauthAccessToken,
		}

		userData := ctx.MustGet("userData").(jwt5.MapClaims)
		adminEmail := userData["email"].(string)
		adminName := userData["name"].(string)

		client := ch.GetClient(&oauthToken)
		calService, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
		if err != nil {
			logger.WithError(err).Error("failed to create calendar service")
			return err
		}

		event := &calendar.Event{
			Summary:     "Rent returned date",
			Description: "The day to return the book you rented",
			Start: &calendar.EventDateTime{
				DateTime: resultRent.ReturnedDate.Format(time.RFC3339),
			},
			End: &calendar.EventDateTime{
				DateTime: resultRent.ReturnedDate.Add(1 * time.Hour).Format(time.RFC3339),
			},
			Creator: &calendar.EventCreator{
				Email:       adminEmail,
				DisplayName: adminName,
			},
			Attendees: []*calendar.EventAttendee{
				{Email: resultRent.User.Email, DisplayName: resultRent.User.Name},
			},
		}

		_, err = calService.Events.Insert("primary", event).Do()
		if err != nil {
			logger.WithError(err).Error("failed to insert event into calendar")
			return err
		}
		logger.Info("end of db transaction")
		return nil
	})

	if err != nil {
		logger.WithError(err).Error("failed to create rent")
		return err
	}

	return nil
}

func (srv *RentService) Reject(ctx *gin.Context, referenceId string) error {
	logger := logrus.WithFields(logrus.Fields{
		"func":         "reject",
		"scope":        "rent service",
		"reference_id": referenceId,
	})
	logger.Info()

	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultRent, err := srv.RentRepository.GetDetail(tx, referenceId)
		if err != nil {
			logger.WithError(err).Error("failed to get rent data")
			return err
		}

		if resultRent.PaymentStatus != PaymentSettlement {
			logger.Error("cannot reject, payment status is not settlement")
			return errors.New("cannot reject, payment status is not settlement")
		}

		srv.mu.Lock()
		defer srv.mu.Unlock()

		res, err := srv.Midtrans.Refund(resultRent)
		if err != nil {
			logger.WithError(err).Error("failed to do refund")
			return err
		}

		resultRent.Status = "rejected"
		_, err = srv.RentRepository.Update(tx, resultRent)
		if err != nil {
			logger.WithError(err).Error("failed to update rent")
			return err
		}

		err = srv.Gomail.SendRefundedPayment(res, resultRent)
		if err != nil {
			logger.WithError(err).Error("failed to send payment refunded email")
			return err
		}

		logger.Info("end of db transaction")
		return nil
	})

	if err != nil {
		logger.WithError(err).Error("failed to reject rent")
		return err
	}

	return nil
}

func (srv *RentService) GetDetail(referenceId string) (*mr.RentResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":         "get_detail",
		"scope":        "rent service",
		"reference_id": referenceId,
	})
	logger.Info()

	result := new(mr.RentResponse)

	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")

		resultRent, err := srv.RentRepository.GetDetail(tx, referenceId)
		if err != nil {
			logger.WithError(err).Error("failed to get rent data")
			return err
		}
		result = modelToResponse(resultRent)

		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.WithError(err).Error("failed to get detail rent data")
		return nil, err
	}

	return result, nil
}

func (srv *RentService) GetList(page pagination.Pagination) (pagination.Pagination, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_list",
		"scope": "rent service",
	})
	var result pagination.Pagination
	logger.Info("data page", page)
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultRepo, err := srv.RentRepository.GetList(tx, page)
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

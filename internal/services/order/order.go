package order

import (
	"book-nest/clients/gomail"
	"book-nest/clients/midtrans"
	"book-nest/internal/constant"
	i "book-nest/internal/interfaces"
	mo "book-nest/internal/models/order"
	ch "book-nest/utils/calendarhelper"
	eh "book-nest/utils/errorhelper"
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

type OrderService struct {
	OrderRepository i.OrderRepository
	BookRepository  i.BookRepository
	UserRepository  i.UserRepository
	DB              *gorm.DB
	Gomail          *gomail.Gomail
	Midtrans        *midtrans.Midtrans
	mu              sync.Mutex
}

func NewOrderService(orderRepository i.OrderRepository, bookRepository i.BookRepository, userRepository i.UserRepository, db *gorm.DB, gomail *gomail.Gomail, midtrans *midtrans.Midtrans) i.OrderService {
	return &OrderService{
		OrderRepository: orderRepository,
		BookRepository:  bookRepository,
		UserRepository:  userRepository,
		DB:              db,
		Gomail:          gomail,
		Midtrans:        midtrans,
	}
}

const (
	userScope  = "user"
	bookScope  = "book"
	orderScope = "order"
)

func (srv *OrderService) Create(input *mo.OrderCreateRequest, userId uuid.UUID) (*mo.OrderResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "create",
		"scope": "order service",
	})

	if input == nil {
		return nil, errors.New("input is nil")
	}

	logger.WithField("data", input)
	result := new(mo.OrderResponse)
	data := requestToModel(input)
	data.UserId = userId
	token := new(string)
	redirect_url := new(string)
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.WithField("data", data).Info("db transaction begin")
		resultBook, err := srv.BookRepository.GetDetail(tx, data.BookId)
		if err != nil {
			eh.FailedGetDetail(logger, err, bookScope)
			return err
		}
		data.Book = resultBook

		resultUser, err := srv.UserRepository.GetDetail(tx, data.UserId)
		if err != nil {
			eh.FailedGetDetail(logger, err, userScope)
			return err
		}
		data.User = resultUser

		fee := calculateRentFee(data.GetDaysBetween(), data.Book.RentFeePerDay)
		data.Fee = fee

		refId := srv.GenerateReferenceId(tx)
		data.ReferenceId = refId

		resultRepo, err := srv.OrderRepository.Create(tx, data)
		if err != nil {
			eh.FailedCreate(logger, err, orderScope)
			return err
		}

		token, redirect_url, err = srv.Midtrans.CreatePayment(data)
		if err != nil {
			eh.FailedCreate(logger, err, "payment")
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
		eh.FailedCreate(logger, err, orderScope)
		return nil, err
	}

	return result, nil
}

func (srv *OrderService) GenerateReferenceId(tx *gorm.DB) string {
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
	order, _ := srv.OrderRepository.GetDetail(tx, refId)
	// check rent with reference id exist, if exist do GenerateReferenceId again
	if order != nil {
		return srv.GenerateReferenceId(tx)
	}

	return refId
}

func (srv *OrderService) Update(input *mo.OrderUpdateRequest) (*mo.OrderResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "update",
		"scope": "order service",
	})

	if input == nil {
		return nil, errors.New("input is nil")
	}

	logger.WithField("data", input).Info()
	result := new(mo.OrderResponse)
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultOrder, err := srv.OrderRepository.GetDetail(tx, input.ReferenceId)
		if err != nil {
			eh.FailedGetDetail(logger, err, orderScope)
			return err
		}
		resultOrder.PaymentStatus = input.PaymentStatus
		if input.PaymentStatus == constant.PaymentRefund {
			resultOrder.Status = "rejected"
		}

		resultRepo, err := srv.OrderRepository.Update(tx, resultOrder)
		if err != nil {
			eh.FailedUpdate(logger, err, orderScope)
			return err
		}

		result = modelToResponse(resultRepo)

		if input.PaymentStatus == constant.PaymentSettlement {
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
		eh.FailedUpdate(logger, err, orderScope)
		return nil, err
	}

	return result, nil
}

func (srv *OrderService) Accept(ctx *gin.Context, referenceId string) error {
	logger := logrus.WithFields(logrus.Fields{
		"func":         "accept",
		"scope":        "order service",
		"reference_id": referenceId,
	})
	logger.Info()

	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultOrder, err := srv.OrderRepository.GetDetail(tx, referenceId)
		if err != nil {
			eh.FailedGetDetail(logger, err, orderScope)
			return err
		}

		if resultOrder.PaymentStatus != constant.PaymentSettlement {
			logger.Error("cannot accpet, payment status is not settlement")
			return errors.New("cannot accpet, payment status is not settlement")
		}

		if resultOrder.Status == "rejected" {
			logger.Error("cannot accept, rent already rejected")
			return errors.New("cannot accept, rent already rejected")
		}

		resultOrder.Book.AvailableAt = &resultOrder.ReturnedDate
		resultOrder.Book.IsAvailable = false

		_, err = srv.BookRepository.Update(tx, resultOrder.Book)
		if err != nil {
			eh.FailedUpdate(logger, err, orderScope)
			return err
		}

		resultOrder.Status = "accepted"
		_, err = srv.OrderRepository.Update(tx, resultOrder)
		if err != nil {
			eh.FailedUpdate(logger, err, orderScope)
			return err
		}

		oauthToken := oauth2.Token{
			AccessToken: resultOrder.User.OauthAccessToken,
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
				DateTime: resultOrder.ReturnedDate.Format(time.RFC3339),
			},
			End: &calendar.EventDateTime{
				DateTime: resultOrder.ReturnedDate.Add(1 * time.Hour).Format(time.RFC3339),
			},
			Creator: &calendar.EventCreator{
				Email:       adminEmail,
				DisplayName: adminName,
			},
			Attendees: []*calendar.EventAttendee{
				{Email: resultOrder.User.Email, DisplayName: resultOrder.User.Name},
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
		logger.WithError(err).Error("failed to accept order")
		return err
	}

	return nil
}

func (srv *OrderService) Reject(ctx *gin.Context, referenceId string) error {
	logger := logrus.WithFields(logrus.Fields{
		"func":         "reject",
		"scope":        "order service",
		"reference_id": referenceId,
	})
	logger.Info()

	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultOrder, err := srv.OrderRepository.GetDetail(tx, referenceId)
		if err != nil {
			eh.FailedGetDetail(logger, err, orderScope)
			return err
		}

		if resultOrder.PaymentStatus != constant.PaymentSettlement {
			logger.Error("cannot reject, payment status is not settlement")
			return errors.New("cannot reject, payment status is not settlement")
		}

		srv.mu.Lock()
		defer srv.mu.Unlock()

		res, err := srv.Midtrans.Refund(resultOrder)
		if err != nil {
			logger.WithError(err).Error("failed to do refund")
			return err
		}

		resultOrder.Status = "rejected"
		_, err = srv.OrderRepository.Update(tx, resultOrder)
		if err != nil {
			eh.FailedUpdate(logger, err, orderScope)
			return err
		}

		err = srv.Gomail.SendRefundedPayment(res, resultOrder)
		if err != nil {
			logger.WithError(err).Error("failed to send payment refunded email")
			return err
		}

		logger.Info("end of db transaction")
		return nil
	})

	if err != nil {
		logger.WithError(err).Error("failed to reject order")
		return err
	}

	return nil
}

func (srv *OrderService) GetDetail(referenceId string) (*mo.OrderResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":         "get_detail",
		"scope":        "order service",
		"reference_id": referenceId,
	})
	logger.Info()

	result := new(mo.OrderResponse)

	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")

		resultOrder, err := srv.OrderRepository.GetDetail(tx, referenceId)
		if err != nil {
			eh.FailedGetDetail(logger, err, orderScope)
			return err
		}
		result = modelToResponse(resultOrder)

		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		eh.FailedGetDetail(logger, err, orderScope)
		return nil, err
	}

	return result, nil
}

func (srv *OrderService) GetList(page pagination.Pagination) (pagination.Pagination, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_list",
		"scope": "order service",
	})
	var result pagination.Pagination
	logger.Info("data page", page)
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultRepo, err := srv.OrderRepository.GetList(tx, page)
		if err != nil {
			eh.FailedGetList(logger, err, orderScope)
			return err
		}
		result = resultRepo
		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		eh.FailedGetList(logger, err, orderScope)
		return result, err
	}

	return result, nil
}

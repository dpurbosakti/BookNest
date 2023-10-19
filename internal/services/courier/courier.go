package courier

import (
	"book-nest/clients/biteship"
	i "book-nest/internal/interfaces"
	mc "book-nest/internal/models/courier"
	eh "book-nest/utils/errorhelper"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CourierService struct {
	CourierRepository i.CourierRepository
	AddressRepository i.AddressRepository
	BookRepository    i.BookRepository
	OrderRepository   i.OrderRepository
	DB                *gorm.DB
	Biteship          *biteship.Biteship
}

func NewCourierService(courierRepository i.CourierRepository, addressRepository i.AddressRepository, orderRepository i.OrderRepository, bookRepository i.BookRepository, db *gorm.DB, biteship *biteship.Biteship) i.CourierService {
	return &CourierService{
		CourierRepository: courierRepository,
		AddressRepository: addressRepository,
		OrderRepository:   orderRepository,
		BookRepository:    bookRepository,
		DB:                db,
		Biteship:          biteship,
	}
}

const (
	bookScope    = "book"
	addressScope = "address"
	courierScope = "courier"
)

func (srv *CourierService) GetBiteshipCourier() error {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_biteship_courier",
		"scope": "courier service",
	})
	biteshipResp, err := srv.Biteship.GetListCourier()
	if err != nil {
		eh.FailedGetList(logger, err, courierScope)
		return err
	}

	instantCouriers := GetInstantCourierOnly(biteshipResp.Couriers)

	err = srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.WithField("data", biteshipResp.Couriers).Info("db transaction begin")
		err := srv.CourierRepository.Create(tx, &instantCouriers)
		if err != nil {
			eh.FailedCreate(logger, err, courierScope)
			return err
		}

		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		eh.FailedGetList(logger, err, courierScope)
		return err
	}

	return nil
}

func (srv *CourierService) GetList() ([]mc.Courier, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_list",
		"scope": "courier service",
	})

	var result []mc.Courier
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultRepo, err := srv.CourierRepository.GetList(tx)
		if err != nil {
			eh.FailedGetList(logger, err, courierScope)
			return err
		}
		if resultRepo == nil {
			return errors.New("couriers not found")
		}
		result = resultRepo
		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		eh.FailedGetList(logger, err, courierScope)
		return nil, err
	}

	return result, nil

}

func (srv *CourierService) CheckRates(userId uuid.UUID, input *mc.CheckRatesRequest) (*biteship.BiteshipCheckRatesResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":    "check_rates",
		"scope":   "courier service",
		"user_id": userId,
	})
	logger.Info()
	result := new(biteship.BiteshipCheckRatesResponse)

	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultBook, err := srv.BookRepository.GetDetail(tx, input.BookId)
		if err != nil {
			eh.FailedGetDetail(logger, err, bookScope)
			return err
		}
		if resultBook == nil {
			return errors.New("book not found")
		}

		resultAddress, err := srv.AddressRepository.GetDetail(tx, input.AddressId)
		if err != nil {
			eh.FailedGetDetail(logger, err, addressScope)
			return err
		}

		if resultAddress.UserId != userId {
			logger.Warn("Unauthorized")
			return errors.New("Unauthorized")
		}

		resultCourier, err := srv.CourierRepository.GetList(tx)
		if err != nil {
			eh.FailedGetList(logger, err, courierScope)
			return err
		}
		payload := checkRatesPayloadBuilder(resultBook, *resultAddress, getCouriersName(resultCourier))

		res, err := srv.Biteship.CheckRates(payload)
		if err != nil {
			logger.WithError(err).Error("failed to check rates")
			return err
		}
		result = res
		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.Error("failed to check rates")
		return nil, err
	}

	return result, nil
}

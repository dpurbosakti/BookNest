package courier

import (
	"book-nest/clients/biteship"
	i "book-nest/internal/interfaces"
	mc "book-nest/internal/models/courier"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CourierService struct {
	CourierRepository i.CourierRepository
	AddressRepository i.AddressRepository
	BookRepository    i.BookRepository
	RentRepository    i.RentRepository
	DB                *gorm.DB
	Biteship          *biteship.Biteship
}

func NewCourierService(courierRepository i.CourierRepository, addressRepository i.AddressRepository, rentRepository i.RentRepository, bookRepository i.BookRepository, db *gorm.DB, biteship *biteship.Biteship) i.CourierService {
	return &CourierService{
		CourierRepository: courierRepository,
		AddressRepository: addressRepository,
		RentRepository:    rentRepository,
		BookRepository:    bookRepository,
		DB:                db,
		Biteship:          biteship,
	}
}

func (srv *CourierService) GetBiteshipCourier() error {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_biteship_courier",
		"scope": "courier service",
	})
	biteshipResp, err := srv.Biteship.GetListCourier()
	if err != nil {
		logger.WithError(err).Error("failed to get list courier")
		return err
	}

	instantCouriers := GetInstantCourierOnly(biteshipResp.Couriers)

	err = srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.WithField("data", biteshipResp.Couriers).Info("db transaction begin")
		err := srv.CourierRepository.Create(tx, &instantCouriers)
		if err != nil {
			logger.WithError(err).Error("failed to create courier")
			return err
		}

		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.WithError(err).Error("failed to create and get list courier")
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
			logger.WithError(err).Error("failed to get list")
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
		logger.Error("failed to get list")
		return nil, err
	}

	return result, nil

}

func (srv *CourierService) CheckRates(userId uuid.UUID, bookId uint) (*biteship.BiteshipCheckRatesResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":    "check_rates",
		"scope":   "courier service",
		"user_id": userId,
	})
	logger.Info()
	result := new(biteship.BiteshipCheckRatesResponse)

	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultBook, err := srv.BookRepository.GetDetail(tx, bookId)
		if err != nil {
			logger.WithError(err).Error("failed to get detail book")
			return err
		}
		if resultBook == nil {
			return errors.New("book not found")
		}

		resultAddress, err := srv.AddressRepository.GetByUserId(tx, userId)
		if err != nil {
			logger.WithError(err).Error("failed to get detail address")
			return err
		}

		resultCourier, err := srv.CourierRepository.GetList(tx)
		if err != nil {
			logger.WithError(err).Error("failed to get list courier")
			return err
		}
		payload := checkRatesPayloadBuilder(resultBook, *resultAddress, getCouriersName(resultCourier))

		res, err := srv.Biteship.CheckRates(payload)
		if err != nil {
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

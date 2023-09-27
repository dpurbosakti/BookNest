package courier

import (
	"book-nest/clients/biteship"
	mc "book-nest/internal/models/courier"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CourierService struct {
	CourierRepository mc.CourierRepository
	DB                *gorm.DB
	Biteship          *biteship.Biteship
}

func NewCourierService(courierRepository mc.CourierRepository, db *gorm.DB, biteship *biteship.Biteship) mc.CourierService {
	return &CourierService{
		CourierRepository: courierRepository,
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

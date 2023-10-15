package address

import (
	i "book-nest/internal/interfaces"
	mad "book-nest/internal/models/address"
	eh "book-nest/utils/errorhelper"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AddressService struct {
	AddressRepository i.AddressRepository
	DB                *gorm.DB
}

func NewAddressService(addressRepository i.AddressRepository, db *gorm.DB) i.AddressService {
	return &AddressService{
		AddressRepository: addressRepository,
		DB:                db,
	}
}

func (srv *AddressService) Create(input *mad.AddressCreateRequest) (*mad.AddressResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "create",
		"scope": "address service",
	})

	if input == nil {
		return nil, errors.New("input is nil")
	}

	logger.WithField("data", input)
	result := new(mad.AddressResponse)
	data := requestToModel(input)
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.WithField("data", data).Info("db transaction begin")
		resultRepo, err := srv.AddressRepository.Create(tx, data)
		if err != nil {
			eh.FailedCreate(logger, err, "address")
			return err
		}

		result = modelToResponse(resultRepo)
		logger.WithField("data", data).Info("end of db transaction")
		return nil
	})
	if err != nil {
		eh.FailedCreate(logger, err, "address")
		return nil, err
	}

	return result, nil
}

func (srv *AddressService) GetDetail(addressId uint) (*mad.AddressResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_detail",
		"scope": "address service",
	})
	logger.WithField("id", addressId)
	result := new(mad.AddressResponse)
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultRepo, err := srv.AddressRepository.GetDetail(tx, addressId)
		if err != nil {
			eh.FailedGetDetail(logger, err, "address")
			return err
		}
		result = modelToResponse(resultRepo)
		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		eh.FailedGetDetail(logger, err, "address")
		return result, err
	}

	return result, nil
}

func (srv *AddressService) Update(input *mad.AddressUpdateRequest, addressId uint, userId uuid.UUID) (*mad.AddressResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "update",
		"scope": "address service",
	})
	logger.WithField("data", input)

	result := new(mad.AddressResponse)

	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultGet, err := srv.AddressRepository.GetDetail(tx, addressId)
		if err != nil {
			logger.WithError(err).Error("failed to get detail")
			return err
		}

		if resultGet.UserId != userId {
			logger.Warn("Unauthorized")
			return errors.New("Unauthorized")
		}
		resultGet.Copier(input)
		resultUpdate, err := srv.AddressRepository.Update(tx, resultGet)
		if err != nil {
			logger.WithError(err).Error("failed to update data")
			return err
		}
		result = modelToResponse(resultUpdate)
		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.Error("failed to update data")
		return result, err
	}

	return result, nil
}

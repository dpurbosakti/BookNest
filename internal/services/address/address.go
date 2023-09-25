package address

import (
	mad "book-nest/internal/models/address"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AddressService struct {
	AddressRepository mad.AddressRepository
	DB                *gorm.DB
}

func NewAddressService(addressRepository mad.AddressRepository, db *gorm.DB) mad.AddressService {
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
			logger.WithError(err).Error("failed to create address")
			return err
		}

		result = modelToResponse(resultRepo)
		logger.WithField("data", data).Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.WithError(err).Error("failed to create address")
		return nil, err
	}

	return result, nil
}

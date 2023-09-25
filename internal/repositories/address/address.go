package address

import (
	mad "book-nest/internal/models/address"

	"gorm.io/gorm"
)

type AddressRepository struct {
}

func NewAddressRepository() mad.AddressRepository {
	return &AddressRepository{}
}

func (repo *AddressRepository) Create(tx *gorm.DB, input *mad.Address) (*mad.Address, error) {
	resultcreate := tx.Create(&input)
	if resultcreate.Error != nil {
		return nil, resultcreate.Error
	}

	return input, nil
}

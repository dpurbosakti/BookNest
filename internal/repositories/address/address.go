package address

import (
	i "book-nest/internal/interfaces"
	mad "book-nest/internal/models/address"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AddressRepository struct {
}

func NewAddressRepository() i.AddressRepository {
	return &AddressRepository{}
}

func (repo *AddressRepository) Create(tx *gorm.DB, input *mad.Address) (*mad.Address, error) {
	resultcreate := tx.Create(&input)
	if resultcreate.Error != nil {
		return nil, resultcreate.Error
	}

	return input, nil
}

func (repo *AddressRepository) Update(tx *gorm.DB, input *mad.Address) (*mad.Address, error) {
	result := tx.Save(&input)
	if result.Error != nil {
		return nil, errors.New("error updating your data")
	}

	return input, nil
}

func (repo *AddressRepository) GetDetail(tx *gorm.DB, addressId uint) (*mad.Address, error) {
	book := new(mad.Address)
	result := tx.First(&book, addressId)
	if result.Error != nil {
		return nil, fmt.Errorf("user id %d not found", addressId)
	}

	return book, nil
}

func (repo *AddressRepository) GetByUserId(tx *gorm.DB, userId uuid.UUID) (*mad.Address, error) {
	address := new(mad.Address)
	result := tx.Where("user_id = ?", userId).First(&address)
	if result.Error != nil {
		return nil, result.Error
	}

	return address, nil
}

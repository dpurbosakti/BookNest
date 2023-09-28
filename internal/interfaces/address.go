package interfaces

import (
	mad "book-nest/internal/models/address"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AddressHandler interface {
	Create(c *gin.Context)
	GetDetail(c *gin.Context)
	Update(c *gin.Context)
}

type AddressService interface {
	Create(input *mad.AddressCreateRequest) (*mad.AddressResponse, error)
	Update(input *mad.AddressUpdateRequest, addressId uint, userId uuid.UUID) (*mad.AddressResponse, error)
	GetDetail(addressId uint) (*mad.AddressResponse, error)
}

type AddressRepository interface {
	Create(tx *gorm.DB, input *mad.Address) (*mad.Address, error)
	Update(tx *gorm.DB, input *mad.Address) (*mad.Address, error)
	GetDetail(tx *gorm.DB, addressId uint) (*mad.Address, error)
	GetByUserId(tx *gorm.DB, userId uuid.UUID) (*mad.Address, error)
}

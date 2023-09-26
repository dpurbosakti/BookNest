package address

import (
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
	Create(input *AddressCreateRequest) (*AddressResponse, error)
	Update(input *AddressUpdateRequest, addressId uint, userId uuid.UUID) (*AddressResponse, error)
	GetDetail(addressId uint) (*AddressResponse, error)
}

type AddressRepository interface {
	Create(tx *gorm.DB, input *Address) (*Address, error)
	Update(tx *gorm.DB, input *Address) (*Address, error)
	GetDetail(tx *gorm.DB, addressId uint) (*Address, error)
}

package address

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AddressHandler interface {
	Create(c *gin.Context)
}

type AddressService interface {
	Create(input *AddressCreateRequest) (*AddressResponse, error)
}

type AddressRepository interface {
	Create(tx *gorm.DB, input *Address) (*Address, error)
}

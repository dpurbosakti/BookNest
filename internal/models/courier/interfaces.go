package courier

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CourierHandler interface {
	GetBiteshipCourier(c *gin.Context)
	GetList(c *gin.Context)
}

type CourierService interface {
	GetBiteshipCourier() error
	GetList() ([]Courier, error)
}

type CourierRepository interface {
	Create(tx *gorm.DB, input *[]Courier) error
	GetList(tx *gorm.DB) ([]Courier, error)
}

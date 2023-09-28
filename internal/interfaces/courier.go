package interfaces

import (
	"book-nest/clients/biteship"
	mc "book-nest/internal/models/courier"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CourierHandler interface {
	GetBiteshipCourier(c *gin.Context)
	GetList(c *gin.Context)
}

type CourierService interface {
	GetBiteshipCourier() error
	GetList() ([]mc.Courier, error)
	CheckRates(referenceId string) (*biteship.BiteshipCheckRatesResponse, error)
}

type CourierRepository interface {
	Create(tx *gorm.DB, input *[]mc.Courier) error
	GetList(tx *gorm.DB) ([]mc.Courier, error)
}

package rent

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RentHandler interface {
	Create(c *gin.Context)
}

type RentService interface {
	Create(input *RentCreateRequest, userId uuid.UUID) (*RentResponse, error)
}
type RentRepository interface {
	Create(tx *gorm.DB, input *Rent) (*Rent, error)
}

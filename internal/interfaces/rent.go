package interfaces

import (
	mr "book-nest/internal/models/rent"
	"book-nest/utils/pagination"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RentHandler interface {
	Create(c *gin.Context)
	Accept(c *gin.Context)
	Reject(c *gin.Context)
	GetDetail(c *gin.Context)
}

type RentService interface {
	Create(input *mr.RentCreateRequest, userId uuid.UUID) (*mr.RentResponse, error)
	GenerateReferenceId(tx *gorm.DB) string
	Update(input *mr.RentUpdateRequest) (*mr.RentResponse, error)
	Accept(ctx *gin.Context, referenceId string) error
	Reject(ctx *gin.Context, referenceId string) error
	GetDetail(referenceId string) (*mr.RentResponse, error)
	GetList(page pagination.Pagination) (pagination.Pagination, error)
}
type RentRepository interface {
	Create(tx *gorm.DB, input *mr.Rent) (*mr.Rent, error)
	GetDetail(tx *gorm.DB, referenceId string) (*mr.Rent, error)
	Update(tx *gorm.DB, input *mr.Rent) (*mr.Rent, error)
	GetList(tx *gorm.DB, page pagination.Pagination) (pagination.Pagination, error)
}

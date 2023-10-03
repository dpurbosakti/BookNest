package interfaces

import (
	mo "book-nest/internal/models/order"
	"book-nest/utils/pagination"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderHandler interface {
	Create(c *gin.Context)
	Accept(c *gin.Context)
	Reject(c *gin.Context)
	GetDetail(c *gin.Context)
}

type OrderService interface {
	Create(input *mo.OrderCreateRequest, userId uuid.UUID) (*mo.OrderResponse, error)
	GenerateReferenceId(tx *gorm.DB) string
	Update(input *mo.OrderUpdateRequest) (*mo.OrderResponse, error)
	Accept(ctx *gin.Context, referenceId string) error
	Reject(ctx *gin.Context, referenceId string) error
	GetDetail(referenceId string) (*mo.OrderResponse, error)
	GetList(page pagination.Pagination) (pagination.Pagination, error)
}
type OrderRepository interface {
	Create(tx *gorm.DB, input *mo.Order) (*mo.Order, error)
	GetDetail(tx *gorm.DB, referenceId string) (*mo.Order, error)
	Update(tx *gorm.DB, input *mo.Order) (*mo.Order, error)
	GetList(tx *gorm.DB, page pagination.Pagination) (pagination.Pagination, error)
}

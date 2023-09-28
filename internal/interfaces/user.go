package interfaces

import (
	mu "book-nest/internal/models/user"
	"book-nest/utils/pagination"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler interface {
	Create(c *gin.Context)
	Verify(c *gin.Context)
	RefreshVerificationCode(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetDetail(c *gin.Context)
	GetList(c *gin.Context)
}

type UserService interface {
	Create(input *mu.UserCreateRequest) (*mu.UserResponse, error)
	Verify(input *mu.UserVerifyRequest) error
	RefreshVerificationCode(input *mu.UserVerificationCodeRequest) error
	Update(input *mu.UserUpdateRequest, userId uuid.UUID) (*mu.UserResponse, error)
	Delete(userId uuid.UUID) error
	GetDetail(userId uuid.UUID) (*mu.UserResponse, error)
	GetList(page pagination.Pagination) (pagination.Pagination, error)
}

type UserRepository interface {
	Create(tx *gorm.DB, input *mu.User) (*mu.User, error)
	CheckDuplicate(tx *gorm.DB, input *mu.User) error
	GetByEmail(tx *gorm.DB, email string) (*mu.User, error)
	Update(tx *gorm.DB, input *mu.User) (*mu.User, error)
	GetDetail(tx *gorm.DB, userId uuid.UUID) (*mu.User, error)
	GetList(tx *gorm.DB, page pagination.Pagination) (pagination.Pagination, error)
	Delete(tx *gorm.DB, userId uuid.UUID) error
}

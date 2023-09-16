package user

import (
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
	Create(input *UserCreateRequest) (*UserResponse, error)
	Verify(input *UserVerifyRequest) error
	RefreshVerificationCode(input *UserVerificationCodeRequest) error
	Update(input *UserUpdateRequest, userId uuid.UUID) (*UserResponse, error)
	Delete(userId uuid.UUID) error
	GetDetail(userId uuid.UUID) (*UserResponse, error)
	GetList(page pagination.Pagination) (pagination.Pagination, error)
}

type UserRepository interface {
	Create(tx *gorm.DB, input *User) (*User, error)
	CheckDuplicate(tx *gorm.DB, input *User) error
	GetByEmail(tx *gorm.DB, email string) (*User, error)
	Update(tx *gorm.DB, input *User) (*User, error)
	GetDetail(tx *gorm.DB, userId uuid.UUID) (*User, error)
	GetList(tx *gorm.DB, page pagination.Pagination) (pagination.Pagination, error)
	Delete(tx *gorm.DB, userId uuid.UUID) error
}

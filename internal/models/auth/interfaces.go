package auth

import (
	"book-nest/internal/models/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler interface {
	Login(c *gin.Context)
	GoogleLogin(c *gin.Context)
	GoogleCallback(c *gin.Context)
	TwitterLogin(c *gin.Context)
	TwitterCallback(c *gin.Context)
	GithubLogin(c *gin.Context)
	GithubCallback(c *gin.Context)
}

type AuthService interface {
	Login(input LoginRequest)
}

type AuthRepository interface {
	Login(tx *gorm.DB, input LoginRequest) (user.User, error)
}

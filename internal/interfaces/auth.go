package interfaces

import (
	ma "book-nest/internal/models/auth"
	mu "book-nest/internal/models/user"

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
	Login(input ma.LoginRequest) (*string, error)
	LoginByGoogle(input *ma.GoogleResponse) (*string, error)
}

type AuthRepository interface {
	Login(tx *gorm.DB, input ma.LoginRequest) (*mu.User, error)
}

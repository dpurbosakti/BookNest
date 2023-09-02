package auth

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	GoogleLogin(c *gin.Context)
	GoogleCallback(c *gin.Context)
}

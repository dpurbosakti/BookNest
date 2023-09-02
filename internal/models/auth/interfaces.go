package auth

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	GoogleLogin(c *gin.Context)
	GoogleCallback(c *gin.Context)
	TwitterLogin(c *gin.Context)
	TwitterCallback(c *gin.Context)
	GithubLogin(c *gin.Context)
	GithubCallback(c *gin.Context)
}

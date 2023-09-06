package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func AdminAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userData := ctx.MustGet("userData").(jwt5.MapClaims)
		role := userData["role"].(string)
		logger := logrus.WithFields(logrus.Fields{"func": "admin_authorization", "data": userData})
		logger.WithField("role", role).Info()
		if role != "admin" {
			logger.Warn("Unauthorized")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "You are not allowed to access this data",
			})
			return
		}

		ctx.Next()
	}
}

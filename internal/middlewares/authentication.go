package middlewares

import (
	"net/http"

	jh "book-nest/utils/jwthelper"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := jh.VerifyToken(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthenticated",
				"message": err.Error(),
			})
			return
		}
		ctx.Set("userData", token)
		ctx.Next()
	}
}

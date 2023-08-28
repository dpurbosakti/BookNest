package handlers

import "github.com/gin-gonic/gin"

// Ping use for health check
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

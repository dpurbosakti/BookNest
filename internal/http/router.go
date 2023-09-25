package http

import (
	"book-nest/internal/handlers"
	"book-nest/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(r *gin.Engine, db *gorm.DB) {
	p := NewPresenter(db)

	// ping
	r.GET("/ping", handlers.Ping)

	// auth
	authGroup := r.Group("/auth")
	authGroup.POST("/login", p.Auth.Login)
	authGroup.GET("/google/login", p.Auth.GoogleLogin)
	authGroup.GET("/google/callback", p.Auth.GoogleCallback)
	authGroup.GET("/twitter/login", p.Auth.TwitterLogin)
	authGroup.GET("/twitter/callback", p.Auth.TwitterCallback)
	authGroup.GET("/github/login", p.Auth.GithubLogin)
	authGroup.GET("/github/callback", p.Auth.GithubCallback)

	// users
	userGroup := r.Group("/users")
	userGroup.POST("", p.User.Create)
	userGroup.POST("/verify-email", p.User.Verify)
	userGroup.POST("/refreshcode", p.User.RefreshVerificationCode)
	userGroup.Use(middlewares.Authentication())
	userGroup.GET("", middlewares.AdminAuthorization(), p.User.GetList)
	userGroup.GET("/detail", p.User.GetDetail)
	userGroup.PUT("/update", p.User.Update)
	userGroup.DELETE("/delete", p.User.Delete)

	// books
	bookGroup := r.Group("/books")
	bookGroup.GET("/:id", p.Book.GetDetail)
	bookGroup.GET("", p.Book.GetList)
	bookGroup.Use(middlewares.Authentication())
	bookGroup.POST("", middlewares.AdminAuthorization(), p.Book.Create)
	bookGroup.PUT("/:id", middlewares.AdminAuthorization(), p.Book.Update)
	bookGroup.DELETE("/:id", middlewares.AdminAuthorization(), p.Book.Delete)

	// rents
	rentGroup := r.Group("/rents")
	rentGroup.POST("/midtrans/callback", p.Rent.MidtransCallback)
	rentGroup.Use(middlewares.Authentication())
	rentGroup.POST("", p.Rent.Create)
	rentGroup.POST("/:reference_id/accept", middlewares.AdminAuthorization(), p.Rent.Accept)
	rentGroup.POST("/:reference_id/reject", middlewares.AdminAuthorization(), p.Rent.Reject)
	rentGroup.GET("/:reference_id", p.Rent.GetDetail)
	rentGroup.GET("/:reference_id", middlewares.AdminAuthorization(), p.Rent.GetList)

	// addresss
	addressGroup := r.Group("/address")
	addressGroup.Use(middlewares.Authentication())
	rentGroup.POST("", p.Address.Create)
}

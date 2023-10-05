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
	bookGroup.POST("/:id/return", middlewares.AdminAuthorization(), p.Book.Return)

	// orders
	rentGroup := r.Group("/orders")
	rentGroup.POST("/midtrans/callback", p.Order.MidtransCallback)
	rentGroup.Use(middlewares.Authentication())
	rentGroup.POST("", p.Order.Create)
	rentGroup.POST("/:reference_id/accept", middlewares.AdminAuthorization(), p.Order.Accept)
	rentGroup.POST("/:reference_id/reject", middlewares.AdminAuthorization(), p.Order.Reject)
	rentGroup.GET("/:reference_id", p.Order.GetDetail)
	rentGroup.GET("", middlewares.AdminAuthorization(), p.Order.GetList)

	// addresss
	addressGroup := r.Group("/address")
	addressGroup.Use(middlewares.Authentication())
	addressGroup.POST("", p.Address.Create)
	addressGroup.GET("/:id", p.Address.GetDetail)
	addressGroup.PUT("/:id", p.Address.Update)

	// couriers
	couriersGroup := r.Group("/couriers")
	couriersGroup.Use(middlewares.Authentication())
	couriersGroup.GET("/biteship", p.Courier.GetBiteshipCourier)
	couriersGroup.GET("", p.Courier.GetList)
	couriersGroup.POST("/checkrates", p.Courier.CheckRates)
}

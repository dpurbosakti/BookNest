package http

import (
	"book-nest/internal/handlers"
	eh "book-nest/utils/emailhelper"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// auth
	authHdl "book-nest/internal/handlers/auth"

	// users
	userHdl "book-nest/internal/handlers/user"
	userRepo "book-nest/internal/repositories/user"
	userSrv "book-nest/internal/services/user"
)

func InitRouter(r *gin.Engine, db *gorm.DB) {

	// auth
	ah := authHdl.NewAuthHandler()

	// emailhelper
	email := eh.NewEmailHelper()

	// users
	ur := userRepo.NewUserRepository()
	us := userSrv.NewUserService(ur, db, email)
	uh := userHdl.NewUserHandler(us)

	// ping
	r.GET("/ping", handlers.Ping)

	// auth
	authGroup := r.Group("/auth")
	authGroup.GET("/google/login", ah.GoogleLogin)
	authGroup.GET("/google/callback", ah.GoogleCallback)
	authGroup.GET("/twitter/login", ah.TwitterLogin)
	authGroup.GET("/twitter/callback", ah.TwitterCallback)
	authGroup.GET("/github/login", ah.GithubLogin)
	authGroup.GET("/github/callback", ah.GithubCallback)

	// users
	userGroup := r.Group("/users")
	userGroup.POST("", uh.Create)
}

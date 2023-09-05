package http

import (
	"book-nest/internal/handlers"
	eh "book-nest/utils/emailhelper"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// auth
	authHdl "book-nest/internal/handlers/auth"
	authRepo "book-nest/internal/repositories/auth"
	authSrv "book-nest/internal/services/auth"

	// users
	userHdl "book-nest/internal/handlers/user"
	userRepo "book-nest/internal/repositories/user"
	userSrv "book-nest/internal/services/user"
)

func InitRouter(r *gin.Engine, db *gorm.DB) {

	// auth
	ar := authRepo.NewAuthRepository()
	as := authSrv.NewAuthService(ar, db)
	ah := authHdl.NewAuthHandler(as)

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
	authGroup.POST("/login", ah.Login)
	authGroup.GET("/google/login", ah.GoogleLogin)
	authGroup.GET("/google/callback", ah.GoogleCallback)
	authGroup.GET("/twitter/login", ah.TwitterLogin)
	authGroup.GET("/twitter/callback", ah.TwitterCallback)
	authGroup.GET("/github/login", ah.GithubLogin)
	authGroup.GET("/github/callback", ah.GithubCallback)

	// users
	userGroup := r.Group("/users")
	userGroup.POST("", uh.Create)
	userGroup.POST("/verify", uh.Verify)
	userGroup.POST("/refreshcode", uh.RefreshVerCode)
}

package http

import (
	"book-nest/internal/handlers"
	eh "book-nest/utils/emailhelper"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// users
	userHdl "book-nest/internal/handlers/user"
	userRepo "book-nest/internal/repositories/user"
	userSrv "book-nest/internal/services/user"
)

func InitRouter(r *gin.Engine, db *gorm.DB) {
	// emailhelper
	email := eh.NewEmailHelper()
	// users
	ur := userRepo.NewUserRepository()
	us := userSrv.NewUserService(ur, db, email)
	uh := userHdl.NewUserHandler(us)

	// ping
	r.GET("/ping", handlers.Ping)

	// users
	userGroup := r.Group("/users")
	userGroup.POST("", uh.Create)
}

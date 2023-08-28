package http

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// users
	userHdl "book-nest/internal/handlers/user"
	userRepo "book-nest/internal/repositories/user"
	userSrv "book-nest/internal/services/user"
)

func InitRouter(r *gin.Engine, db *gorm.DB) {
	// users
	ur := userRepo.NewUserRepository()
	us := userSrv.NewUserService(ur, db)
	uh := userHdl.NewUserHandler(us)

	// user
	userGroup := r.Group("/users")
	userGroup.POST("", uh.Create)
}

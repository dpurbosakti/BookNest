package migration

import (
	"book-nest/internal/models/user"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&user.User{})
}

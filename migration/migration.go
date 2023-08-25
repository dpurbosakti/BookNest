package migration

import (
	"book-nest/internal/models/user"
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatal(err)
	}
}

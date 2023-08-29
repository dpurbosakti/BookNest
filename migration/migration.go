package migration

import (
	"book-nest/internal/models/user"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&user.User{}); err != nil {
		return err
	}
	return nil
}

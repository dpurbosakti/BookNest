package migration

import (
	"book-nest/internal/models/book"
	"book-nest/internal/models/rent"
	"book-nest/internal/models/user"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&user.User{}, &book.Book{}, &rent.Rent{}); err != nil {
		return err
	}
	return nil
}

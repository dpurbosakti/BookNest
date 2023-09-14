package migration

import (
	"book-nest/internal/models/book"
	"book-nest/internal/models/rent"
	"book-nest/internal/models/user"

	"gorm.io/gorm"
)

func MigrateUp(db *gorm.DB) error {
	if err := db.AutoMigrate(&user.User{}, &book.Book{}, &rent.Rent{}); err != nil {
		return err
	}
	return nil
}

func MigrateDown(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&rent.Rent{}, &book.Book{}, &user.User{}); err != nil {
		return err
	}
	return nil
}

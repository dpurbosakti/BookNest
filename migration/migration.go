package migration

import (
	"book-nest/internal/models/address"
	"book-nest/internal/models/book"
	"book-nest/internal/models/courier"
	"book-nest/internal/models/rent"
	"book-nest/internal/models/user"

	"gorm.io/gorm"
)

func MigrateUp(db *gorm.DB) error {
	if err := db.AutoMigrate(&user.User{}, &book.Book{}, &rent.Rent{}, &address.Address{}, &courier.Courier{}); err != nil {
		return err
	}
	return nil
}

func MigrateDown(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&address.Address{}, &courier.Courier{}, &rent.Rent{}, &book.Book{}, &user.User{}); err != nil {
		return err
	}
	return nil
}

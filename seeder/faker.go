package seeder

import (
	"book-nest/internal/models/book"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v4"
	"gorm.io/gorm"
)

var bookTitles = []string{
	"The Great Gatsby",
	"To Kill a Mockingbird",
	"Pride and Prejudice",
	"1984",
	"The Catcher in the Rye",
	"Brave New World",
	"The Lord of the Rings",
	"Harry Potter and the Sorcerer's Stone",
	"Animal Farm",
	"War and Peace",
	"Crime and Punishment",
	"The Hobbit",
	"Fahrenheit 451",
	"Frankenstein",
	"The Grapes of Wrath",
	"Dracula",
	"The Shining",
	"The Hitchhiker's Guide to the Galaxy",
	"Moby-Dick",
	"The Odyssey",
}

func generateFee() float64 {
	// Generate a random integer between 1000 and 5000 (inclusive)
	randomInt := rand.Intn(4001) + 1000 // 4001 is the range (5000 - 1000 + 1)

	// Convert the integer to a float64
	randomFloat := float64(randomInt)
	return randomFloat
}

func generateTitle() string {
	var title string
_:
	faker.FakeData(&title)
	return title
}

func bookFaker(db *gorm.DB, n uint) (result []book.Book) {
	for i := 1; i < int(n); i++ {
		tmp := book.Book{
			Id:            uint(i),
			Title:         bookTitles[i],
			Author:        faker.FirstNameFemale(),
			RentFeePerDay: generateFee(),
			IsAvailable:   true,
			AvailableAt:   nil,
			CreatedAt:     time.Time{},
			UpdatedAt:     time.Time{},
			DeletedAt:     gorm.DeletedAt{},
		}
		result = append(result, tmp)
	}
	return result
}

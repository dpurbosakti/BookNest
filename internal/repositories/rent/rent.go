package rent

import (
	mr "book-nest/internal/models/rent"

	"gorm.io/gorm"
)

type RentRepository struct {
}

func NewRentRepository() mr.RentRepository {
	return &RentRepository{}
}

func (repo *RentRepository) Create(tx *gorm.DB, input *mr.Rent) (*mr.Rent, error) {
	resultcreate := tx.Create(&input)
	if resultcreate.Error != nil {
		return nil, resultcreate.Error
	}

	return input, nil
}

package rent

import (
	mr "book-nest/internal/models/rent"
	"fmt"

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

func (repo *RentRepository) GetDetail(tx *gorm.DB, referenceId string) (*mr.Rent, error) {
	rent := new(mr.Rent)
	rent.ReferenceId = referenceId
	result := tx.First(&rent)
	if result.Error != nil {
		return nil, fmt.Errorf("user id %s not found", referenceId)
	}

	return rent, nil
}

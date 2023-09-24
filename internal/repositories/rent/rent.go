package rent

import (
	mr "book-nest/internal/models/rent"
	"book-nest/utils/pagination"
	"errors"
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

	result := tx.Preload("User").Preload("Book").Where("reference_id = ?", referenceId).First(&rent)
	if result.Error != nil {
		return nil, fmt.Errorf("user id %s not found", referenceId)
	}

	if result.RowsAffected < 1 {
		return nil, nil
	}

	return rent, nil
}

func (repo *RentRepository) Update(tx *gorm.DB, input *mr.Rent) (*mr.Rent, error) {
	result := tx.Save(&input)
	if result.Error != nil {
		return nil, errors.New("error updating your data")
	}

	return input, nil
}

func (repo *RentRepository) GetList(tx *gorm.DB, page pagination.Pagination) (pagination.Pagination, error) {
	var rents []mr.Rent

	tx.Scopes(pagination.Paginate(rents, &page, tx)).Find(&rents)
	page.Rows = rents

	return page, nil
}

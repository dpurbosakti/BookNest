package courier

import (
	mc "book-nest/internal/models/courier"
	"errors"

	"gorm.io/gorm"
)

type CourierRepository struct {
}

func NewCourierRepository() mc.CourierRepository {
	return &CourierRepository{}
}

func (repo *CourierRepository) Create(tx *gorm.DB, input *[]mc.Courier) error {
	resultcreate := tx.Create(&input)
	if resultcreate.Error != nil {
		return resultcreate.Error
	}

	return nil
}

func (repo *CourierRepository) GetList(tx *gorm.DB) ([]mc.Courier, error) {
	var couriers []mc.Courier

	if err := tx.Find(&couriers).Error; err != nil {
		return nil, err
	}

	if len(couriers) < 1 {
		return nil, errors.New("couriers not found")
	}

	return couriers, nil
}

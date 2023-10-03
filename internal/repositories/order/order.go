package order

import (
	i "book-nest/internal/interfaces"
	mo "book-nest/internal/models/order"
	"book-nest/utils/pagination"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type OrderRepository struct {
}

func NewOrderRepository() i.OrderRepository {
	return &OrderRepository{}
}

func (repo *OrderRepository) Create(tx *gorm.DB, input *mo.Order) (*mo.Order, error) {
	resultcreate := tx.Create(&input)
	if resultcreate.Error != nil {
		return nil, resultcreate.Error
	}

	return input, nil
}

func (repo *OrderRepository) GetDetail(tx *gorm.DB, referenceId string) (*mo.Order, error) {
	order := new(mo.Order)

	result := tx.Preload("User").Preload("Book").Where("reference_id = ?", referenceId).First(&order)
	if result.Error != nil {
		return nil, fmt.Errorf("user id %s not found", referenceId)
	}

	if result.RowsAffected < 1 {
		return nil, nil
	}

	return order, nil
}

func (repo *OrderRepository) Update(tx *gorm.DB, input *mo.Order) (*mo.Order, error) {
	result := tx.Save(&input)
	if result.Error != nil {
		return nil, errors.New("error updating your data")
	}

	return input, nil
}

func (repo *OrderRepository) GetList(tx *gorm.DB, page pagination.Pagination) (pagination.Pagination, error) {
	var orders []mo.Order

	tx.Scopes(pagination.Paginate(orders, &page, tx)).Find(&orders)
	page.Rows = orders

	return page, nil
}

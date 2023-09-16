package user

import (
	mu "book-nest/internal/models/user"
	"book-nest/utils/pagination"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
}

func NewUserRepository() mu.UserRepository {
	return &UserRepository{}
}

func (repo *UserRepository) Create(tx *gorm.DB, input *mu.User) (*mu.User, error) {
	resultcreate := tx.Create(&input)
	if resultcreate.Error != nil {
		return nil, resultcreate.Error
	}

	return input, nil
}

func (repo *UserRepository) CheckDuplicate(tx *gorm.DB, input *mu.User) error {
	var count int64
	if resultEmail := tx.Model(&mu.User{}).Where("email = $1 ", input.Email).Count(&count); resultEmail.Error != nil {
		return errors.New("error checking email")
	}
	if count > 0 {
		return errors.New("email already exists in database")
	}

	if resultPhone := tx.Model(&mu.User{}).Where("phone = $1 ", input.Phone).Count(&count); resultPhone.Error != nil {
		return errors.New("error checking phone")
	}
	if count > 0 {
		return errors.New("phone already exists in database")
	}

	return nil
}

func (repo *UserRepository) Update(tx *gorm.DB, input *mu.User) (*mu.User, error) {
	result := tx.Save(&input)
	if result.Error != nil {
		return nil, errors.New("error updating your data")
	}

	return input, nil
}

func (repo *UserRepository) GetByEmail(tx *gorm.DB, email string) (*mu.User, error) {
	user := new(mu.User)
	result := tx.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *UserRepository) GetDetail(tx *gorm.DB, userId uuid.UUID) (*mu.User, error) {
	user := new(mu.User)
	user.Id = userId
	result := tx.Omit("password").Omit("verification_code").First(&user, userId)
	if result.Error != nil {
		return nil, fmt.Errorf("user id %s not found", userId)
	}

	return user, nil
}

func (repo UserRepository) GetList(tx *gorm.DB, page pagination.Pagination) (pagination.Pagination, error) {
	var users []mu.User

	tx.Scopes(pagination.Paginate(users, &page, tx)).Find(&users)
	page.Rows = users

	return page, nil
}

func (repo UserRepository) Delete(tx *gorm.DB, userId uuid.UUID) error {
	var user mu.User
	result := tx.Delete(&user, userId)
	if result.Error != nil {
		return fmt.Errorf("user id %s not found", userId)
	}

	return nil
}

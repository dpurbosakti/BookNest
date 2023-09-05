package user

import (
	mu "book-nest/internal/models/user"
	"errors"

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

func (repo *UserRepository) CheckEmail(tx *gorm.DB, email string) (*mu.User, error) {
	user := new(mu.User)
	result := tx.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, errors.New("error checking email")
	}
	return user, nil
}

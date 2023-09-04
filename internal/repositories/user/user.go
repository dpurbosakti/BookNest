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
	// passwordHashed, errorHash := _bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	// if errorHash != nil {
	// 	fmt.Println("Error hash", errorHash.Error())
	// }
	// user.Password = string(passwordHashed)
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

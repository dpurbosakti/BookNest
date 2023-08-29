package user

import (
	"book-nest/internal/models/user"

	"gorm.io/gorm"
)

type UserRepository struct {
}

func NewUserRepository() user.UserRepository {
	return &UserRepository{}
}

func (repo *UserRepository) Create(tx *gorm.DB, input user.User) (user.User, error) {
	// passwordHashed, errorHash := _bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	// if errorHash != nil {
	// 	fmt.Println("Error hash", errorHash.Error())
	// }
	// user.Password = string(passwordHashed)
	resultcreate := tx.Create(&input)
	if resultcreate.Error != nil {
		return user.User{}, resultcreate.Error
	}

	return input, nil
}

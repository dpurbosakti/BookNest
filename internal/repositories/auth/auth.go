package user

import (
	"book-nest/internal/models/auth"
	"book-nest/internal/models/user"
	"fmt"

	"gorm.io/gorm"
)

type AuthRepository struct {
}

func NewAuthRepository() auth.AuthRepository {
	return &AuthRepository{}
}

func (repo *AuthRepository) Login(tx *gorm.DB, input auth.LoginRequest) (user.User, error) {
	var userData user.User
	result := tx.Where("email = ?", input.Email).First(&userData)
	if result.Error != nil {
		return user.User{}, fmt.Errorf("email %s not found", input.Email)
	}

	return userData, nil
}

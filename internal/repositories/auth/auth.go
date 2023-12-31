package user

import (
	i "book-nest/internal/interfaces"
	"book-nest/internal/models/auth"
	"book-nest/internal/models/user"
	"fmt"

	"gorm.io/gorm"
)

type AuthRepository struct {
}

func NewAuthRepository() i.AuthRepository {
	return &AuthRepository{}
}

func (repo *AuthRepository) Login(tx *gorm.DB, input auth.LoginRequest) (*user.User, error) {
	var userData user.User
	result := tx.Where("email = ?", input.Email).First(&userData)
	if result.Error != nil {
		return nil, fmt.Errorf("account with email: %s are not found", input.Email)
	}

	return &userData, nil
}

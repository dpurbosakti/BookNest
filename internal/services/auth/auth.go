package auth

import (
	ma "book-nest/internal/models/auth"

	"gorm.io/gorm"
)

type AuthService struct {
	AuthRepository ma.AuthRepository
	DB             *gorm.DB
}

func NewAuthService(authRepository ma.AuthRepository, db *gorm.DB) ma.AuthService {
	return &AuthService{
		AuthRepository: authRepository,
		DB:             db,
	}
}

func (srv *AuthService) Login(input ma.LoginRequest) {

}

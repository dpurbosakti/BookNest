package auth

import (
	ma "book-nest/internal/models/auth"
	jh "book-nest/utils/jwthelper"
	ph "book-nest/utils/passwordhelper"
	"errors"
	"fmt"

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

func (srv *AuthService) Login(input ma.LoginRequest) (*string, error) {
	resultRepo, err := srv.AuthRepository.Login(srv.DB, input)
	if err != nil {
		return nil, err
	}
	if resultRepo == nil {
		return nil, fmt.Errorf("account with email: %s are not found", input.Email)
	}

	if !resultRepo.IsVerified {
		return nil, errors.New("account is unverified")
	}
	errCrypt := ph.ComparePassword(resultRepo.Password, input.Password)
	if errCrypt != nil {
		return nil, errors.New("password incorrect")
	}

	token, err := jh.GenereteToken(resultRepo)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

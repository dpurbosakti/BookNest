package auth

import (
	i "book-nest/internal/interfaces"
	ma "book-nest/internal/models/auth"
	jh "book-nest/utils/jwthelper"
	ph "book-nest/utils/passwordhelper"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthService struct {
	AuthRepository i.AuthRepository
	UserRepository i.UserRepository
	DB             *gorm.DB
}

func NewAuthService(authRepository i.AuthRepository, userRepository i.UserRepository, db *gorm.DB) i.AuthService {
	return &AuthService{
		AuthRepository: authRepository,
		UserRepository: userRepository,
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

	token, err := jh.GenerateToken(resultRepo)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (srv *AuthService) LoginByGoogle(input *ma.GoogleResponse) (*string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "login_by_google",
		"scope": "auth service",
	})

	logger.WithField("input", input).Info()
	data := googleResponseToModel(input)
	logger.WithField("data", data).Info()
	var token string
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		resultGet, err := srv.UserRepository.GetByEmail(tx, input.Email)
		logger.WithError(err).Warn("user not found")

		switch err {
		case gorm.ErrRecordNotFound:
			resultCreate, err := srv.UserRepository.Create(tx, data)
			if err != nil {
				return err
			}
			token, err = jh.GenerateToken(resultCreate)
			if err != nil {
				return err
			}
		case nil:
			if !resultGet.IsVerified {
				return errors.New("account is unverified")
			}
			token, err = jh.GenerateToken(resultGet)
			if err != nil {
				return err
			}
			resultGet.OauthAccessToken = data.OauthAccessToken
			_, err := srv.UserRepository.Update(tx, resultGet)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		logger.Error("failed to login")
		return nil, err
	}

	return &token, nil
}

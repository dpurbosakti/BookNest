package user

import (
	mu "book-nest/internal/models/user"
	eh "book-nest/utils/emailhelper"
	ph "book-nest/utils/passwordhelper"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepository mu.UserRepository
	DB             *gorm.DB
	EmailHelper    eh.Emailer
}

func NewUserService(userRepository mu.UserRepository, db *gorm.DB, emailer eh.Emailer) mu.UserService {
	return &UserService{
		UserRepository: userRepository,
		DB:             db,
		EmailHelper:    emailer,
	}
}

func (srv *UserService) Create(input mu.UserCreateRequest) (mu.UserResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "create",
		"scope": "user service",
		"data":  input,
	})
	// var errChan = make(chan error)
	var result mu.UserResponse
	hashPassword, errHash := ph.HashPassword(input.Password)
	if errHash != nil {
		logger.WithError(errHash).Error("failed to hash password")
		return result, errHash
	}
	input.Password = hashPassword
	data := requestToModel(input)
	verificationCode, _ := generateVerCode(6)
	data.IsVerified = false
	data.VerificationCode = verificationCode
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.WithField("data", data).Info("db transaction begin")
		err := srv.UserRepository.CheckDuplicate(tx, data)
		if err != nil {
			logger.WithError(err).Error("failed to create user, there is duplicated data")
			return err
		}
		resultRepo, err := srv.UserRepository.Create(tx, data)
		if err != nil {
			logger.WithError(err).Error("failed to create user")
			return err
		}
		result = modelToResponse(resultRepo)

		// err = eh.SendEmailVerCode(data)
		// if err != nil {
		// 	return errors.New("failed to send email verification code: " + err.Error())
		// }
		if err := srv.EmailHelper.SendEmailVerificationCode(data); err != nil {
			return errors.New("failed to send email: " + err.Error())
		}
		logger.WithField("data", data).Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.WithError(err).Error("failed to create user")
		return mu.UserResponse{}, err
	}

	return result, nil
}

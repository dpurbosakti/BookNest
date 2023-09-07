package user

import (
	mu "book-nest/internal/models/user"
	eh "book-nest/utils/emailhelper"
	"book-nest/utils/pagination"
	ph "book-nest/utils/passwordhelper"
	"errors"

	"github.com/google/uuid"
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

func (srv *UserService) Create(input *mu.UserCreateRequest) (*mu.UserResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "create",
		"scope": "user service",
	})

	if input == nil {
		return nil, errors.New("input is nil")
	}

	logger.WithField("data", input)
	result := new(mu.UserResponse)
	hashPassword, errHash := ph.HashPassword(input.Password)
	if errHash != nil {
		logger.WithError(errHash).Error("failed to hash password")
		return nil, errHash
	}
	input.Password = hashPassword
	data := requestToModel(input)
	verificationCode, _ := generateVerificationCode(6)
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
		if err := srv.EmailHelper.SendEmailVerificationCode(data); err != nil {
			return errors.New("failed to send email: " + err.Error())
		}
		logger.WithField("data", data).Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.WithError(err).Error("failed to create user")
		return nil, err
	}

	return result, nil
}

func (srv *UserService) Verify(input *mu.UserVerifyRequest) error {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "verify",
		"scope": "user service",
	})
	logger.WithField("data", input)
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		result, err := srv.UserRepository.CheckEmail(tx, input.Email)
		if err != nil {
			logger.WithError(err).Error("failed to check email")
			return err
		}
		if result.IsVerified {
			logger.Error("account is already verified")
			return errors.New("account is already verified")
		}
		if result.VerificationCode != input.VerificationCode {
			logger.Error("the verification code you entered is incorrect")
			return errors.New("the verification code you entered is incorrect")
		}

		result.IsVerified = true
		_, errSave := srv.UserRepository.Update(tx, result)
		if errSave != nil {
			logger.WithError(errSave).Error("failed to update data")
			return errSave
		}

		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.WithError(err).Error("failed to verify account")
		return err
	}

	return nil
}

func (srv *UserService) RefreshVerificationCode(input *mu.UserVerificationCodeRequest) error {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "refresh_verification_code",
		"scope": "user service",
	})
	dataUser := new(mu.User)
	logger.WithField("data", input)
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultRepo, err := srv.UserRepository.CheckEmail(tx, input.Email)
		logger.WithError(err).Error("failed to check email")
		if err != nil {
			return err
		}
		if resultRepo.IsVerified {
			logger.Error("account is verified no need to refresh verification code")
			return errors.New("account is verified no need to refresh verification code")
		}
		verCode, _ := generateVerificationCode(6)
		resultRepo.VerificationCode = verCode
		_, err = srv.UserRepository.Update(tx, resultRepo)
		if err != nil {
			logger.WithError(err).Error("failed to update data")
			return err
		}
		dataUser = resultRepo
		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.WithError(err).Error("failed to refresh verification code")
		return err
	}

	err = srv.EmailHelper.SendEmailVerificationCode(dataUser)
	if err != nil {
		logger.Error("failed to send email verification code")
		return errors.New("failed to send email verification code: " + err.Error())
	}
	return nil
}

func (srv *UserService) GetDetail(userId uuid.UUID) (*mu.UserResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_detail",
		"scope": "user service",
	})
	logger.WithField("id", userId)
	result := new(mu.UserResponse)
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultRepo, err := srv.UserRepository.GetDetail(tx, userId)
		logger.WithError(err).Error("failed to get detail")
		if err != nil {
			return err
		}
		result = modelToResponse(resultRepo)
		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.Error("failed to get detail")
		return result, err
	}

	return result, nil
}

func (srv *UserService) GetList(page pagination.Pagination) (pagination.Pagination, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_list",
		"scope": "user service",
	})
	var result pagination.Pagination
	logger.Info("data page", page)
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultRepo, err := srv.UserRepository.GetList(tx, page)
		if err != nil {
			logger.WithError(err).Error("failed to get list")
			return err
		}
		result = resultRepo
		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.Error("failed to get list")
		return result, err
	}

	return result, nil
}

func (srv *UserService) Delete(userId uuid.UUID) error {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "delete",
		"scope": "user service",
	})
	logger.WithField("id", userId)
	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		err := srv.UserRepository.Delete(tx, userId)
		if err != nil {
			logger.WithError(err).Error("failed to delete data")
			return err
		}
		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.Error("failed to delete data")
		return err
	}

	return nil
}

func (srv *UserService) Update(input *mu.UserUpdateRequest, userId uuid.UUID) (*mu.UserResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "update",
		"scope": "user service",
	})
	if input.Password != nil {
		hashPassword, errHash := ph.HashPassword(*input.Password)
		if errHash != nil {
			return nil, errHash
		}
		input.Password = &hashPassword
	}
	logger.WithField("data", input)

	result := new(mu.UserResponse)

	err := srv.DB.Transaction(func(tx *gorm.DB) error {
		logger.Info("db transaction begin")
		resultGet, err := srv.UserRepository.GetDetail(tx, userId)
		if err != nil {
			logger.WithError(err).Error("failed to get detail")
			return err
		}
		copier(resultGet, input)
		resultUpdate, err := srv.UserRepository.Update(tx, resultGet)
		if err != nil {
			logger.WithError(err).Error("failed to update data")
			return err
		}
		result = modelToResponse(resultUpdate)
		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		logger.Error("failed to update data")
		return result, err
	}

	return result, nil
}

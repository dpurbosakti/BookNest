package user

import (
	"book-nest/clients/gomail"
	i "book-nest/internal/interfaces"
	mu "book-nest/internal/models/user"
	eh "book-nest/utils/errorhelper"
	"book-nest/utils/pagination"
	ph "book-nest/utils/passwordhelper"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepository i.UserRepository
	DB             *gorm.DB
	Gomail         *gomail.Gomail
}

func NewUserService(userRepository i.UserRepository, db *gorm.DB, gomail *gomail.Gomail) i.UserService {
	return &UserService{
		UserRepository: userRepository,
		DB:             db,
		Gomail:         gomail,
	}
}

const userScope = "user"

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
			eh.FailedCreate(logger, err, userScope)
			return err
		}
		result = modelToResponse(resultRepo)
		if err := srv.Gomail.SendEmailVerificationCode(data); err != nil {
			logger.WithError(err).Error("failed to send email")
			return errors.New("failed to send email: " + err.Error())
		}
		logger.WithField("data", data).Info("end of db transaction")
		return nil
	})
	if err != nil {
		eh.FailedCreate(logger, err, userScope)
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
		result, err := srv.UserRepository.GetByEmail(tx, input.Email)
		if err != nil {
			eh.FailedGetDetail(logger, err, userScope)
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
			eh.FailedUpdate(logger, err, userScope)
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
		resultRepo, err := srv.UserRepository.GetByEmail(tx, input.Email)
		eh.FailedGetDetail(logger, err, userScope)
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
			eh.FailedUpdate(logger, err, userScope)
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

	err = srv.Gomail.SendEmailVerificationCode(dataUser)
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
		eh.FailedGetDetail(logger, err, userScope)
		if err != nil {
			return err
		}
		result = modelToResponse(resultRepo)
		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		eh.FailedGetDetail(logger, err, userScope)
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
			eh.FailedGetList(logger, err, userScope)
			return err
		}
		result = resultRepo
		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		eh.FailedGetList(logger, err, userScope)
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
			eh.FailedDelete(logger, err, userScope)
			return err
		}
		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		eh.FailedDelete(logger, err, userScope)
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
			eh.FailedGetDetail(logger, err, userScope)
			return err
		}
		resultGet.Copier(input)
		resultUpdate, err := srv.UserRepository.Update(tx, resultGet)
		if err != nil {
			eh.FailedUpdate(logger, err, userScope)
			return err
		}
		result = modelToResponse(resultUpdate)
		logger.Info("end of db transaction")
		return nil
	})
	if err != nil {
		eh.FailedUpdate(logger, err, userScope)
		return result, err
	}

	return result, nil
}

package user

import (
	mu "book-nest/internal/models/user"

	"gorm.io/gorm"
)

type UserService struct {
	UserRepository mu.UserRepository
	DB             *gorm.DB
}

func NewUserService(userRepository mu.UserRepository, db *gorm.DB) mu.UserService {
	return &UserService{
		UserRepository: userRepository,
		DB:             db,
	}
}

func (srv *UserService) Create(input mu.UserCreateRequest) (result mu.UserResponse, err error) {
	// var errChan = make(chan error)
	// hashPassword, errHash := ph.HashPassword(input.Password)
	// if errHash != nil {
	// 	return result, errHash
	// }
	// input.Password = hashPassword
	// data := createRequestToModel(input)
	// verCode, _ := generateVerCode(6)
	// data.IsVerified = false
	// data.VerCode = verCode
	err = srv.DB.Transaction(func(tx *gorm.DB) error {
		// err := srv.UserRepository.CheckDuplicate(tx, data)
		// if err != nil {
		// 	return err
		// }
		// resultRepo, err := srv.UserRepository.Create(tx, data)
		// if err != nil {
		// 	return err
		// }
		// result = modelToResponse(resultRepo)

		// err = eh.SendEmailVerCode(data)
		// if err != nil {
		// 	return errors.New("failed to send email verification code: " + err.Error())
		// }
		return nil
	})
	if err != nil {
		// err = sh.SetError(scope, "create", "error create new data", err.Error())
		return mu.UserResponse{}, err
	}

	return result, nil
}

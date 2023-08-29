package user

import "gorm.io/gorm"

type UserRepository interface {
	Create(tx *gorm.DB, input User) (User, error)
}

type UserService interface {
	Create(input UserCreateRequest) (UserResponse, error)
}

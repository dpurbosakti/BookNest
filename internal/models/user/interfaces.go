package user

import "gorm.io/gorm"

type UserRepository interface {
	Create(tx *gorm.DB, input User) (row int, err error)
}

type UserService interface {
	Create(input UserCreateRequest) (result UserResponse, err error)
}

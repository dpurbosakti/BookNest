package user

import "gorm.io/gorm"

type UserRepository interface {
	Create(tx *gorm.DB, input *User) (*User, error)
	CheckDuplicate(tx *gorm.DB, input *User) error
}

type UserService interface {
	Create(input *UserCreateRequest) (*UserResponse, error)
}

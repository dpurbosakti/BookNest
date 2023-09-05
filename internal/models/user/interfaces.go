package user

import "gorm.io/gorm"

type UserRepository interface {
	Create(tx *gorm.DB, input *User) (*User, error)
	CheckDuplicate(tx *gorm.DB, input *User) error
	CheckEmail(tx *gorm.DB, email string) (*User, error)
	Update(tx *gorm.DB, input *User) (*User, error)
}

type UserService interface {
	Create(input *UserCreateRequest) (*UserResponse, error)
	Verify(input *UserVerifyRequest) (err error)
	RefreshVerificationCode(input *UserVerificationCodeRequest) (err error)
}

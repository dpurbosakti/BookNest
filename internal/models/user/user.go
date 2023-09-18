package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id               uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name             string         `json:"name" gorm:"type:varchar(250);not null"`
	Email            string         `json:"email" gorm:"type:varchar(250);not null;unique"`
	Password         string         `json:"password" gorm:"type:varchar(250);not null"`
	Phone            string         `json:"phone" gorm:"type:varchar(20);not null;unique"`
	Address          string         `json:"address" gorm:"type:varchar(250);not null"`
	Role             string         `json:"role" gorm:"default:user"`
	VerificationCode string         `json:"verification_code"`
	IsVerified       bool           `json:"is_verified"`
	OauthAccessToken string         `json:"oauth_access_token"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (u *User) Copier(input *UserUpdateRequest) {
	if input.Name != nil {
		u.Name = *input.Name
	}

	if input.Address != nil {
		u.Address = *input.Address
	}

	if input.Password != nil {
		u.Password = *input.Password
	}

	if input.Phone != nil {
		u.Phone = *input.Phone
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if len(u.Name) < 2 {
		return errors.New("name must be at least 3 characters")
	}
	return nil
}

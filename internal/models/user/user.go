package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id               uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name             string         `json:"name" gorm:"type:varchar(250);not null"`
	Email            string         `json:"email" gorm:"type:varchar(250);not null"`
	Password         string         `json:"password" gorm:"type:varchar(250);not null"`
	Phone            string         `json:"phone" gorm:"type:varchar(20);not null"`
	Address          string         `json:"address" gorm:"type:varchar(250);not null"`
	Role             string         `json:"role" gorm:"default:user"`
	VerificationCode string         `json:"verificationCode"` //verification code
	IsVerified       bool           `json:"isVerified"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

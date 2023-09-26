package address

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Address struct {
	Id         uint64         `json:"id"`
	UserId     uuid.UUID      `json:"user_id"`
	Address    string         `json:"address"`
	Notes      string         `json:"notes"`
	Latitude   float64        `json:"latitude"`
	Longitude  float64        `json:"longitude"`
	PostalCode string         `json:"postal_code"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}

func (a *Address) Copier(input *AddressUpdateRequest) {
	if input.Address != nil {
		a.Address = *input.Address
	}

	if input.Notes != nil {
		a.Notes = *input.Notes
	}

	if input.Latitude != nil {
		a.Latitude = *input.Latitude
	}

	if input.Longitude != nil {
		a.Longitude = *input.Longitude
	}

	if input.PostalCode != nil {
		a.PostalCode = *input.PostalCode
	}
}

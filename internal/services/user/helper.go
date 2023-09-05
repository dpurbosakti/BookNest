package user

import (
	"book-nest/internal/models/user"
	"crypto/rand"
	"math/big"
)

// mappers
func requestToModel(input *user.UserCreateRequest) *user.User {
	return &user.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Phone:    input.Phone,
		Address:  input.Address,
	}
}

func modelToResponse(input *user.User) *user.UserResponse {
	return &user.UserResponse{
		Id:               input.Id,
		Name:             input.Name,
		Email:            input.Email,
		Phone:            input.Phone,
		Address:          input.Address,
		Role:             input.Role,
		VerificationCode: input.VerificationCode,
		IsVerified:       input.IsVerified,
		CreatedAt:        input.CreatedAt,
		UpdatedAt:        input.UpdatedAt,
		DeletedAt:        input.DeletedAt,
	}
}

// verification code generator
func generateVerificationCode(length int) (string, error) {
	seed := "012345679"
	byteSlice := make([]byte, length)

	for i := 0; i < length; i++ {
		max := big.NewInt(int64(len(seed)))
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}

		byteSlice[i] = seed[num.Int64()]
	}

	return string(byteSlice), nil
}

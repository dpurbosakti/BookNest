package auth

import (
	ma "book-nest/internal/models/auth"
	mu "book-nest/internal/models/user"
)

func googleResponseToModel(input *ma.GoogleResponse) *mu.User {
	return &mu.User{
		Email:      input.Email,
		Name:       input.Name,
		IsVerified: input.VerifiedEmail,
	}
}

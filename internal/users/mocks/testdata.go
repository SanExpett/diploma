package mocks

import (
	"time"

	"github.com/SanExpett/diploma/internal/domain"
)

func NewMockUser() domain.User {
	return domain.User{
		Uuid:         "1",
		Email:        "cakethefake@gmail.com",
		Avatar:       "",
		Name:         "Danya",
		Password:     "123456789",
		IsAdmin:      true,
		RegisteredAt: time.Now(),
		Birthday:     time.Now(),
	}
}

func NewMockUserSignUp() domain.UserSignUp {
	return domain.UserSignUp{
		Email:    "cakethefake@gmail.com",
		Name:     "Danya",
		Password: "123456789",
	}
}

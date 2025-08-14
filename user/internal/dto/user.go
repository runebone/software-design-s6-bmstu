package dto

import (
	"time"
	"user/internal/entity"

	"github.com/google/uuid"
)

type UserDTO struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateUserDTO struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"` // NOTE: unencrypted
}

type UpdateUserDTO struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
}

func ToUserDTO(user entity.User) UserDTO {
	return UserDTO{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Role:         user.Role,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
	}
}

func ToUserDTOs(users []entity.User) []UserDTO {
	userDTOs := make([]UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = ToUserDTO(user)
	}
	return userDTOs
}

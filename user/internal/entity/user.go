package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	Role         string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

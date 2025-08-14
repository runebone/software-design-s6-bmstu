package entity

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	CreatedAt time.Time
}

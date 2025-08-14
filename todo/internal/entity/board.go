package entity

import (
	"time"

	"github.com/google/uuid"
)

type Board struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

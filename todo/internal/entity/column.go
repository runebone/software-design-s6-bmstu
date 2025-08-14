package entity

import (
	"time"

	"github.com/google/uuid"
)

type Column struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	BoardID   uuid.UUID
	Title     string
	Position  float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

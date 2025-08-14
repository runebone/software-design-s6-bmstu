package entity

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	ColumnID    uuid.UUID
	Title       string
	Description string
	Position    float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

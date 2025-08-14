package dto

import (
	"time"
	"todo/internal/entity"

	"github.com/google/uuid"
)

type CreateCardRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	ColumnID    uuid.UUID `json:"column_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Position    float64   `json:"position"`
}

type Card struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	ColumnID    uuid.UUID `json:"column_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Position    float64   `json:"position"`
	CreatedAt   time.Time `json:"created_at"`
}

type UpdateCardRequest struct {
	ID          uuid.UUID `json:"id"`
	ColumnID    uuid.UUID `json:"column_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Position    float64   `json:"position,omitempty"`
}

func ToCardDTO(card *entity.Card) Card {
	return Card{
		ID:          card.ID,
		UserID:      card.UserID,
		ColumnID:    card.ColumnID,
		Title:       card.Title,
		Description: card.Description,
		Position:    card.Position,
		CreatedAt:   card.CreatedAt,
	}
}

func ToCardDTOs(cards []entity.Card) []Card {
	cardDTOs := make([]Card, len(cards))
	for i, card := range cards {
		cardDTOs[i] = ToCardDTO(&card)
	}
	return cardDTOs
}

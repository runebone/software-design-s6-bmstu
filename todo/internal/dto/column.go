package dto

import (
	"todo/internal/entity"

	"github.com/google/uuid"
)

type CreateColumnRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	BoardID  uuid.UUID `json:"board_id"`
	Title    string    `json:"title"`
	Position float64   `json:"position"`
}

type Column struct {
	ID       uuid.UUID `json:"id"`
	BoardID  uuid.UUID `json:"board_id"`
	Title    string    `json:"title"`
	Position float64   `json:"position"`
}

type UpdateColumnRequest struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title,omitempty"`
	Position float64   `json:"position,omitempty"`
}

func ToColumnDTO(column *entity.Column) Column {
	return Column{
		ID:       column.ID,
		BoardID:  column.BoardID,
		Title:    column.Title,
		Position: column.Position,
	}
}

func ToColumnDTOs(columns []entity.Column) []Column {
	columnDTOs := make([]Column, len(columns))
	for i, column := range columns {
		columnDTOs[i] = ToColumnDTO(&column)
	}
	return columnDTOs
}

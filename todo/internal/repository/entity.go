package repository

import (
	"time"
	"todo/internal/entity"

	"github.com/google/uuid"
)

type Board struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Column struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	BoardID   uuid.UUID `db:"board_id"`
	Title     string    `db:"title"`
	Position  float64   `db:"position"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Card struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	ColumnID    uuid.UUID `db:"column_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Position    float64   `db:"position"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func RepoBoard(e entity.Board) Board {
	return Board{
		ID:        e.ID,
		UserID:    e.UserID,
		Title:     e.Title,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func RepoColumn(e entity.Column) Column {
	return Column{
		ID:        e.ID,
		UserID:    e.UserID,
		BoardID:   e.BoardID,
		Title:     e.Title,
		Position:  e.Position,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func RepoCard(e entity.Card) Card {
	return Card{
		ID:          e.ID,
		UserID:      e.UserID,
		ColumnID:    e.ColumnID,
		Title:       e.Title,
		Description: e.Description,
		Position:    e.Position,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func BoardToEntity(r Board) entity.Board {
	return entity.Board{
		ID:        r.ID,
		UserID:    r.UserID,
		Title:     r.Title,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func ColumnToEntity(r Column) entity.Column {
	return entity.Column{
		ID:        r.ID,
		UserID:    r.UserID,
		BoardID:   r.BoardID,
		Title:     r.Title,
		Position:  r.Position,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func CardToEntity(r Card) entity.Card {
	return entity.Card{
		ID:          r.ID,
		UserID:      r.UserID,
		ColumnID:    r.ColumnID,
		Title:       r.Title,
		Description: r.Description,
		Position:    r.Position,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

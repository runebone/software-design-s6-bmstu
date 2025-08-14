package repository

import (
	"context"
	"time"
	"todo/internal/entity"

	"github.com/google/uuid"
)

type BoardRepository interface {
	CreateBoard(ctx context.Context, board *entity.Board) error
	GetBoardByID(ctx context.Context, id uuid.UUID) (*entity.Board, error)
	GetBoardsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.Board, error)
	UpdateBoard(ctx context.Context, board *entity.Board) error
	DeleteBoard(ctx context.Context, id uuid.UUID) error
}

type ColumnRepository interface {
	CreateColumn(ctx context.Context, column *entity.Column) error
	GetColumnByID(ctx context.Context, id uuid.UUID) (*entity.Column, error)
	GetColumnsByBoard(ctx context.Context, boardID uuid.UUID, limit, offset int) ([]entity.Column, error)
	UpdateColumn(ctx context.Context, column *entity.Column) error
	DeleteColumn(ctx context.Context, id uuid.UUID) error
}

type CardRepository interface {
	CreateCard(ctx context.Context, card *entity.Card) error
	GetCardByID(ctx context.Context, id uuid.UUID) (*entity.Card, error)
	GetCardsByColumn(ctx context.Context, columnID uuid.UUID, limit, offset int) ([]entity.Card, error)
	GetNewCards(ctx context.Context, from, to time.Time) ([]entity.Card, error)
	UpdateCard(ctx context.Context, card *entity.Card) error
	MoveCard(ctx context.Context, card *entity.Card) error
	DeleteCard(ctx context.Context, id uuid.UUID) error
}

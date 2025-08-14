package todo

import (
	"aggregator/internal/dto"
	"context"
	"time"
)

type TodoService interface {
	GetNewCards(ctx context.Context, from, to time.Time) ([]dto.Card, error)

	GetBoards(ctx context.Context, userID string) ([]dto.Board, error)
	GetColumns(ctx context.Context, boardID string) ([]dto.Column, error)
	GetCards(ctx context.Context, columnID string) ([]dto.Card, error)
	GetCard(ctx context.Context, id string) (*dto.Card, error)

	CreateBoard(ctx context.Context, board dto.Board) error
	CreateColumn(ctx context.Context, column dto.Column) error
	CreateCard(ctx context.Context, card dto.Card) error

	UpdateBoard(ctx context.Context, board *dto.Board) error
	UpdateColumn(ctx context.Context, column *dto.Column) error
	UpdateCard(ctx context.Context, card *dto.Card) error

	DeleteBoard(ctx context.Context, id string) error
	DeleteColumn(ctx context.Context, id string) error
	DeleteCard(ctx context.Context, id string) error
}

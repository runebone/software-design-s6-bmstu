package service

import (
	"cli/internal/dto"
	"context"
)

type AggregatorService interface {
	Register(ctx context.Context, username, email, password string) (*dto.Tokens, error)
	Login(ctx context.Context, email, password string) (*dto.Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
	Validate(ctx context.Context, token string) (*dto.ValidateTokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error

	ShowBoards(ctx context.Context) ([]dto.Board, error)
	ShowBoard(ctx context.Context, boardID string) ([]dto.Column, error)
	ShowColumn(ctx context.Context, columnID string) ([]dto.Card, error)
	ShowCard(ctx context.Context, cardID string) (*dto.Card, error)

	CreateBoard(ctx context.Context, board dto.Board) error
	CreateColumn(ctx context.Context, column dto.Column) error
	CreateCard(ctx context.Context, card dto.Card) error

	UpdateBoard(ctx context.Context, board *dto.Board) error
	UpdateColumn(ctx context.Context, column *dto.Column) error
	UpdateCard(ctx context.Context, card *dto.Card) error

	DeleteBoard(ctx context.Context, id string) error
	DeleteColumn(ctx context.Context, id string) error
	DeleteCard(ctx context.Context, id string) error

	Stats(ctx context.Context, from, to string) ([]dto.NewUsersAndCardsStats, error)
}

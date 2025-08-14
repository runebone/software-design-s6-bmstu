package usecase

import (
	"cli/internal/dto"
	"context"
)

type Client interface {
	Register(ctx context.Context, username, email, password string) (*dto.Tokens, error)
	Login(ctx context.Context, email, password string) (*dto.Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
	Validate(ctx context.Context, token string) (*dto.ValidateTokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error

	// context with value tokens
	ShowBoards(ctx context.Context)
	ShowBoard(ctx context.Context, boardID string)
	ShowColumn(ctx context.Context, columnID string)
	ShowCard(ctx context.Context, cardID string)

	CreateBoard(ctx context.Context, title string)
	CreateColumn(ctx context.Context, boardID, title string)
	CreateCard(ctx context.Context, columnID, title, description string)

	UpdateBoard(ctx context.Context, boardID, title string)
	UpdateColumn(ctx context.Context, columnID, title string)
	UpdateCardTitle(ctx context.Context, cardID, title string)
	UpdateCardDescription(ctx context.Context, cardID, description string)
	MoveCard(ctx context.Context, cardIDstr, columnIDstr string)

	DeleteBoard(ctx context.Context, id string)
	DeleteColumn(ctx context.Context, id string)
	DeleteCard(ctx context.Context, id string)

	Stats(ctx context.Context, from, to string)
}

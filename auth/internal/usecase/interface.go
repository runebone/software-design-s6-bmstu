package usecase

import (
	"auth/internal/dto"
	"context"
)

type AuthUsecase interface {
	Register(ctx context.Context, username, email, password string) (*dto.Tokens, error)
	Login(ctx context.Context, email, password string) (*dto.Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponse, error)
	ValidateToken(ctx context.Context, token string) (string, string, error)
	Logout(ctx context.Context, refreshToken string) error
}

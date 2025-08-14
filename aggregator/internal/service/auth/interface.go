package auth

import (
	"aggregator/internal/dto"
	"context"
)

type AuthService interface {
	Register(ctx context.Context, username, email, password string) (*dto.Tokens, error)
	Login(ctx context.Context, email, password string) (*dto.Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
	ValidateToken(ctx context.Context, token string) (*dto.ValidateTokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

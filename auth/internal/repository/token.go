package repository

import (
	"auth/internal/entity"
	"context"
)

type TokenRepository interface {
	Save(ctx context.Context, token *entity.Token) error
	Delete(ctx context.Context, tokenID string) error
	FindByToken(ctx context.Context, token string) (*entity.Token, error)
}

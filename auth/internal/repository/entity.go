package repository

import (
	"auth/internal/entity"
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Token     string    `db:"token"`
	CreatedAt time.Time `db:"created_at"`
}

func RepoToken(e entity.Token) Token {
	return Token{
		ID:        e.ID,
		UserID:    e.UserID,
		Token:     e.Token,
		CreatedAt: e.CreatedAt,
	}
}

func TokenToEntity(r Token) entity.Token {
	return entity.Token{
		ID:        r.ID,
		UserID:    r.UserID,
		Token:     r.Token,
		CreatedAt: r.CreatedAt,
	}
}

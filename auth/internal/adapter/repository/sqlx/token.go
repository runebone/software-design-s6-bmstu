package sqlx

import (
	"auth/internal/entity"
	"auth/internal/repository"
	"context"

	"github.com/jmoiron/sqlx"
)

type SQLXTokenRepository struct {
	db *sqlx.DB
}

func NewSQLXTokenRepository(db *sqlx.DB) *SQLXTokenRepository {
	return &SQLXTokenRepository{
		db: db,
	}
}

func (r *SQLXTokenRepository) Save(ctx context.Context, token *entity.Token) error {
	repoToken := repository.RepoToken(*token)

	query := `
        INSERT INTO tokens (id, user_id, token, created_at)
		VALUES (:id, :user_id, :token, :created_at)
    `

	_, err := r.db.NamedExecContext(ctx, query, repoToken)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLXTokenRepository) Delete(ctx context.Context, tokenID string) error {
	query := `DELETE FROM tokens WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, tokenID)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLXTokenRepository) FindByToken(ctx context.Context, tokenValue string) (*entity.Token, error) {
	query := `SELECT id, user_id, token, created_at FROM tokens WHERE token = $1`

	var repoToken repository.Token

	err := r.db.GetContext(ctx, &repoToken, query, tokenValue)
	if err != nil {
		return nil, err
	}

	token := repository.TokenToEntity(repoToken)

	return &token, nil
}

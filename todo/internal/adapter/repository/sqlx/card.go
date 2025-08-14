package repository

import (
	"context"
	"time"
	"todo/internal/entity"
	"todo/internal/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SQLXCardRepository struct {
	db *sqlx.DB
}

func NewSQLXCardRepository(db *sqlx.DB) *SQLXCardRepository {
	return &SQLXCardRepository{db: db}
}

func (r *SQLXCardRepository) CreateCard(ctx context.Context, card *entity.Card) error {
	query := `
	INSERT INTO cards (id, column_id, user_id, title, description, position, created_at, updated_at)
	VALUES (:id, :column_id, :user_id, :title, :description, :position, :created_at, :updated_at)
	`

	repoCard := repository.RepoCard(*card)

	_, err := r.db.NamedExecContext(ctx, query, repoCard)

	return err
}

func (r *SQLXCardRepository) GetCardByID(ctx context.Context, id uuid.UUID) (*entity.Card, error) {
	query := `
	SELECT * FROM cards WHERE id = $1
	`

	var repoCard repository.Card
	err := r.db.GetContext(ctx, &repoCard, query, id)

	if err != nil {
		return nil, err
	}

	card := repository.CardToEntity(repoCard)

	return &card, nil
}

func (r *SQLXCardRepository) GetCardsByColumn(ctx context.Context, columnID uuid.UUID, limit, offset int) ([]entity.Card, error) {
	query := `
	SELECT * FROM cards WHERE column_id = $1
	ORDER BY created_at ASC
	LIMIT $2
	OFFSET $3
	`

	var repoCards []repository.Card
	err := r.db.SelectContext(ctx, &repoCards, query, columnID, limit, offset)

	if err != nil {
		return nil, err
	}

	cards := make([]entity.Card, len(repoCards))
	for i, c := range repoCards {
		cards[i] = repository.CardToEntity(c)
	}

	return cards, nil
}

func (r *SQLXCardRepository) UpdateCard(ctx context.Context, card *entity.Card) error {
	query := `
    UPDATE cards SET
	title = :title,
	description = :description,
	position = :position,
	updated_at = :updated_at
    WHERE id = :id
    `

	repoCard := repository.RepoCard(*card)

	_, err := r.db.NamedExecContext(ctx, query, repoCard)

	return err
}

func (r *SQLXCardRepository) MoveCard(ctx context.Context, card *entity.Card) error {
	query := `
    UPDATE cards SET
	column_id = :column_id,
	updated_at = :updated_at
    WHERE id = :id
    `

	repoCard := repository.RepoCard(*card)

	_, err := r.db.NamedExecContext(ctx, query, repoCard)

	return err
}

func (r *SQLXCardRepository) DeleteCard(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM cards WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}

func (r *SQLXCardRepository) GetNewCards(ctx context.Context, from, to time.Time) ([]entity.Card, error) {
	query := `
	SELECT * FROM cards
	WHERE $1 <= created_at AND created_at <= $2
	`

	var repoCards []repository.Card
	err := r.db.SelectContext(ctx, &repoCards, query, from, to)

	if err != nil {
		return nil, err
	}

	cards := make([]entity.Card, len(repoCards))
	for i, c := range repoCards {
		cards[i] = repository.CardToEntity(c)
	}

	return cards, nil
}

package repository

import (
	"context"
	"todo/internal/entity"
	"todo/internal/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SQLXBoardRepository struct {
	db *sqlx.DB
}

func NewSQLXBoardRepository(db *sqlx.DB) *SQLXBoardRepository {
	return &SQLXBoardRepository{db: db}
}

func (r *SQLXBoardRepository) CreateBoard(ctx context.Context, board *entity.Board) error {
	repoBoard := repository.RepoBoard(*board)

	query := `
    INSERT INTO boards (id, user_id, title, created_at, updated_at)
	VALUES (:id, :user_id, :title, :created_at, :updated_at)
    `

	_, err := r.db.NamedExecContext(ctx, query, repoBoard)

	return err
}

func (r *SQLXBoardRepository) GetBoardByID(ctx context.Context, id uuid.UUID) (*entity.Board, error) {
	query := `
	SELECT * FROM boards WHERE id = $1
	`

	var repoBoard repository.Board
	err := r.db.GetContext(ctx, &repoBoard, query, id)

	if err != nil {
		return nil, err
	}

	board := repository.BoardToEntity(repoBoard)

	return &board, err
}

func (r *SQLXBoardRepository) GetBoardsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.Board, error) {
	query := `
	SELECT * FROM boards WHERE user_id = $1
	ORDER BY created_at ASC
	LIMIT $2
	OFFSET $3
	`

	var repoBoards []repository.Board
	err := r.db.SelectContext(ctx, &repoBoards, query, userID, limit, offset)

	if err != nil {
		return nil, err
	}

	boards := make([]entity.Board, len(repoBoards))
	for i, b := range repoBoards {
		boards[i] = repository.BoardToEntity(b)
	}

	return boards, nil
}

func (r *SQLXBoardRepository) UpdateBoard(ctx context.Context, board *entity.Board) error {
	repoBoard := repository.RepoBoard(*board)

	query := `
    UPDATE boards SET
	title = :title,
	updated_at = :updated_at
    WHERE id = :id
    `

	_, err := r.db.NamedExecContext(ctx, query, repoBoard)

	return err
}

func (r *SQLXBoardRepository) DeleteBoard(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM boards WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}

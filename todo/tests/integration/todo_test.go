package integration_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	sqlxRepository "todo/internal/adapter/repository/sqlx"
	"todo/internal/entity"
	"todo/internal/repository"
	"todo/internal/usecase"
	v1 "todo/internal/usecase/v1"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"
)

var db *sqlx.DB

type testSetup struct {
	ctx        context.Context
	boardRepo  repository.BoardRepository
	columnRepo repository.ColumnRepository
	cardRepo   repository.CardRepository
	uc         usecase.TodoUseCase
}

func sqlxSetup() *testSetup {
	ctx := context.TODO()
	boardRepo := sqlxRepository.NewSQLXBoardRepository(db)
	columnRepo := sqlxRepository.NewSQLXColumnRepository(db)
	cardRepo := sqlxRepository.NewSQLXCardRepository(db)
	uc := v1.NewTodoUseCase(boardRepo, columnRepo, cardRepo)

	return &testSetup{
		ctx:        ctx,
		boardRepo:  boardRepo,
		columnRepo: columnRepo,
		cardRepo:   cardRepo,
		uc:         uc,
	}
}

func applyMigrations(dsn string) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("Failed to create migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///home/rukost/University/software-design-s6-bmstu.git/lab4/src/todo/migrations/sql",
		"nigger",
		driver,
	)
	if err != nil {
		return fmt.Errorf("Failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Failed to apply migrations: %w", err)
	}

	return nil
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
			"TZ":                "Europe/Moscow",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	dbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: dbReq,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start PostgreSQL container: %v", err)
	}
	defer dbContainer.Terminate(ctx)

	host, err := dbContainer.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to get container host: %v", err)
	}

	port, err := dbContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("Failed to get container port: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=testuser password=testpass dbname=testdb sslmode=disable", host, port.Port())
	db, err = sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}

	if err := applyMigrations(dsn); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	code := m.Run()

	db.Close()
	os.Exit(code)
}

// CreateBoard(ctx context.Context, board *entity.Board) error
func resetDatabase() error {
	_, err := db.Exec(`
	TRUNCATE TABLE boards RESTART IDENTITY CASCADE
	`)
	if err == nil {
		_, err = db.Exec(`
		TRUNCATE TABLE columns RESTART IDENTITY CASCADE
		`)
	}
	if err == nil {
		_, err = db.Exec(`
		TRUNCATE TABLE cards RESTART IDENTITY CASCADE
		`)
	}
	return err
}

func TestCreate(t *testing.T) {
	ts := sqlxSetup()
	resetDatabase()

	userID := uuid.New()

	board := entity.Board{
		UserID: userID,
		Title:  "Board Title",
	}

	err := ts.uc.CreateBoard(ts.ctx, &board)

	if err != nil {
		log.Fatalf("Failed to execute CreateBoard usecase: %v", err)
	}

	var createdBoard repository.Board
	err = db.GetContext(ts.ctx, &createdBoard, `
		SELECT * FROM boards WHERE user_id = $1
	`, userID)

	if err != nil {
		log.Fatalf("Failed to select created board: %v", err)
	}

	assert.Equal(t, board.UserID, createdBoard.UserID)
	assert.Equal(t, board.Title, createdBoard.Title)

	boardID := createdBoard.ID

	column := entity.Column{
		UserID:   userID,
		BoardID:  boardID,
		Title:    "Column Title",
		Position: 0,
	}

	err = ts.uc.CreateColumn(ts.ctx, &column)

	if err != nil {
		log.Fatalf("Failed to execute CreateColumn usecase: %v", err)
	}

	var createdColumn repository.Column
	err = db.GetContext(ts.ctx, &createdColumn, `
		SELECT * FROM columns WHERE user_id = $1
	`, userID)

	if err != nil {
		log.Fatalf("Failed to select created column: %v", err)
	}

	assert.Equal(t, column.UserID, createdColumn.UserID)
	assert.Equal(t, column.BoardID, createdColumn.BoardID)
	assert.Equal(t, column.Title, createdColumn.Title)

	columnID := createdColumn.ID

	card := entity.Card{
		UserID:   userID,
		ColumnID: columnID,
		Title:    "Card Title",
		Position: 0,
	}

	err = ts.uc.CreateCard(ts.ctx, &card)

	if err != nil {
		log.Fatalf("Failed to execute CreateCard usecase: %v", err)
	}

	var createdCard repository.Card
	err = db.GetContext(ts.ctx, &createdCard, `
		SELECT * FROM cards WHERE user_id = $1
	`, userID)

	if err != nil {
		log.Fatalf("Failed to select created card: %v", err)
	}

	assert.Equal(t, card.UserID, createdCard.UserID)
	assert.Equal(t, card.ColumnID, createdCard.ColumnID)
	assert.Equal(t, card.Title, createdCard.Title)
}

// GetBoardByID(ctx context.Context, id uuid.UUID) (*entity.Board, error)
func TestGetByID(t *testing.T) {
	ts := sqlxSetup()
	resetDatabase()

	boardID := uuid.New()
	board := repository.Board{
		ID:     boardID,
		UserID: uuid.New(),
		Title:  "Board Title",
	}
	query := `
		INSERT INTO boards (id, user_id, title)
		VALUES (:id, :user_id, :title)
	`
	_, err := db.NamedExecContext(ts.ctx, query, &board)

	if err != nil {
		log.Fatalf("Failed to insert into boards: %v", err)
	}

	gotBoard, err := ts.uc.GetBoardByID(ts.ctx, boardID)

	if err != nil {
		log.Fatalf("Failed to execute GetBoardByID usecase: %v", err)
	}

	assert.Equal(t, board.ID, gotBoard.ID)
	assert.Equal(t, board.UserID, gotBoard.UserID)
	assert.Equal(t, board.Title, gotBoard.Title)

	// TODO: Columns, cards
}

// UpdateBoard(ctx context.Context, board *entity.Board) error
func TestUpdate(t *testing.T) {
	ts := sqlxSetup()
	resetDatabase()

	boardID := uuid.New()
	board := repository.Board{
		ID:     boardID,
		UserID: uuid.New(),
		Title:  "Board Title",
	}
	query := `
		INSERT INTO boards (id, user_id, title)
		VALUES (:id, :user_id, :title)
	`
	_, err := db.NamedExecContext(ts.ctx, query, &board)

	if err != nil {
		log.Fatalf("Failed to insert into boards: %v", err)
	}

	newBoard := repository.BoardToEntity(board)
	newBoard.Title = "New Board Title"

	err = ts.uc.UpdateBoard(ts.ctx, &newBoard)

	if err != nil {
		log.Fatalf("Failed to execute UpdateBoard usecase: %v", err)
	}

	var updatedBoard repository.Board

	err = db.GetContext(ts.ctx, &updatedBoard, `
		SELECT * FROM boards WHERE id = $1
	`, boardID)

	if err != nil {
		log.Fatalf("Failed to select updated board: %v", err)
	}

	assert.Equal(t, newBoard.Title, updatedBoard.Title)

	// TODO: Columns, cards
}

// DeleteBoard(ctx context.Context, id uuid.UUID) error
func TestDelete(t *testing.T) {
	ts := sqlxSetup()
	resetDatabase()

	boardID := uuid.New()
	board := repository.Board{
		ID:     boardID,
		UserID: uuid.New(),
		Title:  "Board Title",
	}
	query := `
		INSERT INTO boards (id, user_id, title)
		VALUES (:id, :user_id, :title)
	`
	_, err := db.NamedExecContext(ts.ctx, query, &board)

	if err != nil {
		log.Fatalf("Failed to insert into boards: %v", err)
	}

	err = ts.uc.DeleteBoard(ts.ctx, boardID)

	if err != nil {
		log.Fatalf("Failed to execute DeleteBoard usecase: %v", err)
	}

	var tmp repository.Board
	err = db.GetContext(ts.ctx, &tmp, `
		SELECT * FROM boards WHERE id = $1
	`, boardID)

	assert.NotNil(t, err)

	// TODO: Columns, cards
}

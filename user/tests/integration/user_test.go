package integration_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	sqlxRepository "user/internal/adapter/repository/sqlx"
	"user/internal/entity"
	"user/internal/repository"
	"user/internal/usecase"
	v1 "user/internal/usecase/v1"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

var db *sqlx.DB

type testSetup struct {
	ctx  context.Context
	repo repository.UserRepository
	uc   usecase.UserUseCase
}

func sqlxSetup() *testSetup {
	ctx := context.TODO()
	repo := sqlxRepository.NewSQLXUserRepository(db)
	uc := v1.NewUserUseCase(repo)

	return &testSetup{
		ctx:  ctx,
		repo: repo,
		uc:   uc,
	}
}

func applyMigrations(dsn string) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("Failed to create migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///home/rukost/University/software-design-s6-bmstu.git/lab4/src/user/migrations/sql",
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

func resetDatabase() error {
	_, err := db.Exec(`
	TRUNCATE TABLE users RESTART IDENTITY CASCADE
	`)
	return err
}

func TestCreateUser(t *testing.T) {
	ts := sqlxSetup()
	resetDatabase()

	user := entity.User{
		Username:     "username",
		Email:        "email@test.com",
		PasswordHash: "Pa$$w0rD",
	}

	err := ts.uc.CreateUser(ts.ctx, user)

	if err != nil {
		log.Fatalf("Failed to execute CreateUser usecase: %v", err)
	}

	var createdUser repository.User

	err = db.GetContext(ts.ctx, &createdUser, `
		SELECT * FROM users WHERE username = $1
	`, user.Username)

	if err != nil {
		log.Fatalf("Failed to select created user: %v", err)
	} else {
		assert.Equal(t, user.Username, createdUser.Username)
		assert.Equal(t, user.Email, createdUser.Email)
		err = bcrypt.CompareHashAndPassword([]byte(createdUser.PasswordHash), []byte(user.PasswordHash))
		assert.Nil(t, err)
	}
}

func TestGetUserByID(t *testing.T) {
	ts := sqlxSetup()
	resetDatabase()

	id := uuid.New()

	query := `
    INSERT INTO users (id, username, email, password_hash)
    VALUES (:id, :username, :email, :password_hash)
    `
	_, err := db.NamedExecContext(ts.ctx, query, &repository.User{
		ID:           id,
		Username:     "username",
		Email:        "email@test.com",
		PasswordHash: "somePasswordHash",
	})
	if err != nil {
		log.Fatalf("Failed to insert into users: %v", err)
	}

	user, err := ts.uc.GetUserByID(ts.ctx, id)

	if err != nil {
		log.Fatalf("Failed to execute GetUserByID usecase: %v", err)
	} else {
		assert.Equal(t, "username", user.Username)
		assert.Equal(t, "email@test.com", user.Email)
	}
}

func TestUpdateUser(t *testing.T) {
	ts := sqlxSetup()
	resetDatabase()

	id := uuid.New()

	query := `
    INSERT INTO users (id, username, email, password_hash)
    VALUES (:id, :username, :email, :password_hash)
    `
	_, err := db.NamedExecContext(ts.ctx, query, &repository.User{
		ID:           id,
		Username:     "username",
		Email:        "email@test.com",
		PasswordHash: "somePasswordHash",
	})
	if err != nil {
		log.Fatalf("Failed to insert into users: %v", err)
	}

	updatedUserInfo := &entity.User{
		ID:       id,
		Username: "new_username",
		Email:    "new_email@test.com",
	}

	err = ts.uc.UpdateUser(ts.ctx, updatedUserInfo)

	if err != nil {
		log.Fatalf("Failed to execute UpdateUser usecase: %v", err)
	}

	var updatedUser repository.User

	err = db.GetContext(ts.ctx, &updatedUser, `
		SELECT * FROM users WHERE id = $1
	`, id)

	if err != nil {
		log.Fatalf("Failed to select updated user: %v", err)
	} else {
		assert.Equal(t, "new_username", updatedUser.Username)
		assert.Equal(t, "new_email@test.com", updatedUser.Email)
	}
}

func TestDeleteUser(t *testing.T) {
	ts := sqlxSetup()
	resetDatabase()

	id := uuid.New()

	query := `
    INSERT INTO users (id, username, email, password_hash)
    VALUES (:id, :username, :email, :password_hash)
    `
	_, err := db.NamedExecContext(ts.ctx, query, &repository.User{
		ID:           id,
		Username:     "username",
		Email:        "email@test.com",
		PasswordHash: "somePasswordHash",
	})
	if err != nil {
		log.Fatalf("Failed to insert into users: %v", err)
	}

	err = ts.uc.DeleteUser(ts.ctx, id)

	if err != nil {
		log.Fatalf("Failed to execute DeleteUser usecase: %v", err)
	} else {
		var tmp repository.User
		err = db.GetContext(ts.ctx, &tmp, `
			SELECT * FROM users WHERE id = $1
		`, id)
		assert.NotNil(t, err)
	}
}

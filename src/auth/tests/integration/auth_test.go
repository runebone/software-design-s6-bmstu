package integration_test

import (
	sqlxRepository "auth/internal/adapter/repository/sqlx"
	"auth/internal/adapter/service/tokengen/jwt"
	"auth/internal/dto"
	"auth/internal/entity"
	"auth/internal/repository"
	"auth/internal/service/tokengen"
	"auth/internal/usecase"
	v1 "auth/internal/usecase/v1"
	"auth/mocks"
	"context"
	"fmt"
	"log"
	"os"
	"testing"

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
	ctx      context.Context
	repo     repository.TokenRepository
	userSvc  *mocks.UserService
	tokenSvc tokengen.TokenService
	uc       usecase.AuthUsecase
}

func sqlxSetup() *testSetup {
	ctx := context.TODO()
	repo := sqlxRepository.NewSQLXTokenRepository(db)

	userSvc := new(mocks.UserService)

	tokenSvc := jwt.NewJWTService("secret", 15*60, 7*24*60*60)

	uc := v1.NewAuthUseCase(repo, userSvc, tokenSvc)

	return &testSetup{
		ctx:      ctx,
		repo:     repo,
		userSvc:  userSvc,
		tokenSvc: tokenSvc,
		uc:       uc,
	}
}

func applyMigrations(dsn string) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("Failed to create migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///home/rukost/University/software-design-s6-bmstu.git/lab4/src/auth/migrations/sql",
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
	TRUNCATE TABLE tokens RESTART IDENTITY CASCADE
	`)
	return err
}

// Login(ctx context.Context, email, password string) (*dto.LoginResponse, error)
func TestLogin(t *testing.T) {
	ts := sqlxSetup()
	resetDatabase()

	id := uuid.New()
	username := "username"
	email := "test@email.com"
	password := "Pa$$w0rD"
	password_hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &dto.User{
		ID:           id,
		Username:     username,
		Email:        email,
		PasswordHash: string(password_hash),
	}

	// GetUserByEmail(ctx context.Context, email string) (*dto.User, error)
	ts.userSvc.On("GetUserByEmail", ts.ctx, email).Return(user, nil)

	loginResponse, err := ts.uc.Login(ts.ctx, email, password)
	if err != nil {
		log.Fatalf("Failed to execute Login usecase: %v", err)
	}

	var token entity.Token
	err = db.GetContext(ts.ctx, &token, `
		SELECT * FROM tokens WHERE user_id = $1
	`, id)
	if err != nil {
		log.Fatalf("Failed saving refresh token to repository: %v", err)
	}

	accessToken := loginResponse.AccessToken
	refreshToken := loginResponse.RefreshToken

	user_id, err := ts.tokenSvc.ValidateToken(ts.ctx, accessToken)
	if err != nil {
		log.Fatalf("Failed to validate access token: %v", err)
	}

	assert.Equal(t, id.String(), user_id)

	user_id, err = ts.tokenSvc.ValidateToken(ts.ctx, refreshToken)
	if err != nil {
		log.Fatalf("Failed to validate refresh token: %v", err)
	}

	assert.Equal(t, id.String(), user_id)
}

// Refresh(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponse, error)
func TestRefresh(t *testing.T) {
	ts := sqlxSetup()
	resetDatabase()

	id, _ := uuid.Parse("4457092b-2b06-4443-81af-323c2e67d000")
	refreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2NzgzMDksImlhdCI6MTcyNzYzMDMwOSwic3ViIjoiNDQ1NzA5MmItMmIwNi00NDQzLTgxYWYtMzIzYzJlNjdkMDAwIiwidHlwZSI6InJlZnJlc2gifQ._tO3m754wHTIqJqp3WyYFeHFFuOSYLCjMCuU167xn7U"

	refreshTokenResponse, err := ts.uc.Refresh(ts.ctx, refreshToken)
	if err != nil {
		log.Fatalf("Failed to execute Refresh usecase: %v", err)
	}

	accessToken := refreshTokenResponse.AccessToken

	user_id, err := ts.tokenSvc.ValidateToken(ts.ctx, accessToken)
	if err != nil {
		log.Fatalf("Failed to validate access token: %v", err)
	}

	assert.Equal(t, id.String(), user_id)
}

// Logout(ctx context.Context, refreshToken string) error
func TestLogout(t *testing.T) {
	ts := sqlxSetup()
	resetDatabase()

	id := uuid.New()
	user_id := "4457092b-2b06-4443-81af-323c2e67d000"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2NzgzMDksImlhdCI6MTcyNzYzMDMwOSwic3ViIjoiNDQ1NzA5MmItMmIwNi00NDQzLTgxYWYtMzIzYzJlNjdkMDAwIiwidHlwZSI6InJlZnJlc2gifQ._tO3m754wHTIqJqp3WyYFeHFFuOSYLCjMCuU167xn7U"

	query := `
    INSERT INTO tokens (id, user_id, token)
    VALUES ($1, $2, $3)
    `
	_, err := db.ExecContext(ts.ctx, query, id, user_id, token)
	if err != nil {
		log.Fatalf("Failed to insert into tokens: %v", err)
	}

	err = ts.uc.Logout(ts.ctx, token)

	if err != nil {
		log.Fatalf("Failed to execute Logout usecase: %v", err)
	} else {
		var tmp entity.Token
		err = db.GetContext(ts.ctx, &tmp, `
			SELECT * FROM tokens WHERE user_id = $1
		`, user_id)
		assert.NotNil(t, err)
	}
}

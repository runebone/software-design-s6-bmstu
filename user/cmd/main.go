package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	_ "time/tzdata"
	"user/internal/adapter/database"
	"user/internal/adapter/logger"
	mongoRepo "user/internal/adapter/repository/mongo"
	sqlxRepo "user/internal/adapter/repository/sqlx"
	api "user/internal/api/v1"
	"user/internal/config"
	handler "user/internal/handler/v1"
	"user/internal/middleware"
	"user/internal/repository"
	usecase "user/internal/usecase/v1"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatalf("Couldn't set timezone: %v", err)
	}
	time.Local = loc
}

type dbrepo interface {
	DB() (any, error)
	Repo(any) any
}

type postgres struct {
	cfg *config.Config
}

// func (p *postgres) DB(cfg config.PostgresConfig) (*sqlx.DB, error) {
func (p *postgres) DB() (any, error) {
	return database.NewPostgresDB(p.cfg.User.Postgres)
}

// func (p *postgres) Repo(db *sqlx.DB) *sqlxRepo.SQLXUserRepository {
func (p *postgres) Repo(db any) any {
	return sqlxRepo.NewSQLXUserRepository(db.(*sqlx.DB))
}

type mongodb struct {
	cfg *config.Config
}

// func (m *mongodb) DB(cfg config.MongoConfig) (*mongo.Database, error) {
func (m *mongodb) DB() (any, error) {
	return database.NewMongoDB(m.cfg.User.Mongo)
}

// func (m *mongodb) Repo(db *mongo.Database) *mongoRepo.MongoUserRepository {
func (m *mongodb) Repo(db any) any {
	return mongoRepo.NewMongoUserRepository(db.(*mongo.Database))
}

func main() {
	config, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Println("Error reading config (config.toml)")
	}

	dbmap := make(map[string]dbrepo)
	dbmap["postgres"] = &postgres{cfg: config}
	dbmap["mongo"] = &mongodb{cfg: config}

	dbRepo := dbmap[config.User.Database]

	// db, err := database.NewPostgresDB(config.User.Postgres)
	db, err := dbRepo.DB()
	if err != nil {
		log.Fatalf("Couldn't connect to database, exiting: %w", err)
	}

	logger := logger.NewZapLogger(config.User.Log)

	// repo := sqlxRepo.NewSQLXUserRepository(db)
	repo := dbRepo.Repo(db).(repository.UserRepository)
	uc := usecase.NewUserUseCase(repo, logger)

	userHandler := handler.NewUserHandler(uc)
	router := mux.NewRouter()
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)
	router.Use(loggingMiddleware.Middleware)
	api.InitializeV1Routes(router, userHandler)

	localPort := fmt.Sprintf("%d", config.User.LocalPort)
	exposedPort := fmt.Sprintf("%d", config.User.ExposedPort)

	log.Printf("Starting server on :%s\n", exposedPort)
	http.ListenAndServe(":"+localPort, router)
}

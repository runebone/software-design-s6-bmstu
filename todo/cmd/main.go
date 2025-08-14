package main

import (
	"fmt"
	"time"
	_ "time/tzdata"
	"todo/internal/adapter/database"
	"todo/internal/adapter/logger"

	"log"
	"net/http"
	sqlxRepo "todo/internal/adapter/repository/sqlx"
	api "todo/internal/api/v1"
	"todo/internal/config"
	handler "todo/internal/handler/v1"
	"todo/internal/middleware"
	usecase "todo/internal/usecase/v1"

	"github.com/gorilla/mux"
)

func init() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatalf("Couldn't set timezone: %v", err)
	}
	time.Local = loc
}

func main() {
	config, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Println("Error reading config (config.toml)")
	}

	db, err := database.NewPostgresDB(config.Todo.Postgres)
	if err != nil {
		log.Println("Couldn't connect to database, exiting")
		return
	}

	logger := logger.NewZapLogger(config.Todo.Log)

	boardRepo := sqlxRepo.NewSQLXBoardRepository(db)
	columnRepo := sqlxRepo.NewSQLXColumnRepository(db)
	cardRepo := sqlxRepo.NewSQLXCardRepository(db)

	uc := usecase.NewTodoUseCase(boardRepo, columnRepo, cardRepo, logger)

	userHandler := handler.NewTodoHandler(uc, config.Pagination)
	router := mux.NewRouter()
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)
	router.Use(loggingMiddleware.Middleware)
	api.InitializeV1Routes(router, userHandler)

	localPort := fmt.Sprintf("%d", config.Todo.LocalPort)
	exposedPort := fmt.Sprintf("%d", config.Todo.ExposedPort)

	log.Printf("Starting server on :%s\n", exposedPort)
	http.ListenAndServe(":"+localPort, router)
}

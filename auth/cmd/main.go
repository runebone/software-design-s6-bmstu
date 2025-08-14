package main

import (
	"auth/internal/adapter/database"
	"auth/internal/adapter/logger"
	"fmt"
	"time"
	_ "time/tzdata"

	sqlxRepo "auth/internal/adapter/repository/sqlx"
	"auth/internal/adapter/service/tokengen/jwt"
	user "auth/internal/adapter/service/user/http"
	api "auth/internal/api/v1"
	"auth/internal/config"
	handler "auth/internal/handler/v1"
	"auth/internal/middleware"
	usecase "auth/internal/usecase/v1"
	"log"
	"net/http"

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

	logger := logger.NewZapLogger(config.Auth.Log)

	db, err := database.NewPostgresDB(config.Auth.Postgres)
	if err != nil {
		log.Println("Couldn't connect to database, exiting")
		return
	}

	repo := sqlxRepo.NewSQLXTokenRepository(db)

	baseURL := fmt.Sprintf("http://%s:%d/%s", config.User.ContainerName, config.User.LocalPort, config.User.BaseURL)

	userService := user.NewHTTPUserService(baseURL, 2*time.Second)
	tokenService := jwt.NewJWTService(
		config.Auth.Token.Secret,
		time.Duration(config.Auth.Token.AccessTTL)*time.Second,
		time.Duration(config.Auth.Token.RefreshTTL)*time.Second,
	)

	uc := usecase.NewAuthUseCase(repo, userService, tokenService, logger)

	userHandler := handler.NewAuthHandler(uc)
	router := mux.NewRouter()
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)
	router.Use(loggingMiddleware.Middleware)
	api.InitializeV1Routes(router, userHandler)

	localPort := fmt.Sprintf("%d", config.Auth.LocalPort)
	exposedPort := fmt.Sprintf("%d", config.Auth.ExposedPort)

	log.Printf("Starting server on :%s\n", exposedPort)
	http.ListenAndServe(":"+localPort, router)
}

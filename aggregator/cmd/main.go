package main

import (
	"aggregator/internal/adapter/logger"
	"aggregator/internal/middleware"
	"fmt"
	"time"
	_ "time/tzdata"

	httpAuth "aggregator/internal/adapter/service/auth/http"
	httpTodo "aggregator/internal/adapter/service/todo/http"
	httpUser "aggregator/internal/adapter/service/user/http"
	api "aggregator/internal/api/v1"
	"aggregator/internal/config"
	h "aggregator/internal/handler/v1"
	"aggregator/internal/service/auth"
	"aggregator/internal/service/todo"
	"aggregator/internal/service/user"
	v1 "aggregator/internal/usecase/v1"
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

	logger := logger.NewZapLogger(config.Aggregator.Log)

	var userSvc user.UserService
	{
		baseURL := fmt.Sprintf("http://%s:%d/%s", config.User.ContainerName, config.User.LocalPort, config.User.BaseURL)
		userSvc = httpUser.NewUserService(baseURL, 2*time.Second, logger)
	}

	var authSvc auth.AuthService
	{
		baseURL := fmt.Sprintf("http://%s:%d/%s", config.Auth.ContainerName, config.Auth.LocalPort, config.Auth.BaseURL)
		authSvc = httpAuth.NewAuthService(baseURL, 2*time.Second, logger)
	}

	var todoSvc todo.TodoService
	{
		baseURL := fmt.Sprintf("http://%s:%d/%s", config.Todo.ContainerName, config.Todo.LocalPort, config.Todo.BaseURL)
		todoSvc = httpTodo.NewTodoService(baseURL, 2*time.Second, logger)
	}

	uc := v1.NewAggregatorUseCase(userSvc, authSvc, todoSvc, logger)
	handler := h.NewAggregatorHandler(uc)

	router := mux.NewRouter()
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)
	router.Use(loggingMiddleware.Middleware)
	authMiddleware := middleware.NewAuthMiddleware(authSvc)
	api.InitializeV1Routes(router, handler, authMiddleware)

	localPort := fmt.Sprintf("%d", config.Aggregator.LocalPort)
	exposedPort := fmt.Sprintf("%d", config.Aggregator.ExposedPort)

	log.Printf("Starting server on :%s\n", exposedPort)
	http.ListenAndServe(":"+localPort, router)
}

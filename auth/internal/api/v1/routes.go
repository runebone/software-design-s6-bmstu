package v1

import (
	v1 "auth/internal/handler/v1"

	"github.com/gorilla/mux"
)

func InitializeV1Routes(router *mux.Router, authHandler *v1.AuthHandler) {
	router.HandleFunc("/api/v1/register", authHandler.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/v1/login", authHandler.LoginHandler).Methods("POST")
	router.HandleFunc("/api/v1/refresh", authHandler.RefreshTokenHandler).Methods("POST")
	router.HandleFunc("/api/v1/validate", authHandler.ValidateTokenHandler).Methods("POST")
	router.HandleFunc("/api/v1/logout", authHandler.LogoutHandler).Methods("POST")
}

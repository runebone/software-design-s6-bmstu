package v1

import (
	v1 "user/internal/handler/v1"

	"github.com/gorilla/mux"
)

func InitializeV1Routes(router *mux.Router, userHandler *v1.UserHandler) {
	router.HandleFunc("/api/v1/users", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/api/v1/users/new", userHandler.GetNewUsers).Methods("GET")
	router.HandleFunc("/api/v1/users/{id}", userHandler.GetUserByID).Methods("GET")
	router.HandleFunc("/api/v1/users", userHandler.GetUsers).Methods("GET")
	// TODO: GetUsersBatch

	// TODO: PUT and DELETE requests should require authroization
	router.HandleFunc("/api/v1/users", userHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/api/v1/users", userHandler.DeleteUser).Methods("DELETE")
}

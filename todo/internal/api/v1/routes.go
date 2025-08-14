package v1

import (
	v1 "todo/internal/handler/v1"

	"github.com/gorilla/mux"
)

func InitializeV1Routes(router *mux.Router, todoHandler *v1.TodoHandler) {
	router.HandleFunc("/api/v1/boards", todoHandler.CreateBoard).Methods("POST")
	router.HandleFunc("/api/v1/boards/{id}", todoHandler.GetBoardByID).Methods("GET")
	router.HandleFunc("/api/v1/boards", todoHandler.GetBoardsByUser).Methods("GET")
	router.HandleFunc("/api/v1/boards", todoHandler.UpdateBoard).Methods("PUT")
	router.HandleFunc("/api/v1/boards", todoHandler.DeleteBoard).Methods("DELETE")

	router.HandleFunc("/api/v1/columns", todoHandler.CreateColumn).Methods("POST")
	router.HandleFunc("/api/v1/columns/{id}", todoHandler.GetColumnByID).Methods("GET")
	router.HandleFunc("/api/v1/columns", todoHandler.GetColumnsByBoard).Methods("GET")
	router.HandleFunc("/api/v1/columns", todoHandler.UpdateColumn).Methods("PUT")
	router.HandleFunc("/api/v1/columns", todoHandler.DeleteColumn).Methods("DELETE")

	router.HandleFunc("/api/v1/cards", todoHandler.CreateCard).Methods("POST")
	router.HandleFunc("/api/v1/cards/new", todoHandler.GetNewCards).Methods("GET")
	router.HandleFunc("/api/v1/cards/{id}", todoHandler.GetCardByID).Methods("GET")
	router.HandleFunc("/api/v1/cards", todoHandler.GetCardsByColumn).Methods("GET")
	router.HandleFunc("/api/v1/cards", todoHandler.UpdateCard).Methods("PUT")
	router.HandleFunc("/api/v1/cards", todoHandler.DeleteCard).Methods("DELETE")
}

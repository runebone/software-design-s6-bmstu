package v1

import (
	v1 "aggregator/internal/handler/v1"
	"aggregator/internal/middleware"

	"github.com/gorilla/mux"
)

func InitializeV1Routes(router *mux.Router, aggHandler *v1.AggregatorHandler, authMiddleware *middleware.AuthMiddleware) {
	router.HandleFunc("/api/v1/register", aggHandler.Register).Methods("POST")
	router.HandleFunc("/api/v1/login", aggHandler.Login).Methods("POST")
	router.HandleFunc("/api/v1/refresh", aggHandler.Refresh).Methods("POST")
	router.HandleFunc("/api/v1/validate", aggHandler.Validate).Methods("POST")
	router.HandleFunc("/api/v1/logout", aggHandler.Logout).Methods("POST")

	authRoutes := router.PathPrefix("/api/v1").Subrouter()
	authRoutes.Use(authMiddleware.Middleware)

	authRoutes.HandleFunc("/boards", aggHandler.GetBoards).Methods("GET")      // Boards
	authRoutes.HandleFunc("/board/{id}", aggHandler.GetBoard).Methods("GET")   // Columns + cards
	authRoutes.HandleFunc("/column/{id}", aggHandler.GetColumn).Methods("GET") // Cards
	authRoutes.HandleFunc("/card/{id}", aggHandler.GetCard).Methods("GET")     // Card + description

	authRoutes.HandleFunc("/board", aggHandler.CreateBoard).Methods("POST")
	authRoutes.HandleFunc("/column", aggHandler.CreateColumn).Methods("POST")
	authRoutes.HandleFunc("/card", aggHandler.CreateCard).Methods("POST")

	authRoutes.HandleFunc("/board", aggHandler.UpdateBoard).Methods("PUT")
	authRoutes.HandleFunc("/column", aggHandler.UpdateColumn).Methods("PUT")
	authRoutes.HandleFunc("/card", aggHandler.UpdateCard).Methods("PUT")

	authRoutes.HandleFunc("/board/{id}", aggHandler.DeleteBoard).Methods("DELETE")
	authRoutes.HandleFunc("/column/{id}", aggHandler.DeleteColumn).Methods("DELETE")
	authRoutes.HandleFunc("/card/{id}", aggHandler.DeleteCard).Methods("DELETE")

	authRoutes.HandleFunc("/stats/{from}/{to}", aggHandler.GetStats).Methods("GET")
	authRoutes.HandleFunc("/stats/{from}", aggHandler.GetStats).Methods("GET")
	authRoutes.HandleFunc("/stats", aggHandler.GetStats).Methods("GET")
}

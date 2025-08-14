package handler

import "net/http"

type AggregatorHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	Validate(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)

	GetBoards(w http.ResponseWriter, r *http.Request)
	GetBoard(w http.ResponseWriter, r *http.Request)
	GetColumn(w http.ResponseWriter, r *http.Request)
	GetCard(w http.ResponseWriter, r *http.Request)
	GetStats(w http.ResponseWriter, r *http.Request)

	CreateBoard(w http.ResponseWriter, r *http.Request)
	CreateColumn(w http.ResponseWriter, r *http.Request)
	CreateCard(w http.ResponseWriter, r *http.Request)

	UpdateBoard(w http.ResponseWriter, r *http.Request)
	UpdateColumn(w http.ResponseWriter, r *http.Request)
	UpdateCard(w http.ResponseWriter, r *http.Request)

	DeleteBoard(w http.ResponseWriter, r *http.Request)
	DeleteColumn(w http.ResponseWriter, r *http.Request)
	DeleteCard(w http.ResponseWriter, r *http.Request)
}

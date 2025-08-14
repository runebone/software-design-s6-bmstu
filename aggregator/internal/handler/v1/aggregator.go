package v1

import (
	"aggregator/internal/dto"
	"aggregator/internal/middleware"
	"aggregator/internal/usecase"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const layout = "02-01-2006"

var (
	ErrInvalidRequestBody error = errors.New("invalid request body")
	ErrNoUserID           error = errors.New("couldn't get userID from context")
	ErrBadUserID          error = errors.New("couldn't parse userID")
	ErrNoRole             error = errors.New("couldn't get role from context")
	ErrNotAdmin           error = errors.New("not admin")
)

type AggregatorHandler struct {
	uc usecase.AggregatorUseCase
}

func NewAggregatorHandler(uc usecase.AggregatorUseCase) *AggregatorHandler {
	return &AggregatorHandler{
		uc: uc,
	}
}

func (h *AggregatorHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	tokens, err := h.uc.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(tokens)
}

func (h *AggregatorHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	tokens, err := h.uc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(tokens)
}

func (h *AggregatorHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.uc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *AggregatorHandler) Validate(w http.ResponseWriter, r *http.Request) {
	var req dto.ValidateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.uc.Validate(r.Context(), req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *AggregatorHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	err := h.uc.Logout(r.Context(), req.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
}

func (h *AggregatorHandler) GetBoards(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, ErrNoUserID.Error(), http.StatusUnauthorized)
		return
	}

	boards, err := h.uc.GetBoards(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	json.NewEncoder(w).Encode(boards)
}

func (h *AggregatorHandler) GetBoard(w http.ResponseWriter, r *http.Request) {
	boardID := mux.Vars(r)["id"]

	columns, err := h.uc.GetColumns(r.Context(), boardID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	json.NewEncoder(w).Encode(columns)
}

func (h *AggregatorHandler) GetColumn(w http.ResponseWriter, r *http.Request) {
	columnID := mux.Vars(r)["id"]

	cards, err := h.uc.GetCards(r.Context(), columnID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	json.NewEncoder(w).Encode(cards)
}

func (h *AggregatorHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	cardID := mux.Vars(r)["id"]

	card, err := h.uc.GetCard(r.Context(), cardID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	json.NewEncoder(w).Encode(card)
}

func (h *AggregatorHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	// XXX: Role check better should be in another role checking middleware
	role, ok := middleware.GetRoleFromContext(r.Context())
	if !ok {
		http.Error(w, ErrNoRole.Error(), http.StatusUnauthorized)
		return
	}

	if role != "admin" {
		http.Error(w, ErrNotAdmin.Error(), http.StatusUnauthorized)
		return
	}

	// XXX: Read default values from config
	fromStr, ok := mux.Vars(r)["from"]
	if !ok {
		fromStr = "01-01-1980"
	}
	from, _ := time.Parse(layout, fromStr)

	toStr, ok := mux.Vars(r)["to"]
	if !ok {
		toStr = "01-01-9999"
	}
	to, _ := time.Parse(layout, toStr)

	stats, err := h.uc.GetStats(r.Context(), from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

func (h *AggregatorHandler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	userIDstr, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, ErrNoUserID.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		http.Error(w, ErrBadUserID.Error(), http.StatusUnauthorized)
		return
	}

	board := dto.Board{
		UserID: userID,
		Title:  req.Title,
	}

	err = h.uc.CreateBoard(r.Context(), board)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AggregatorHandler) CreateColumn(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateColumnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	userIDstr, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, ErrNoUserID.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		http.Error(w, ErrBadUserID.Error(), http.StatusUnauthorized)
		return
	}

	column := dto.Column{
		UserID:  userID,
		BoardID: req.BoardID,
		Title:   req.Title,
	}

	err = h.uc.CreateColumn(r.Context(), column)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AggregatorHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	userIDstr, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, ErrNoUserID.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		http.Error(w, ErrBadUserID.Error(), http.StatusUnauthorized)
		return
	}

	card := dto.Card{
		UserID:      userID,
		ColumnID:    req.ColumnID,
		Title:       req.Title,
		Description: req.Description,
	}

	err = h.uc.CreateCard(r.Context(), card)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AggregatorHandler) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	board := dto.Board{
		ID:    req.ID,
		Title: req.Title,
	}

	err := h.uc.UpdateBoard(r.Context(), &board)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
}

func (h *AggregatorHandler) UpdateColumn(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateColumnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	userIDstr, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, ErrNoUserID.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		http.Error(w, ErrBadUserID.Error(), http.StatusUnauthorized)
		return
	}

	column := dto.Column{
		ID:      req.ID,
		UserID:  userID,
		BoardID: req.BoardID,
		Title:   req.Title,
	}

	err = h.uc.UpdateColumn(r.Context(), &column)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
}

func (h *AggregatorHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	userIDstr, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, ErrNoUserID.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		http.Error(w, ErrBadUserID.Error(), http.StatusUnauthorized)
		return
	}

	card := dto.Card{
		ID:          req.ID,
		UserID:      userID,
		ColumnID:    req.ColumnID,
		Title:       req.Title,
		Description: req.Description,
	}

	err = h.uc.UpdateCard(r.Context(), &card)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
}

func (h *AggregatorHandler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := h.uc.DeleteBoard(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
}

func (h *AggregatorHandler) DeleteColumn(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := h.uc.DeleteColumn(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
}

func (h *AggregatorHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := h.uc.DeleteCard(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
}

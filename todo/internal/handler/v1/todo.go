package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"todo/internal/config"
	"todo/internal/dto"
	"todo/internal/entity"
	"todo/internal/usecase"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	ErrInvalidUserID   = "invalid user id"
	ErrInvalidBoardID  = "invalid board id"
	ErrInvalidColumnID = "invalid column id"
	ErrInvalidCardID   = "invalid card id"
	ErrInvalidFromDate = "invalid <<from>> date"
	ErrInvalidToDate   = "invalid <<to>> date"
)

type TodoHandler struct {
	todoUseCase usecase.TodoUseCase
	config      config.PaginationConfig
}

func NewTodoHandler(todoUseCase usecase.TodoUseCase, config config.PaginationConfig) *TodoHandler {
	return &TodoHandler{todoUseCase: todoUseCase, config: config}
}

func (h *TodoHandler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateBoardRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	board := &entity.Board{
		UserID: input.UserID,
		Title:  input.Title,
	}

	err := h.todoUseCase.CreateBoard(r.Context(), board)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *TodoHandler) GetBoardByID(w http.ResponseWriter, r *http.Request) {
	boardID := mux.Vars(r)["id"]
	id, err := uuid.Parse(boardID)

	if err != nil {
		http.Error(w, ErrInvalidBoardID, http.StatusBadRequest)
		return
	}

	board, err := h.todoUseCase.GetBoardByID(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	boardDTO := dto.ToBoardDTO(board)

	json.NewEncoder(w).Encode(boardDTO)
}

func (h *TodoHandler) GetBoardsByUser(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userID := query.Get("user_id")
	id, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, ErrInvalidUserID, http.StatusBadRequest)
		return
	}

	limit := h.config.Limit
	if _, ok := query["limit"]; ok {
		limitStr := query.Get("limit")
		limitInt, err := strconv.Atoi(limitStr)
		if err == nil {
			limit = limitInt
		}
	}

	offset := h.config.Offset
	if _, ok := query["offset"]; ok {
		offsetStr := query.Get("offset")
		offsetInt, err := strconv.Atoi(offsetStr)
		if err == nil {
			offset = offsetInt
		}
	}

	boards, err := h.todoUseCase.GetBoardsByUser(r.Context(), id, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	boardDTOs := dto.ToBoardDTOs(boards)

	json.NewEncoder(w).Encode(boardDTOs)
}

func (h *TodoHandler) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	var input dto.UpdateBoardRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	board := &entity.Board{
		ID:    input.ID,
		Title: input.Title,
	}

	err := h.todoUseCase.UpdateBoard(r.Context(), board)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TodoHandler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	boardID := query.Get("id")
	id, err := uuid.Parse(boardID)

	if err != nil {
		http.Error(w, ErrInvalidBoardID, http.StatusBadRequest)
		return
	}

	err = h.todoUseCase.DeleteBoard(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TodoHandler) CreateColumn(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateColumnRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	column := &entity.Column{
		UserID:   input.UserID,
		BoardID:  input.BoardID,
		Title:    input.Title,
		Position: input.Position,
	}

	err := h.todoUseCase.CreateColumn(r.Context(), column)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *TodoHandler) GetColumnByID(w http.ResponseWriter, r *http.Request) {
	columnID := mux.Vars(r)["id"]
	id, err := uuid.Parse(columnID)

	if err != nil {
		http.Error(w, ErrInvalidColumnID, http.StatusBadRequest)
		return
	}

	column, err := h.todoUseCase.GetColumnByID(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	columnDTO := dto.ToColumnDTO(column)

	json.NewEncoder(w).Encode(columnDTO)
}

func (h *TodoHandler) GetColumnsByBoard(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	boardID := query.Get("board_id")
	id, err := uuid.Parse(boardID)
	if err != nil {
		http.Error(w, ErrInvalidBoardID, http.StatusBadRequest)
		return
	}

	limit := h.config.Limit
	if _, ok := query["limit"]; ok {
		limitStr := query.Get("limit")
		limitInt, err := strconv.Atoi(limitStr)
		if err == nil {
			limit = limitInt
		}
	}

	offset := h.config.Offset
	if _, ok := query["offset"]; ok {
		offsetStr := query.Get("offset")
		offsetInt, err := strconv.Atoi(offsetStr)
		if err == nil {
			offset = offsetInt
		}
	}

	columns, err := h.todoUseCase.GetColumnsByBoard(r.Context(), id, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	columnDTOs := dto.ToColumnDTOs(columns)

	json.NewEncoder(w).Encode(columnDTOs)
}

func (h *TodoHandler) UpdateColumn(w http.ResponseWriter, r *http.Request) {
	var input dto.UpdateColumnRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	column := &entity.Column{
		ID:       input.ID,
		Title:    input.Title,
		Position: input.Position,
	}

	err := h.todoUseCase.UpdateColumn(r.Context(), column)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TodoHandler) DeleteColumn(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	columnID := query.Get("id")
	id, err := uuid.Parse(columnID)

	if err != nil {
		http.Error(w, ErrInvalidColumnID, http.StatusBadRequest)
		return
	}

	err = h.todoUseCase.DeleteColumn(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TodoHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateCardRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	card := &entity.Card{
		UserID:      input.UserID,
		ColumnID:    input.ColumnID,
		Title:       input.Title,
		Description: input.Description,
		Position:    input.Position,
	}

	err := h.todoUseCase.CreateCard(r.Context(), card)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *TodoHandler) GetCardByID(w http.ResponseWriter, r *http.Request) {
	cardID := mux.Vars(r)["id"]
	id, err := uuid.Parse(cardID)

	if err != nil {
		http.Error(w, ErrInvalidCardID, http.StatusBadRequest)
		return
	}

	card, err := h.todoUseCase.GetCardByID(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cardDTO := dto.ToCardDTO(card)

	json.NewEncoder(w).Encode(cardDTO)
}

func (h *TodoHandler) GetCardsByColumn(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	columnID := query.Get("column_id")
	id, err := uuid.Parse(columnID)
	if err != nil {
		http.Error(w, ErrInvalidColumnID, http.StatusBadRequest)
		return
	}

	limit := h.config.Limit
	if _, ok := query["limit"]; ok {
		limitStr := query.Get("limit")
		limitInt, err := strconv.Atoi(limitStr)
		if err == nil {
			limit = limitInt
		}
	}

	offset := h.config.Offset
	if _, ok := query["offset"]; ok {
		offsetStr := query.Get("offset")
		offsetInt, err := strconv.Atoi(offsetStr)
		if err == nil {
			offset = offsetInt
		}
	}

	cards, err := h.todoUseCase.GetCardsByColumn(r.Context(), id, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cardDTOs := dto.ToCardDTOs(cards)

	json.NewEncoder(w).Encode(cardDTOs)
}

func (h *TodoHandler) GetNewCards(w http.ResponseWriter, r *http.Request) {
	layout := "02-01-2006" // DD-MM-YYYY

	fromParam := r.URL.Query().Get("from")
	from, err := time.Parse(layout, fromParam)
	if err != nil {
		http.Error(w, ErrInvalidFromDate, http.StatusBadRequest)
		return
	}

	toParam := r.URL.Query().Get("to")
	to, err := time.Parse(layout, toParam)
	if err != nil {
		http.Error(w, ErrInvalidToDate, http.StatusBadRequest)
		return
	}

	cards, err := h.todoUseCase.GetNewCards(r.Context(), from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	cardDTOs := dto.ToCardDTOs(cards)

	json.NewEncoder(w).Encode(cardDTOs)
}

func (h *TodoHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	var input dto.UpdateCardRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	card := &entity.Card{
		ID:          input.ID,
		ColumnID:    input.ColumnID,
		Title:       input.Title,
		Description: input.Description,
		Position:    input.Position,
	}

	err := h.todoUseCase.UpdateCard(r.Context(), card)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TodoHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	cardID := query.Get("id")
	id, err := uuid.Parse(cardID)

	if err != nil {
		http.Error(w, ErrInvalidCardID, http.StatusBadRequest)
		return
	}

	err = h.todoUseCase.DeleteCard(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

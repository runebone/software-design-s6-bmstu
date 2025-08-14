package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"user/internal/dto"
	"user/internal/entity"
	"user/internal/repository"
	"user/internal/usecase"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := entity.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: input.Password, // NOTE: Hashing will be done in UseCase
	}

	err := h.userUseCase.CreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(dto.UserDTO{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(dto.UserDTO{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	})
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	id := new(string)
	email := new(string)
	username := new(string)

	if _, ok := query["id"]; ok {
		*id = query.Get("id")
	} else {
		id = nil
	}

	if _, ok := query["email"]; ok {
		*email = query.Get("email")
	} else {
		email = nil
	}

	if _, ok := query["username"]; ok {
		*username = query.Get("username")
	} else {
		username = nil
	}

	filter := repository.UserFilter{
		ID:       id,
		Email:    email,
		Username: username,
	}

	users, err := h.userUseCase.GetUsers(r.Context(), filter)
	if err != nil {
		http.Error(w, "Users not found", http.StatusNotFound)
		return
	}

	userDTOs := dto.ToUserDTOs(users)

	json.NewEncoder(w).Encode(userDTOs)
}

func (h *UserHandler) GetUsersBatch(w http.ResponseWriter, r *http.Request) {
	limitParam := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}

	offsetParam := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		http.Error(w, "Invalid offset", http.StatusBadRequest)
		return
	}

	users, err := h.userUseCase.GetUsersBatch(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "Users not found", http.StatusNotFound)
		return
	}

	userDTOs := dto.ToUserDTOs(users)

	json.NewEncoder(w).Encode(userDTOs)
}

func (h *UserHandler) GetNewUsers(w http.ResponseWriter, r *http.Request) {
	layout := "02-01-2006" // DD-MM-YYYY

	fromParam := r.URL.Query().Get("from")
	from, err := time.Parse(layout, fromParam)
	if err != nil {
		http.Error(w, "Invalid <<from>> date", http.StatusBadRequest)
		return
	}

	toParam := r.URL.Query().Get("to")
	to, err := time.Parse(layout, toParam)
	if err != nil {
		http.Error(w, "Invalid <<to>> date", http.StatusBadRequest)
		return
	}

	users, err := h.userUseCase.GetNewUsers(r.Context(), from, to)
	if err != nil {
		http.Error(w, "New users not found", http.StatusNotFound)
		return
	}

	userDTOs := dto.ToUserDTOs(users)

	json.NewEncoder(w).Encode(userDTOs)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var input dto.UpdateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	idParam := r.URL.Query().Get("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user := &entity.User{
		ID: id,
	}
	if input.Username != nil {
		user.Username = *input.Username
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	user.UpdatedAt = time.Now()

	err = h.userUseCase.UpdateUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.UserDTO{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.userUseCase.DeleteUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}

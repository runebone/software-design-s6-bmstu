package dto

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Card struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	ColumnID    uuid.UUID `json:"column_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Position    float64   `json:"position"`
	CreatedAt   time.Time `json:"created_at"`
}

type Board struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Title  string    `json:"title"`
}

type Column struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	BoardID  uuid.UUID `json:"board_id"`
	Title    string    `json:"title"`
	Position float64   `json:"position"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

type CreateBoardRequest struct {
	Title string `json:"title"`
}

type CreateColumnRequest struct {
	BoardID uuid.UUID `json:"board_id"`
	Title   string    `json:"title"`
}

type CreateCardRequest struct {
	ColumnID    uuid.UUID `json:"column_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
}

type UpdateBoardRequest struct {
	CreateBoardRequest
}

type UpdateColumnRequest struct {
	CreateColumnRequest
}

type UpdateCardRequest struct {
	CreateCardRequest
}

type UserBase struct {
	ID       uuid.UUID
	Username string
	Email    string
}

type CardBase struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
}

type NewUsersAndCardsStats struct {
	Date               time.Time  `json:"date"`
	Users              []UserBase `json:"users"`
	Cards              []CardBase `json:"cards"`
	NumCardsByNewUsers int        `json:"num_cards_by_new_users"`
}

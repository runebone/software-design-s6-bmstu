package dto

import (
	"aggregator/internal/entity"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func UserToEntity(userDTO *User) entity.User {
	return entity.User{
		ID:       userDTO.ID,
		Username: userDTO.Username,
		Email:    userDTO.Email,
	}
}

func ToUserEntities(userDTOs []User) []entity.User {
	users := make([]entity.User, len(userDTOs))
	for i, userDTO := range userDTOs {
		users[i] = UserToEntity(&userDTO)
	}
	return users
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

func CardToEntity(cardDTO *Card) entity.Card {
	return entity.Card{
		ID:          cardDTO.ID,
		UserID:      cardDTO.UserID,
		Title:       cardDTO.Title,
		Description: cardDTO.Description,
	}
}

func ToCardEntities(cardDTOs []Card) []entity.Card {
	cards := make([]entity.Card, len(cardDTOs))
	for i, cardDTO := range cardDTOs {
		cards[i] = CardToEntity(&cardDTO)
	}
	return cards
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
	ID uuid.UUID `json:"id"`
	CreateBoardRequest
}

type UpdateColumnRequest struct {
	ID uuid.UUID `json:"id"`
	CreateColumnRequest
}

type UpdateCardRequest struct {
	ID uuid.UUID `json:"id"`
	CreateCardRequest
}

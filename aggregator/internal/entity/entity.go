package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Username string
	Email    string
}

type Card struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
}

// XXX: JSON hints should be in DTO, but I have to finish this stuff ASAP
type NewUsersAndCardsStats struct {
	Date               time.Time `json:"date"`
	Users              []User    `json:"users"`
	Cards              []Card    `json:"cards"`
	NumCardsByNewUsers int       `json:"num_cards_by_new_users"`
}

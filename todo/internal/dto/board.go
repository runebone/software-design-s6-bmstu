package dto

import (
	"todo/internal/entity"

	"github.com/google/uuid"
)

type CreateBoardRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Title  string    `json:"title"`
}

type Board struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type UpdateBoardRequest struct {
	Board
}

func ToBoardDTO(board *entity.Board) Board {
	return Board{
		ID:    board.ID,
		Title: board.Title,
	}
}

func ToBoardDTOs(boards []entity.Board) []Board {
	boardDTOs := make([]Board, len(boards))
	for i, board := range boards {
		boardDTOs[i] = ToBoardDTO(&board)
	}
	return boardDTOs
}

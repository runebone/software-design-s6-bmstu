package v1

import (
	"context"
	"errors"
	"fmt"
	"time"
	"todo/internal/common/logger"
	"todo/internal/entity"
	"todo/internal/repository"
	"todo/internal/usecase"

	"github.com/google/uuid"
)

var (
	ErrBoardEmptyTitle        = errors.New("board should have a title")
	ErrBoardNoUserID          = errors.New("board should have a user id")
	ErrColumnEmptyTitle       = errors.New("column should have a title")
	ErrColumnNoUserID         = errors.New("column should have a user id")
	ErrColumnNoBoardID        = errors.New("column should have a board id")
	ErrColumnNegativePosition = errors.New("column cannot have a negative position")
	ErrNegativeLimitOrOffset  = errors.New("limit and offset cannot be negative")
	ErrZeroLimit              = errors.New("limit cannot be zero")
	ErrCardNoUserID           = errors.New("card should have a user id")
	ErrCardNoColumnID         = errors.New("card should have a column id")
	ErrCardNegativePosition   = errors.New("card cannot have a negative position")
	ErrCardEmptyTitle         = errors.New("card should have a title")
	ErrGetBoardByID           = errors.New("failed to get board by id")
	ErrGetBoardsByUser        = errors.New("failed to get boards by user")
	ErrUpdateBoard            = errors.New("failed to update board")
	ErrDeleteBoard            = errors.New("failed to delete board")
	ErrGetColumnByID          = errors.New("failed to get column by id")
	ErrGetColumnsByBoard      = errors.New("failed to get columns by board")
	ErrUpdateColumn           = errors.New("failed to update column")
	ErrDeleteColumn           = errors.New("failed to delete column")
	ErrGetCardByID            = errors.New("failed to get card by id")
	ErrGetCardsByColumn       = errors.New("failed to get cards by column")
	ErrUpdateCard             = errors.New("failed to update card")
	ErrDeleteCard             = errors.New("failed to delete card")
	ErrInvalidTimeRange       = errors.New("<<from>> cannot be greater than <<to>> date")
	ErrGetNewCards            = errors.New("failed to get new cards")
)

type todoUseCase struct {
	boardRepo  repository.BoardRepository
	columnRepo repository.ColumnRepository
	cardRepo   repository.CardRepository
	log        logger.Logger
}

func NewTodoUseCase(
	boardRepo repository.BoardRepository,
	columnRepo repository.ColumnRepository,
	cardRepo repository.CardRepository,
	log logger.Logger,
) usecase.TodoUseCase {
	return &todoUseCase{
		boardRepo:  boardRepo,
		columnRepo: columnRepo,
		cardRepo:   cardRepo,
		log:        log,
	}
}

func (uc *todoUseCase) CreateBoard(ctx context.Context, board *entity.Board) error {
	header := "CreateBoard: "

	uc.log.Info(ctx, header+"Usecase called; Validating board", "board", board)

	err := validateBoard(board)

	if err != nil {
		info := "Validation failed"
		uc.log.Info(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	board.ID = uuid.New()
	board.CreatedAt = time.Now()
	board.UpdatedAt = time.Now()

	uc.log.Info(ctx, header+"Successful validation; Assigned uuid to board", "uuid", board.ID)

	uc.log.Info(ctx, header+"Making request to board repo (CreateBoard)", "board", board)

	err = uc.boardRepo.CreateBoard(ctx, board)

	if err != nil {
		info := "Failed to create board"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Board successfully created")

	return nil
}

func validateBoard(board *entity.Board) error {
	if board.Title == "" {
		return ErrBoardEmptyTitle
	}

	if board.UserID == uuid.Nil {
		return ErrBoardNoUserID
	}

	return nil
}

func (uc *todoUseCase) GetBoardByID(ctx context.Context, id uuid.UUID) (*entity.Board, error) {
	header := "GetBoardByID: "

	uc.log.Info(ctx, header+"Usecase called; Making request to board repo (GetBoardByID)", "id", id)

	board, err := uc.boardRepo.GetBoardByID(ctx, id)

	if err != nil {
		info := "Failed to get board by id"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got board", "board", board)

	return board, nil
}

func (uc *todoUseCase) GetBoardsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.Board, error) {
	header := "GetBoardsByUser: "

	uc.log.Info(ctx, header+"Usecase called; Validating limit and offset", "userID", userID, "limit", limit, "offset", offset)

	err := validateLimitAndOffset(limit, offset)

	if err != nil {
		info := "Validation failed"
		uc.log.Info(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Making request to board repo (GetBoardsByUser)", "userID", userID, "limit", limit, "offset", offset)

	boards, err := uc.boardRepo.GetBoardsByUser(ctx, userID, limit, offset)

	if err != nil {
		info := "Failed to get boards by user"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got boards", "boards", boards)

	return boards, nil
}

func (uc *todoUseCase) UpdateBoard(ctx context.Context, board *entity.Board) error {
	header := "UpdateBoard: "

	uc.log.Info(ctx, header+"Usecase called; Validating board", "board", board)

	err := validateBoard(board)
	if err == ErrBoardNoUserID {
		err = nil
	}

	if err != nil {
		info := "Validation failed"
		uc.log.Info(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	board.UpdatedAt = time.Now()

	uc.log.Info(ctx, header+"Making request to board repo (UpdateBoard)", "board", board)

	err = uc.boardRepo.UpdateBoard(ctx, board)

	if err != nil {
		info := "Failed to update board"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Board successfully updated")

	return nil
}

func (uc *todoUseCase) DeleteBoard(ctx context.Context, id uuid.UUID) error {
	header := "DeleteBoard: "

	uc.log.Info(ctx, header+"Usecase called; Making request to board repo (DeleteBoard)", "id", id)

	err := uc.boardRepo.DeleteBoard(ctx, id)

	if err != nil {
		info := "Failed to delete board"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successfully deleted board")

	return nil
}

func (uc *todoUseCase) CreateColumn(ctx context.Context, column *entity.Column) error {
	header := "CreateColumn: "

	uc.log.Info(ctx, header+"Usecase called; Validating column", "column", column)

	err := validateColumn(column)

	if err != nil {
		info := "Validation failed"
		uc.log.Info(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	column.ID = uuid.New()
	column.CreatedAt = time.Now()
	column.UpdatedAt = time.Now()

	uc.log.Info(ctx, header+"Successful validation; assigned uuid to column", "uuid", column.ID)

	uc.log.Info(ctx, header+"Making request to column repo (CreateColumn)", "column", column)

	err = uc.columnRepo.CreateColumn(ctx, column)

	if err != nil {
		info := "Failed to make request to repo"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Column successfully created")

	return nil
}

func validateColumn(column *entity.Column) error {
	if column.Title == "" {
		return ErrColumnEmptyTitle
	}

	if column.UserID == uuid.Nil {
		return ErrColumnNoUserID
	}

	if column.BoardID == uuid.Nil {
		return ErrColumnNoBoardID
	}

	if column.Position < 0 {
		return ErrColumnNegativePosition
	}

	return nil
}

func (uc *todoUseCase) GetColumnByID(ctx context.Context, id uuid.UUID) (*entity.Column, error) {
	header := "GetColumnByID: "

	uc.log.Info(ctx, header+"Usecase called; Making request to column repo (GetColumnByID)", "id", id)

	column, err := uc.columnRepo.GetColumnByID(ctx, id)

	if err != nil {
		info := "Failed to get column by id"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got column", "column", column)

	return column, nil
}

func (uc *todoUseCase) GetColumnsByBoard(ctx context.Context, boardID uuid.UUID, limit, offset int) ([]entity.Column, error) {
	header := "GetColumnsByBoard: "

	uc.log.Info(ctx, header+"Usecase called; Validating limit and offset", "boardID", boardID, "limit", limit, "offset", offset)

	err := validateLimitAndOffset(limit, offset)

	if err != nil {
		info := "Validation failed"
		uc.log.Info(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Validation successful; Making request to column repo (GetColumnsByBoard)", "boardID", boardID, "limit", limit, "offset", offset)

	columns, err := uc.columnRepo.GetColumnsByBoard(ctx, boardID, limit, offset)

	if err != nil {
		info := "Failed to get columns by board"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got columns", "columns", columns)

	return columns, nil
}

func validateLimitAndOffset(limit, offset int) error {
	if limit < 0 || offset < 0 {
		return ErrNegativeLimitOrOffset
	}

	if limit == 0 {
		return ErrZeroLimit
	}

	return nil
}

func (uc *todoUseCase) UpdateColumn(ctx context.Context, column *entity.Column) error {
	header := "UpdateColumn: "

	uc.log.Info(ctx, header+"Usecase called; Validating column", "column", column)

	err := validateColumn(column)
	if err == ErrColumnNoUserID || err == ErrColumnNoBoardID {
		err = nil
	}

	if err != nil {
		info := "Validation failed"
		uc.log.Info(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	column.UpdatedAt = time.Now()

	uc.log.Info(ctx, header+"Validation successful; Making request to column repo (UpdateColumn)", "column", column)

	err = uc.columnRepo.UpdateColumn(ctx, column)

	if err != nil {
		info := "Failed to update column"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Column successfully updated")

	return nil
}

func (uc *todoUseCase) DeleteColumn(ctx context.Context, id uuid.UUID) error {
	header := "DeleteColumn: "

	uc.log.Info(ctx, header+"Usecase called; Making request to column repo (DeleteColumn)", "id", id)

	err := uc.columnRepo.DeleteColumn(ctx, id)

	if err != nil {
		info := "Failed to delete column"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successfully deleted column")

	return nil
}

func (uc *todoUseCase) CreateCard(ctx context.Context, card *entity.Card) error {
	header := "CreateCard: "

	uc.log.Info(ctx, header+"Usecase called; Validating card", "card", card)

	err := validateCard(card)

	if err != nil {
		info := "Validation failed"
		uc.log.Info(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	card.ID = uuid.New()
	card.CreatedAt = time.Now()
	card.UpdatedAt = time.Now()

	uc.log.Info(ctx, header+"Successful validation; Assigned uuid to card", "uuid", card.ID)

	uc.log.Info(ctx, header+"Making request to card repo (CreateCard)", "card", card)

	err = uc.cardRepo.CreateCard(ctx, card)

	if err != nil {
		info := "Failed to create card"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	return nil
}

func validateCard(card *entity.Card) error {
	if card.UserID == uuid.Nil {
		return ErrCardNoUserID
	}

	if card.ColumnID == uuid.Nil {
		return ErrCardNoColumnID
	}

	if card.Position < 0 {
		return ErrCardNegativePosition
	}

	if card.Title == "" {
		return ErrCardEmptyTitle
	}

	return nil
}

func (uc *todoUseCase) GetCardByID(ctx context.Context, id uuid.UUID) (*entity.Card, error) {
	header := "GetCardByID: "

	uc.log.Info(ctx, header+"Usecase called; Making request to card repo (GetCardByID)", "id", id)

	card, err := uc.cardRepo.GetCardByID(ctx, id)

	if err != nil {
		info := "Failed to get card by id"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got card", "card", card)

	return card, nil
}

func (uc *todoUseCase) GetCardsByColumn(ctx context.Context, columnID uuid.UUID, limit, offset int) ([]entity.Card, error) {
	header := "GetCardsByColumn: "

	uc.log.Info(ctx, header+"Usecase called; Validating limit and offset", "columnID", columnID, "limit", limit, "offset", offset)

	err := validateLimitAndOffset(limit, offset)

	if err != nil {
		info := "Validation failed"
		uc.log.Info(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successful validation; Making request to card repo (GetCardsByColumn)", "columnID", columnID, "limit", limit, "offset", offset)

	cards, err := uc.cardRepo.GetCardsByColumn(ctx, columnID, limit, offset)

	if err != nil {
		info := "Failed to get cards by column"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got cards", "cards", cards)

	return cards, nil
}

func (uc *todoUseCase) GetNewCards(ctx context.Context, from, to time.Time) ([]entity.Card, error) {
	header := "GetNewCards: "

	uc.log.Info(ctx, header+"Usecase called; Validating from, to dates", "from", from, "to", to)

	err := validateFromToDate(from, to)

	if err != nil {
		info := "Validation failed"
		uc.log.Info(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successful validation; Making request to card repo (GetNewCards)", "from", from, "to", to)

	cards, err := uc.cardRepo.GetNewCards(ctx, from, to)

	if err != nil {
		info := "Failed get new cards"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got cards", "cards", cards)

	return cards, nil
}

func validateFromToDate(from, to time.Time) error {
	if from.Unix() > to.Unix() {
		return ErrInvalidTimeRange
	}

	return nil
}

func (uc *todoUseCase) UpdateCard(ctx context.Context, card *entity.Card) error {
	header := "UpdateCard: "

	uc.log.Info(ctx, header+"Usecase called; Validating card", "card", card)

	err := validateCard(card)
	if err == ErrCardNoColumnID || err == ErrCardNoUserID {
		err = nil
	}

	if err != nil {
		info := "Validation failed"
		uc.log.Info(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	card.UpdatedAt = time.Now()

	uc.log.Info(ctx, header+"Successful validation; Making request to card repo (UpdateCard)", "card", card)

	if card.ColumnID == uuid.Nil {
		err = uc.cardRepo.UpdateCard(ctx, card)
	} else {
		err = uc.cardRepo.MoveCard(ctx, card)
	}

	if err != nil {
		info := "Failed to update card"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Card successfully updated")

	return nil
}

func (uc *todoUseCase) DeleteCard(ctx context.Context, id uuid.UUID) error {
	header := "DeleteCard: "

	uc.log.Info(ctx, header+"Usecase called; Making request to card repo (DeleteCard)", "id", id)

	err := uc.cardRepo.DeleteCard(ctx, id)

	if err != nil {
		info := "Failed to delete card"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Card successfully deleted")

	return nil
}

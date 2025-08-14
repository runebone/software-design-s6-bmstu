package v1

import (
	"aggregator/internal/common/logger"
	"aggregator/internal/dto"
	"aggregator/internal/entity"
	"aggregator/internal/service/auth"
	"aggregator/internal/service/todo"
	"aggregator/internal/service/user"
	"aggregator/internal/usecase"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	dateLayout          string = "2006-01-02"
	ErrInvalidTimeRange error  = errors.New("<<from>> should be not greater than <<to>>")
	ErrGetNewUsers      error  = errors.New("failed to get new users")
	ErrGetNewCards      error  = errors.New("failed to get new cards")
	ErrRegister         error  = errors.New("failed to register")
	ErrLogin            error  = errors.New("failed to login")
	ErrRefresh          error  = errors.New("failed to refresh")
	ErrValidate         error  = errors.New("failed to validate token")
	ErrLogout           error  = errors.New("failed to logout")
	ErrGetBoards        error  = errors.New("failed to get boards")
	ErrGetColumns       error  = errors.New("failed to get columns")
	ErrGetCards         error  = errors.New("failed to get cards")
	ErrGetCard          error  = errors.New("failed to get card")
	ErrCreateBoard      error  = errors.New("failed to create board")
	ErrCreateColumn     error  = errors.New("failed to create column")
	ErrCreateCard       error  = errors.New("failed to create card")
	ErrUpdateBoard      error  = errors.New("failed to update board")
	ErrUpdateColumn     error  = errors.New("failed to update column")
	ErrUpdateCard       error  = errors.New("failed to update card")
	ErrDeleteBoard      error  = errors.New("failed to delete board")
	ErrDeleteColumn     error  = errors.New("failed to delete column")
	ErrDeleteCard       error  = errors.New("failed to delete card")
)

type AggregatorUseCase struct {
	userSvc user.UserService
	authSvc auth.AuthService
	todoSvc todo.TodoService
	log     logger.Logger
}

func NewAggregatorUseCase(
	userSvc user.UserService,
	authSvc auth.AuthService,
	todoSvc todo.TodoService,
	log logger.Logger,
) usecase.AggregatorUseCase {
	return &AggregatorUseCase{
		userSvc: userSvc,
		authSvc: authSvc,
		todoSvc: todoSvc,
		log:     log,
	}
}

func (uc *AggregatorUseCase) GetStats(ctx context.Context, from, to time.Time) ([]entity.NewUsersAndCardsStats, error) {
	header := "GetStats: "

	uc.log.Info(ctx, header+"Usecase called; Validating from, to dates", "from", from, "to", to)

	var err error

	if from.After(to) {
		err = ErrInvalidTimeRange
	}

	if err != nil {
		info := "Validation failed"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successful validation; Making request to user service (GetNewUsers)", "from", from, "to", to)

	users, err := uc.userSvc.GetNewUsers(ctx, from, to)
	if err != nil {
		info := "Failed to get new users"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got users", "users", users)

	dateUsersMap, newUsersMap := getDateUsersAndNewUsersMap(users)

	uc.log.Info(ctx, header+"Making request to todo service (GetNewCards)", "from", from, "to", to)

	cards, err := uc.todoSvc.GetNewCards(ctx, from, to)
	if err != nil {
		info := "Failed to get new cards"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got cards", "cards", cards)

	dateCardsMap, numCardsByNewUsersMap := getDateCardsAndNumCardsByNewUsersMap(cards, newUsersMap)

	stats := getStats(dateUsersMap, dateCardsMap, numCardsByNewUsersMap)

	uc.log.Info(ctx, header+"Success; Generated stats", "stats", stats)

	return stats, nil
}

func getDateUsersAndNewUsersMap(users []dto.User) (map[string][]entity.User, map[uuid.UUID]dto.User) {
	dateUsersMap := map[string][]entity.User{}
	newUsersMap := map[uuid.UUID]dto.User{}

	for _, userDTO := range users {
		dateKey := userDTO.CreatedAt.Format(dateLayout)
		user := dto.UserToEntity(&userDTO)

		dateUsersMap[dateKey] = append(dateUsersMap[dateKey], user)

		newUsersMap[user.ID] = userDTO
	}

	return dateUsersMap, newUsersMap
}

func getDateCardsAndNumCardsByNewUsersMap(cards []dto.Card, newUsersMap map[uuid.UUID]dto.User) (map[string][]entity.Card, map[string]int) {
	dateCardsMap := map[string][]entity.Card{}
	numCardsByNewUsersMap := map[string]int{}

	for _, cardDTO := range cards {
		dateKey := cardDTO.CreatedAt.Format(dateLayout)
		card := dto.CardToEntity(&cardDTO)

		dateCardsMap[dateKey] = append(dateCardsMap[dateKey], card)

		if user, ok := newUsersMap[card.UserID]; ok && user.CreatedAt.Format(dateLayout) == cardDTO.CreatedAt.Format(dateLayout) {
			numCardsByNewUsersMap[dateKey]++
		}
	}

	return dateCardsMap, numCardsByNewUsersMap
}

func getStats(dateUsersMap map[string][]entity.User, dateCardsMap map[string][]entity.Card, numCardsByNewUsersMap map[string]int) []entity.NewUsersAndCardsStats {
	stats := []entity.NewUsersAndCardsStats{}

	dates := mergeKeys(dateUsersMap, dateCardsMap)

	for dateKey := range dates {
		date, _ := time.Parse(dateLayout, dateKey)

		stat := entity.NewUsersAndCardsStats{
			Date:               date,
			Users:              dateUsersMap[dateKey],
			Cards:              dateCardsMap[dateKey],
			NumCardsByNewUsers: numCardsByNewUsersMap[dateKey],
		}

		stats = append(stats, stat)
	}

	return stats
}

func mergeKeys(dateUsersMap map[string][]entity.User, dateCardsMap map[string][]entity.Card) map[string]bool {
	dates := map[string]bool{}

	for k := range dateUsersMap {
		dates[k] = true
	}

	for k := range dateCardsMap {
		dates[k] = true
	}

	return dates
}

func (uc *AggregatorUseCase) Register(ctx context.Context, username, email, password string) (*dto.Tokens, error) {
	header := "Register: "

	uc.log.Info(ctx, header+"Usecase called; Making request to auth service", "username", username, "email", email, "password", password)

	tokens, err := uc.authSvc.Register(ctx, username, email, password)

	if err != nil {
		info := "Failed to register"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successful register; Got tokens", "tokens", tokens)

	return tokens, nil
}

func (uc *AggregatorUseCase) Login(ctx context.Context, email, password string) (*dto.Tokens, error) {
	header := "Login: "

	uc.log.Info(ctx, header+"Usecase called; Making request to auth service", "email", email, "password", password)

	tokens, err := uc.authSvc.Login(ctx, email, password)

	if err != nil {
		info := "Failed to login"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successful login; Got tokens", "tokens", tokens)

	return tokens, nil
}

func (uc *AggregatorUseCase) Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	header := "Refresh: "

	uc.log.Info(ctx, header+"Usecase called; Making request to auth service", "refreshToken", refreshToken)

	resp, err := uc.authSvc.Refresh(ctx, refreshToken)

	if err != nil {
		info := "Failed to refresh"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successful refresh", "refreshResponse", resp)

	return resp, nil
}

func (uc *AggregatorUseCase) Validate(ctx context.Context, token string) (*dto.ValidateTokenResponse, error) {
	header := "Validate: "

	uc.log.Info(ctx, header+"Usecase called; Making request to auth service (ValidateToken)", "token", token)

	resp, err := uc.authSvc.ValidateToken(ctx, token)

	if err != nil {
		info := "Failed to validate token"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successful validation", "validateTokenResponse", resp)

	return resp, nil
}

func (uc *AggregatorUseCase) Logout(ctx context.Context, refreshToken string) error {
	header := "Logout: "

	uc.log.Info(ctx, header+"Usecase called; Making request to auth service", "refreshToken", refreshToken)

	err := uc.authSvc.Logout(ctx, refreshToken)

	if err != nil {
		info := "Failed to logout"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successful logout")

	return nil
}

func (uc *AggregatorUseCase) GetBoards(ctx context.Context, userID string) ([]dto.Board, error) {
	header := "GetBoards: "

	uc.log.Info(ctx, header+"Usecase called; Making request to todo service", "userID", userID)

	boards, err := uc.todoSvc.GetBoards(ctx, userID)

	if err != nil {
		info := "Failed to get boards"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got boards", "boards", boards)

	return boards, nil
}

func (uc *AggregatorUseCase) GetColumns(ctx context.Context, boardID string) ([]dto.Column, error) {
	header := "GetColumns: "

	uc.log.Info(ctx, header+"Usecase called; Making request to todo service", "boardID", boardID)

	columns, err := uc.todoSvc.GetColumns(ctx, boardID)

	if err != nil {
		info := "Failed to get columns"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got columns", "columns", columns)

	return columns, nil
}

func (uc *AggregatorUseCase) GetCards(ctx context.Context, columnID string) ([]dto.Card, error) {
	header := "GetCards: "

	uc.log.Info(ctx, header+"Usecase called; Making request to todo service", "columnID", columnID)

	cards, err := uc.todoSvc.GetCards(ctx, columnID)

	if err != nil {
		info := "Failed to get cards"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got cards", "cards", cards)

	return cards, nil
}

func (uc *AggregatorUseCase) GetCard(ctx context.Context, id string) (*dto.Card, error) {
	header := "GetCard: "

	uc.log.Info(ctx, header+"Usecase called; Making request to todo service", "id", id)

	card, err := uc.todoSvc.GetCard(ctx, id)

	if err != nil {
		info := "Failed to get card"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got card", "card", card)

	return card, nil
}

func (uc *AggregatorUseCase) CreateBoard(ctx context.Context, board dto.Board) error {
	header := "CreateBoard: "

	uc.log.Info(ctx, header+"Usecase called; Making request to todo service", "board", board)

	err := uc.todoSvc.CreateBoard(ctx, board)

	if err != nil {
		info := "Failed to create board"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successfully created board")

	return nil
}

func (uc *AggregatorUseCase) CreateColumn(ctx context.Context, column dto.Column) error {
	header := "CreateColumn: "

	uc.log.Info(ctx, header+"Usecase called; Making request to todo service", "column", column)

	err := uc.todoSvc.CreateColumn(ctx, column)

	if err != nil {
		info := "Failed to create column"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successfully created column")

	return nil
}

func (uc *AggregatorUseCase) CreateCard(ctx context.Context, card dto.Card) error {
	header := "CreateCard: "

	uc.log.Info(ctx, header+"Usecase called; Making request to todo service", "card", card)

	err := uc.todoSvc.CreateCard(ctx, card)

	if err != nil {
		info := "Failed to create card"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successfully created card")

	return nil
}

func (uc *AggregatorUseCase) UpdateBoard(ctx context.Context, board *dto.Board) error {
	header := "UpdateBoard: "

	uc.log.Info(ctx, header+"Usecase called; Making request to todo service", "board", board)

	err := uc.todoSvc.UpdateBoard(ctx, board)

	if err != nil {
		info := "Failed to update board"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successfully updated board")

	return nil
}

func (uc *AggregatorUseCase) UpdateColumn(ctx context.Context, column *dto.Column) error {
	header := "UpdateColumn: "

	uc.log.Info(ctx, header+"Usecase called; Making request to todo service", "column", column)

	err := uc.todoSvc.UpdateColumn(ctx, column)

	if err != nil {
		info := "Failed to update column"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successfully updated column")

	return nil
}

func (uc *AggregatorUseCase) UpdateCard(ctx context.Context, card *dto.Card) error {
	header := "UpdateCard: "

	uc.log.Info(ctx, header+"Usecase called; Making request to todo service", "card", card)

	err := uc.todoSvc.UpdateCard(ctx, card)

	if err != nil {
		info := "Failed to update card"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successfully updated card")

	return nil
}

func (uc *AggregatorUseCase) DeleteBoard(ctx context.Context, id string) error {
	header := "DeleteBoard: "

	uc.log.Info(ctx, header+"Usecase called; Making request to todo service", "id", id)

	err := uc.todoSvc.DeleteBoard(ctx, id)

	if err != nil {
		info := "Failed to delete board"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successfully deleted board")

	return nil
}

func (uc *AggregatorUseCase) DeleteColumn(ctx context.Context, id string) error {
	header := "DeleteColumn: "

	uc.log.Info(ctx, header+"Usecase called; Making request to todo service", "id", id)

	err := uc.todoSvc.DeleteColumn(ctx, id)

	if err != nil {
		info := "Failed to delete column"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successfully deleted column")

	return nil
}

func (uc *AggregatorUseCase) DeleteCard(ctx context.Context, id string) error {
	header := "DeleteCard: "

	uc.log.Info(ctx, header+"Usecase called; Making request to todo service", "id", id)

	err := uc.todoSvc.DeleteCard(ctx, id)

	if err != nil {
		info := "Failed to delete card"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successfully deleted card")

	return nil
}

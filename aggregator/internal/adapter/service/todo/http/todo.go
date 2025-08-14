package http

import (
	"aggregator/internal/common/logger"
	"aggregator/internal/dto"
	"aggregator/internal/service/todo"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const layout string = "02-01-2006"

var (
	ErrDecodeResponse func(error) error = func(err error) error {
		return fmt.Errorf("Failed to decode response: %w", err)
	}
	ErrGetNewCards  error = errors.New("failed to get new cards")
	ErrGetBoards    error = errors.New("failed to get boards")
	ErrGetColumns   error = errors.New("failed to get columns")
	ErrGetCards     error = errors.New("failed to get cards")
	ErrGetCard      error = errors.New("failed to get card")
	ErrCreateBoard  error = errors.New("failed to create board")
	ErrCreateColumn error = errors.New("failed to create column")
	ErrCreateCard   error = errors.New("failed to create card")
	ErrUpdateBoard  error = errors.New("failed to update board")
	ErrUpdateColumn error = errors.New("failed to update column")
	ErrUpdateCard   error = errors.New("failed to update card")
	ErrDeleteBoard  error = errors.New("failed to delete board")
	ErrDeleteColumn error = errors.New("failed to delete column")
	ErrDeleteCard   error = errors.New("failed to delete card")
)

type TodoService struct {
	baseURL    string
	httpClient *http.Client
	log        logger.Logger
}

func NewTodoService(baseURL string, timeout time.Duration, logger logger.Logger) todo.TodoService {
	return &TodoService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		log: logger,
	}
}

func (s *TodoService) GetNewCards(ctx context.Context, from, to time.Time) ([]dto.Card, error) {
	fromStr := from.Format(layout)
	toStr := to.Format(layout)
	url := fmt.Sprintf("%s/cards/new?from=%s&to=%s", s.baseURL, fromStr, toStr)

	method := http.MethodGet
	resp, err := s.makeRequest(ctx, method, url, nil)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = ErrGetNewCards
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	var cards []dto.Card
	if err := json.NewDecoder(resp.Body).Decode(&cards); err != nil {
		err = ErrDecodeResponse(err)
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	return cards, nil
}

func (s *TodoService) GetBoards(ctx context.Context, userID string) ([]dto.Board, error) {
	url := fmt.Sprintf("%s/boards?user_id=%s", s.baseURL, userID)

	method := http.MethodGet
	resp, err := s.makeRequest(ctx, method, url, nil)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = ErrGetBoards
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	var boards []dto.Board
	if err := json.NewDecoder(resp.Body).Decode(&boards); err != nil {
		err = ErrDecodeResponse(err)
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	return boards, nil
}

func (s *TodoService) GetColumns(ctx context.Context, boardID string) ([]dto.Column, error) {
	url := fmt.Sprintf("%s/columns?board_id=%s", s.baseURL, boardID)

	method := http.MethodGet
	resp, err := s.makeRequest(ctx, method, url, nil)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url)
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = ErrGetColumns
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	var columns []dto.Column
	if err := json.NewDecoder(resp.Body).Decode(&columns); err != nil {
		err = ErrDecodeResponse(err)
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	return columns, nil
}

func (s *TodoService) GetCards(ctx context.Context, columnID string) ([]dto.Card, error) {
	url := fmt.Sprintf("%s/cards?column_id=%s", s.baseURL, columnID)

	method := http.MethodGet
	resp, err := s.makeRequest(ctx, method, url, nil)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrGetCards
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	var cards []dto.Card
	if err := json.NewDecoder(resp.Body).Decode(&cards); err != nil {
		err = ErrDecodeResponse(err)
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	return cards, nil
}

func (s *TodoService) GetCard(ctx context.Context, id string) (*dto.Card, error) {
	url := fmt.Sprintf("%s/cards/%s", s.baseURL, id)

	method := http.MethodGet
	resp, err := s.makeRequest(ctx, method, url, nil)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrGetCard
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	var card dto.Card
	if err := json.NewDecoder(resp.Body).Decode(&card); err != nil {
		err = ErrDecodeResponse(err)
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	return &card, nil
}

func (s *TodoService) CreateBoard(ctx context.Context, board dto.Board) error {
	url := fmt.Sprintf("%s/boards", s.baseURL)

	data := board

	method := http.MethodPost
	resp, err := s.makeRequest(ctx, method, url, data)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url, "data", data)
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		err = ErrCreateBoard
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

func (s *TodoService) CreateColumn(ctx context.Context, column dto.Column) error {
	url := fmt.Sprintf("%s/columns", s.baseURL)

	data := column

	method := http.MethodPost
	resp, err := s.makeRequest(ctx, method, url, data)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url, "data", data)
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		err = ErrCreateColumn
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

func (s *TodoService) CreateCard(ctx context.Context, card dto.Card) error {
	url := fmt.Sprintf("%s/cards", s.baseURL)

	data := card

	method := http.MethodPost
	resp, err := s.makeRequest(ctx, method, url, data)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url, "data", data)
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		err = ErrCreateCard
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

func (s *TodoService) UpdateBoard(ctx context.Context, board *dto.Board) error {
	url := fmt.Sprintf("%s/boards", s.baseURL)

	data := board

	method := http.MethodPut
	resp, err := s.makeRequest(ctx, method, url, data)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url, "data", data)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrUpdateBoard
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

func (s *TodoService) UpdateColumn(ctx context.Context, column *dto.Column) error {
	url := fmt.Sprintf("%s/columns", s.baseURL)

	data := column

	method := http.MethodPut
	resp, err := s.makeRequest(ctx, method, url, data)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url, "data", data)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrUpdateColumn
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

func (s *TodoService) UpdateCard(ctx context.Context, card *dto.Card) error {
	url := fmt.Sprintf("%s/cards", s.baseURL)

	data := card

	method := http.MethodPut
	resp, err := s.makeRequest(ctx, method, url, data)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url, "data", data)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrUpdateCard
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

func (s *TodoService) DeleteBoard(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/boards?id=%s", s.baseURL, id)

	method := http.MethodDelete
	resp, err := s.makeRequest(ctx, method, url, nil)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrDeleteBoard
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

func (s *TodoService) DeleteColumn(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/columns?id=%s", s.baseURL, id)

	method := http.MethodDelete
	resp, err := s.makeRequest(ctx, method, url, nil)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrDeleteColumn
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

func (s *TodoService) DeleteCard(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/cards?id=%s", s.baseURL, id)

	method := http.MethodDelete
	resp, err := s.makeRequest(ctx, method, url, nil)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrDeleteCard
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

func (s *TodoService) makeRequest(ctx context.Context, method, url string, data any) (*http.Response, error) {
	jsonBody, err := json.Marshal(data)
	if err != nil {
		err = fmt.Errorf("error marshaling user data: %w", err)
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		err = fmt.Errorf("error creating request: %w", err)
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("error sending request: %w", err)
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	return resp, nil
}

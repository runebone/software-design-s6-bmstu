package http

import (
	"bytes"
	"cli/internal/common/logger"
	"cli/internal/dto"
	"cli/internal/service"
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
	ErrUnauthorized error = errors.New("Unauthorized")
	ErrRegister     error = errors.New("User wasn't created")
	ErrLogin        error = errors.New("Failed to log in")
	ErrRefresh      error = errors.New("Failed to refresh")
	ErrValidate     error = errors.New("Failed to validate token")
	ErrLogout       error = errors.New("Failed to log out")
	ErrGetNewCards  error = errors.New("Failed to get new cards")
	ErrGetBoards    error = errors.New("Failed to get boards")
	ErrGetColumns   error = errors.New("Failed to get columns")
	ErrGetCards     error = errors.New("Failed to get cards")
	ErrGetCard      error = errors.New("Failed to get card")
	ErrCreateBoard  error = errors.New("Failed to create board")
	ErrCreateColumn error = errors.New("Failed to create column")
	ErrCreateCard   error = errors.New("Failed to create card")
	ErrUpdateBoard  error = errors.New("Failed to update board")
	ErrUpdateColumn error = errors.New("Failed to update column")
	ErrUpdateCard   error = errors.New("Failed to update card")
	ErrDeleteBoard  error = errors.New("Failed to delete board")
	ErrDeleteColumn error = errors.New("Failed to delete column")
	ErrDeleteCard   error = errors.New("Failed to delete card")
)

type AggregatorService struct {
	baseURL    string
	httpClient *http.Client
	log        logger.Logger
}

func NewAggregatorService(baseURL string, timeout time.Duration, logger logger.Logger) service.AggregatorService {
	return &AggregatorService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		log: logger,
	}
}

// Register(ctx context.Context, username, email, password string) (*dto.Tokens, error)
func (s *AggregatorService) Register(ctx context.Context, username, email, password string) (*dto.Tokens, error) {
	url := fmt.Sprintf("%s/register", s.baseURL)

	data := dto.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	s.log.Info(ctx, "Making register request", "url", url, "data", data)

	method := http.MethodPost
	resp, err := s.makeRequest(ctx, method, url, data)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url, "data", data)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		err = ErrRegister
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	var tokens dto.Tokens
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		err = ErrDecodeResponse(err)
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	return &tokens, nil
}

// Login(ctx context.Context, email, password string) (*dto.Tokens, error)
func (s *AggregatorService) Login(ctx context.Context, email, password string) (*dto.Tokens, error) {
	url := fmt.Sprintf("%s/login", s.baseURL)

	data := dto.LoginRequest{
		Email:    email,
		Password: password,
	}

	s.log.Info(ctx, "Making login request", "url", url, "data", data)

	method := http.MethodPost
	resp, err := s.makeRequest(ctx, method, url, data)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url, "data", data)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrLogin
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	var tokens dto.Tokens
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		err = ErrDecodeResponse(err)
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	return &tokens, nil
}

// Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
func (s *AggregatorService) Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	url := fmt.Sprintf("%s/refresh", s.baseURL)

	data := dto.RefreshRequest{
		RefreshToken: refreshToken,
	}

	s.log.Info(ctx, "Making refresh request", "url", url, "data", data)

	method := http.MethodPost
	resp, err := s.makeRequest(ctx, method, url, data)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url, "data", data)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrRefresh
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	var token dto.RefreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		err = ErrDecodeResponse(err)
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	return &token, nil
}

// Validate(ctx context.Context, token string) (*dto.ValidateTokenResponse, error)
func (s *AggregatorService) Validate(ctx context.Context, token string) (*dto.ValidateTokenResponse, error) {
	url := fmt.Sprintf("%s/validate", s.baseURL)

	data := dto.ValidateTokenRequest{
		Token: token,
	}

	s.log.Info(ctx, "Making validate request", "url", url, "data", data)

	method := http.MethodPost
	resp, err := s.makeRequest(ctx, method, url, data)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url, "data", data)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrValidate
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	var userData dto.ValidateTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		err = ErrDecodeResponse(err)
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	return &userData, nil
}

// Logout(ctx context.Context, refreshToken string) error
func (s *AggregatorService) Logout(ctx context.Context, refreshToken string) error {
	url := fmt.Sprintf("%s/logout", s.baseURL)

	data := dto.LogoutRequest{
		RefreshToken: refreshToken,
	}

	s.log.Info(ctx, "Making logout request", "url", url, "data", data)

	method := http.MethodPost
	resp, err := s.makeRequest(ctx, method, url, data)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url, "data", data)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = ErrLogout
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

// ShowBoards(ctx context.Context) ([]dto.Board, error)
func (s *AggregatorService) ShowBoards(ctx context.Context) ([]dto.Board, error) {
	url := fmt.Sprintf("%s/boards", s.baseURL)

	method := http.MethodGet
	resp, err := s.makeRequest(ctx, method, url, nil)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
		s.log.Error(ctx, err.Error())
		return nil, err
	}

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

// ShowBoard(ctx context.Context, boardID string) ([]dto.Column, error)
func (s *AggregatorService) ShowBoard(ctx context.Context, boardID string) ([]dto.Column, error) {
	url := fmt.Sprintf("%s/board/%s", s.baseURL, boardID)

	method := http.MethodGet
	resp, err := s.makeRequest(ctx, method, url, nil)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
		s.log.Error(ctx, err.Error())
		return nil, err
	}

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

// ShowColumn(ctx context.Context, columnID string) ([]dto.Card, error)
func (s *AggregatorService) ShowColumn(ctx context.Context, columnID string) ([]dto.Card, error) {
	url := fmt.Sprintf("%s/column/%s", s.baseURL, columnID)

	method := http.MethodGet
	resp, err := s.makeRequest(ctx, method, url, nil)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
		s.log.Error(ctx, err.Error())
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

// ShowCard(ctx context.Context, cardID string) (*dto.Card, error)
func (s *AggregatorService) ShowCard(ctx context.Context, cardID string) (*dto.Card, error) {
	url := fmt.Sprintf("%s/card/%s", s.baseURL, cardID)

	method := http.MethodGet
	resp, err := s.makeRequest(ctx, method, url, nil)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
		s.log.Error(ctx, err.Error())
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

// CreateBoard(ctx context.Context, board dto.Board) error
func (s *AggregatorService) CreateBoard(ctx context.Context, board dto.Board) error {
	url := fmt.Sprintf("%s/board", s.baseURL)

	data := board

	method := http.MethodPost
	resp, err := s.makeRequest(ctx, method, url, data)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url, "data", data)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
		s.log.Error(ctx, err.Error())
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		err = ErrCreateBoard
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

// CreateColumn(ctx context.Context, column dto.Column) error
func (s *AggregatorService) CreateColumn(ctx context.Context, column dto.Column) error {
	url := fmt.Sprintf("%s/column", s.baseURL)

	data := column

	method := http.MethodPost
	resp, err := s.makeRequest(ctx, method, url, data)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url, "data", data)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
		s.log.Error(ctx, err.Error())
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		err = ErrCreateBoard
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

// CreateCard(ctx context.Context, card dto.Card) error
func (s *AggregatorService) CreateCard(ctx context.Context, card dto.Card) error {
	url := fmt.Sprintf("%s/card", s.baseURL)

	data := card

	method := http.MethodPost
	resp, err := s.makeRequest(ctx, method, url, data)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url, "data", data)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
		s.log.Error(ctx, err.Error())
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		err = ErrCreateBoard
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

// UpdateBoard(ctx context.Context, board *dto.Board) error
func (s *AggregatorService) UpdateBoard(ctx context.Context, board *dto.Board) error {
	url := fmt.Sprintf("%s/board", s.baseURL)

	data := *board

	method := http.MethodPut
	resp, err := s.makeRequest(ctx, method, url, data)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url, "data", data)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
		s.log.Error(ctx, err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrUpdateBoard
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

// UpdateColumn(ctx context.Context, column *dto.Column) error
func (s *AggregatorService) UpdateColumn(ctx context.Context, column *dto.Column) error {
	url := fmt.Sprintf("%s/column", s.baseURL)

	data := *column

	method := http.MethodPut
	resp, err := s.makeRequest(ctx, method, url, data)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url, "data", data)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
		s.log.Error(ctx, err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrUpdateColumn
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

// UpdateCard(ctx context.Context, card *dto.Card) error
func (s *AggregatorService) UpdateCard(ctx context.Context, card *dto.Card) error {
	url := fmt.Sprintf("%s/card", s.baseURL)

	data := *card

	method := http.MethodPut
	resp, err := s.makeRequest(ctx, method, url, data)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url, "data", data)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
		s.log.Error(ctx, err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrUpdateCard
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

// DeleteBoard(ctx context.Context, id string) error
func (s *AggregatorService) DeleteBoard(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/board/%s", s.baseURL, id)

	method := http.MethodDelete
	resp, err := s.makeRequest(ctx, method, url, nil)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
		s.log.Error(ctx, err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrDeleteBoard
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

// DeleteColumn(ctx context.Context, id string) error
func (s *AggregatorService) DeleteColumn(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/column/%s", s.baseURL, id)

	method := http.MethodDelete
	resp, err := s.makeRequest(ctx, method, url, nil)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
		s.log.Error(ctx, err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrDeleteColumn
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

// DeleteCard(ctx context.Context, id string) error
func (s *AggregatorService) DeleteCard(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/card/%s", s.baseURL, id)

	method := http.MethodDelete
	resp, err := s.makeRequest(ctx, method, url, nil)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
		s.log.Error(ctx, err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrDeleteCard
		s.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

// Stats(ctx context.Context, from, to string) ([]dto.NewUsersAndCardsStats, error)
func (s *AggregatorService) Stats(ctx context.Context, from, to string) ([]dto.NewUsersAndCardsStats, error) {
	url := fmt.Sprintf("%s/stats/%s/%s", s.baseURL, from, to)

	method := http.MethodGet
	resp, err := s.makeRequest(ctx, method, url, nil)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrGetNewCards
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	var stats []dto.NewUsersAndCardsStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		err = ErrDecodeResponse(err)
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	return stats, nil
}

func (s *AggregatorService) makeRequest(ctx context.Context, method, url string, data any) (*http.Response, error) {
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

	tokens, ok := ctx.Value("tokens").(*dto.Tokens)
	if ok {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokens.AccessToken))
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

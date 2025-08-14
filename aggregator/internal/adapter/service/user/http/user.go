package http

import (
	"aggregator/internal/common/logger"
	"aggregator/internal/dto"
	"aggregator/internal/service/user"
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
	ErrGetNewUsers    error             = errors.New("failed to get new users")
	ErrDecodeResponse func(error) error = func(err error) error {
		return fmt.Errorf("Failed to decode response: %w", err)
	}
)

type UserService struct {
	baseURL    string
	httpClient *http.Client
	log        logger.Logger
}

func NewUserService(baseURL string, timeout time.Duration, logger logger.Logger) user.UserService {
	return &UserService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		log: logger,
	}
}

func (s *UserService) GetNewUsers(ctx context.Context, from, to time.Time) ([]dto.User, error) {
	fromStr := from.Format(layout)
	toStr := to.Format(layout)
	url := fmt.Sprintf("%s/users/new?from=%s&to=%s", s.baseURL, fromStr, toStr)

	s.log.Info(ctx, "Making GetNewUsers request", "url", url)

	method := http.MethodGet
	resp, err := s.makeRequest(ctx, method, url, nil)
	if err != nil {
		s.log.Error(ctx, "Error making the request", "method", method, "url", url)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = ErrGetNewUsers
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	var users []dto.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		err = ErrDecodeResponse(err)
		s.log.Error(ctx, err.Error())
		return nil, err
	}

	return users, nil
}

func (s *UserService) makeRequest(ctx context.Context, method, url string, data any) (*http.Response, error) {
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

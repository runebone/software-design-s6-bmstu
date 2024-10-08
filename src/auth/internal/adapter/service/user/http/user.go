package user

import (
	"auth/internal/dto"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type HTTPUserService struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPUserService(baseURL string, timeoutSec time.Duration) *HTTPUserService {
	return &HTTPUserService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeoutSec * 10e9,
		},
	}
}

func (s *HTTPUserService) GetUserByEmail(ctx context.Context, email string) (*dto.User, error) {
	url := fmt.Sprintf("%s/users?email=%s", s.baseURL, email)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user: status %d", resp.StatusCode)
	}

	var users []dto.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var user dto.User
	if len(users) == 0 {
		return nil, errors.New("couldn't find user by email")
	} else if len(users) > 1 {
		return nil, errors.New("several users with the same email found")
	}
	user = users[0]

	return &user, nil
}

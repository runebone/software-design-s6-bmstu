package http

import (
	"aggregator/internal/common/logger"
	"aggregator/internal/dto"
	"aggregator/internal/service/auth"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrRegister       error             = errors.New("User wasn't created")
	ErrLogin          error             = errors.New("Failed to log in")
	ErrRefresh        error             = errors.New("Failed to refresh")
	ErrValidate       error             = errors.New("Failed to validate token")
	ErrLogout         error             = errors.New("Failed to log out")
	ErrDecodeResponse func(error) error = func(err error) error {
		return fmt.Errorf("Failed to decode response: %w", err)
	}
)

type AuthService struct {
	baseURL    string
	httpClient *http.Client
	log        logger.Logger
}

func NewAuthService(baseURL string, timeout time.Duration, logger logger.Logger) auth.AuthService {
	return &AuthService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		log: logger,
	}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) (*dto.Tokens, error) {
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

func (s *AuthService) Login(ctx context.Context, email, password string) (*dto.Tokens, error) {
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

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
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

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*dto.ValidateTokenResponse, error) {
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

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
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

func (s *AuthService) makeRequest(ctx context.Context, method, url string, data any) (*http.Response, error) {
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

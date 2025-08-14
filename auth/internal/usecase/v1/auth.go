package v1

import (
	"auth/internal/common/logger"
	"auth/internal/dto"
	"auth/internal/entity"
	"auth/internal/repository"
	"auth/internal/service/tokengen"
	"auth/internal/service/user"
	"auth/internal/usecase"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrIncorrectPassword    error = errors.New("incorrect password")
	ErrGetAccessToken       error = errors.New("couldn't get access token")
	ErrGenerateRefreshToken error = errors.New("couldn't generate refresh token")
	ErrSaveRefreshToken     error = errors.New("couldn't save refresh token")
	ErrInvalidRefreshToken  error = errors.New("invalid refresh token")
	ErrGenerateAccessToken  error = errors.New("couldn't generate access token")
	ErrFindRefreshToken     error = errors.New("couldn't find refresh token")
	ErrDeleteRefreshToken   error = errors.New("couldn't delete refresh token")
	ErrCreateUser           error = errors.New("couldn't create user")
	ErrValidateToken        error = errors.New("couldn't validate token")
)

type authUseCase struct {
	tokenRepo    repository.TokenRepository
	userService  user.UserService
	tokenService tokengen.TokenService
	log          logger.Logger
}

func NewAuthUseCase(
	tokenRepo repository.TokenRepository,
	userService user.UserService,
	tokenService tokengen.TokenService,
	log logger.Logger,
) usecase.AuthUsecase {
	return &authUseCase{
		tokenRepo:    tokenRepo,
		userService:  userService,
		tokenService: tokenService,
		log:          log,
	}
}

func (uc *authUseCase) Register(ctx context.Context, username, email, password string) (*dto.Tokens, error) {
	header := "Register: "

	uc.log.Info(ctx, header+"Usecase called; Making request to user service (CreateUser)", "username", username, "email", email, "password", password)

	err := uc.userService.CreateUser(ctx, username, email, password)

	if err != nil {
		info := ErrCreateUser.Error()
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Making request to Login usecase", "email", email, "password", password)

	tokens, err := uc.Login(ctx, email, password)

	if err != nil {
		info := "Failed to login"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Logged in", "tokens", tokens)

	return tokens, nil
}

func (uc *authUseCase) Login(ctx context.Context, email, password string) (*dto.Tokens, error) {
	header := "Login: "

	uc.log.Info(ctx, header+"Usecase called; Making request to user service (GetUserByEmail)", "email", email, "password", password)

	user, err := uc.userService.GetUserByEmail(ctx, email)
	if err != nil {
		info := "Failed to get user by email"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got user", "user", user)

	// NOTE: Better to do this using something like uc.userService.VerifyPassword
	// and remove PasswordHash from User dto in user service
	if !validatePassword(password, user.PasswordHash) {
		err := ErrIncorrectPassword
		uc.log.Info(ctx, header+err.Error(), "password", password, "userPasswordHash", user.PasswordHash)
		return nil, err
	}

	userID := user.ID.String()
	role := user.Role

	uc.log.Info(ctx, header+"Making request to token service (GenerateAccessToken)", "userID", userID, "role", role)

	accessToken, err := uc.tokenService.GenerateAccessToken(ctx, userID, role)

	if err != nil {
		info := "Failed to generate access token"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got access token", "accessToken", accessToken)

	uc.log.Info(ctx, header+"Making request to token service (GenerateRefreshToken)", "userID", userID, "role", role)

	refreshToken, err := uc.tokenService.GenerateRefreshToken(ctx, userID, role)

	if err != nil {
		info := "Failed to generate refresh token"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got refresh token", "refreshToken", refreshToken)

	token := &entity.Token{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		CreatedAt: time.Now(),
	}

	uc.log.Info(ctx, header+"Making request to token repo (Save)", "refreshToken", token)

	err = uc.tokenRepo.Save(ctx, token)

	if err != nil {
		info := "Failed to save token"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	return &dto.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func validatePassword(password, passwordHash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) == nil
}

func (uc *authUseCase) Refresh(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponse, error) {
	header := "Refresh: "

	uc.log.Info(ctx, header+"Usecase called; Making request to token service (ValidateToken)", "refreshToken", refreshToken)

	userID, role, err := uc.tokenService.ValidateToken(ctx, refreshToken)
	if err != nil {
		info := "Validation failed"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got user data from token", "userID", userID, "role", role)

	uc.log.Info(ctx, header+"Making request to token service (GenerateAccessToken)", "userID", userID, "role", role)

	newAccessToken, err := uc.tokenService.GenerateAccessToken(ctx, userID, role)
	if err != nil {
		info := "Failed to generate access token"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return nil, fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got new accessToken", "accessToken", newAccessToken)

	return &dto.RefreshTokenResponse{
		AccessToken: newAccessToken,
	}, nil
}

func (uc *authUseCase) ValidateToken(ctx context.Context, token string) (string, string, error) {
	header := "ValidateToken: "

	uc.log.Info(ctx, header+"Usecase called; Making request to token service (ValidateToken)", "token", token)

	userID, role, err := uc.tokenService.ValidateToken(ctx, token)

	if err != nil {
		info := "Validation failed"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return "", "", fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Got user data", "userID", userID, "role", role)

	return userID, role, nil
}

func (uc *authUseCase) Logout(ctx context.Context, refreshToken string) error {
	header := "Logout: "

	uc.log.Info(ctx, header+"Usecase called; Making request to token repo (FindByToken)", "refreshToken", refreshToken)

	token, err := uc.tokenRepo.FindByToken(ctx, refreshToken)

	if err != nil {
		info := "Failed to find token id by token"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	tokenID := token.ID.String()

	uc.log.Info(ctx, header+"Found token id to delete", "tokenID", tokenID)

	uc.log.Info(ctx, header+"Making request to token repo (Delete)", "tokenID", tokenID)

	err = uc.tokenRepo.Delete(ctx, tokenID)

	if err != nil {
		info := "Failed to delete token"
		uc.log.Error(ctx, header+info, "err", err.Error())
		return fmt.Errorf(header+info+": %w", err)
	}

	uc.log.Info(ctx, header+"Successfully logged out (Deleted refresh token from DB)")

	return nil
}

package v1_test

import (
	"auth/internal/dto"
	"auth/internal/entity"
	"auth/internal/usecase"
	v1 "auth/internal/usecase/v1"
	"auth/mocks"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type testSetup struct {
	ctx           context.Context
	mockTokenRepo *mocks.TokenRepository
	mockUserSvc   *mocks.UserService
	mockTokenSvc  *mocks.TokenService
	uc            usecase.AuthUsecase
}

func setup() *testSetup {
	ctx := context.TODO()

	mockTokenRepo := new(mocks.TokenRepository)
	mockUserSvc := new(mocks.UserService)
	mockTokenSvc := new(mocks.TokenService)

	authUseCase := v1.NewAuthUseCase(mockTokenRepo, mockUserSvc, mockTokenSvc)

	return &testSetup{
		ctx:           ctx,
		mockTokenRepo: mockTokenRepo,
		mockUserSvc:   mockUserSvc,
		mockTokenSvc:  mockTokenSvc,
		uc:            authUseCase,
	}
}

func TestLogin(t *testing.T) {
	ts := setup()

	mockUserSvcFnOk := func(email string, user *dto.User) {
		ts.mockUserSvc.On("GetUserByEmail", ts.ctx, email).Return(user, nil)
	}
	// mockUserSvcFnErr := func(email string, user *dto.User) {
	// 	ts.mockUserSvc.On("GetUserByEmail", ts.ctx, email).Return(nil, errors.New(""))
	// }

	// Save(ctx context.Context, token *entity.Token) error
	mockTokenRepoSaveFnOk := func(token *entity.Token) {
		ts.mockTokenRepo.On("Save", ts.ctx, mock.Anything).Return(nil)
	}
	// mockTokenRepoSaveFnErr := func(token *entity.Token) {
	// 	ts.mockTokenRepo.On("Save", ts.ctx, token).Return(errors.New(""))
	// }

	// // Delete(ctx context.Context, tokenID string) error
	// mockTokenRepoDeleteFnOk := func(tokenID string) {
	// 	ts.mockTokenRepo.On("Delete", ts.ctx, tokenID).Return(nil)
	// }
	// mockTokenRepoDeleteFnErr := func(tokenID string) {
	// 	ts.mockTokenRepo.On("Delete", ts.ctx, tokenID).Return(errors.New(""))
	// }

	// // FindByToken(ctx context.Context, token string) (*entity.Token, error)
	// mockTokenRepoFindByTokenFnOk := func(token string, retToken *entity.Token) {
	// 	ts.mockTokenRepo.On("FindByToken", ts.ctx, token).Return(retToken, nil)
	// }
	// mockTokenRepoFindByTokenFnErr := func(token string, retToken *entity.Token) {
	// 	ts.mockTokenRepo.On("FindByToken", ts.ctx, token).Return(nil, errors.New(""))
	// }

	// GenerateAccessToken(ctx context.Context, userID string) (string, error)
	mockTokenSvcGenerateAccessTokenOk := func(userID, role, accessToken string) {
		ts.mockTokenSvc.On("GenerateAccessToken", ts.ctx, userID, role).Return(accessToken, nil)
	}
	// mockTokenSvcGenerateAccessTokenErr := func(userID, accessToken string) {
	// 	ts.mockTokenSvc.On("GenerateAccessToken", ts.ctx, userID).Return(nil, errors.New(""))
	// }

	// GenerateRefreshToken(ctx context.Context, userID string) (string, error)
	mockTokenSvcGenerateRefreshTokenOk := func(userID, role, refreshToken string) {
		ts.mockTokenSvc.On("GenerateRefreshToken", ts.ctx, userID, role).Return(refreshToken, nil)
	}
	// mockTokenSvcGenerateRefreshTokenErr := func(userID, accessToken string) {
	// 	ts.mockTokenSvc.On("GenerateRefreshToken", ts.ctx, userID).Return(nil, errors.New(""))
	// }

	// // ValidateToken(ctx context.Context, token string) (string, error) // Returns userID
	// mockTokenSvcValidateTokenOk := func(token, userID string) {
	// 	ts.mockTokenSvc.On("ValidateToken", ts.ctx, token).Return(userID, nil)
	// }
	// mockTokenSvcValidateTokenErr := func(token, userID string) {
	// 	ts.mockTokenSvc.On("ValidateToken", ts.ctx, token).Return(nil, errors.New(""))
	// }

	tests := []struct {
		name                             string
		userID                           uuid.UUID
		tokenID                          uuid.UUID
		token                            string
		username                         string
		email                            string
		role                             string
		password                         string
		createdAt                        time.Time
		responseFn                       func(accessToken, refreshToken string) *dto.Tokens
		accessToken                      string
		refreshToken                     string
		mockUserSvcFn                    func(email string, user *dto.User)
		mockTokenSvcGenerateAccessToken  func(userID, role, accessToken string)
		mockTokenSvcGenerateRefreshToken func(userID, role, refreshToken string)
		mockTokenRepoSaveFn              func(token *entity.Token)
		wantErr                          bool
		errMsg                           string
	}{
		{
			name:      "success",
			userID:    uuid.New(),
			tokenID:   uuid.New(),
			token:     "alskdjflksadjflkjdsf",
			username:  "username",
			email:     "success@email.com",
			role:      "user",
			createdAt: time.Now(),
			password:  "Pa$$w0rD",
			responseFn: func(accessToken string, refreshToken string) *dto.Tokens {
				return &dto.Tokens{
					AccessToken:  accessToken,
					RefreshToken: refreshToken,
				}
			},
			mockUserSvcFn:                    mockUserSvcFnOk,
			mockTokenSvcGenerateAccessToken:  mockTokenSvcGenerateAccessTokenOk,
			mockTokenSvcGenerateRefreshToken: mockTokenSvcGenerateRefreshTokenOk,
			mockTokenRepoSaveFn:              mockTokenRepoSaveFnOk,
			wantErr:                          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			passwordHash, _ := hashPassword(tt.password)
			user := &dto.User{
				ID:           tt.userID,
				Username:     tt.username,
				Email:        tt.email,
				Role:         tt.role,
				PasswordHash: passwordHash,
			}
			tt.mockUserSvcFn(tt.email, user)
			tt.mockTokenSvcGenerateAccessToken(tt.userID.String(), tt.role, tt.accessToken)
			tt.mockTokenSvcGenerateRefreshToken(tt.userID.String(), tt.role, tt.refreshToken)
			token := &entity.Token{
				ID:        tt.tokenID,
				UserID:    tt.userID,
				Token:     tt.token,
				CreatedAt: tt.createdAt,
			}
			tt.mockTokenRepoSaveFn(token)

			expectedLoginResponse := tt.responseFn(tt.accessToken, tt.refreshToken)
			actualLoginResponse, err := ts.uc.Login(ts.ctx, tt.email, tt.password)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, expectedLoginResponse, actualLoginResponse)
				// ts.mockUserSvc.AssertCalled(t, "GetNewUsers", ts.ctx, mock.Anything, mock.Anything, mock.Anything)
			}

			t.Cleanup(func() {
				ts.mockUserSvc.ExpectedCalls = nil
				ts.mockUserSvc.Calls = nil

				ts.mockTokenRepo.ExpectedCalls = nil
				ts.mockTokenRepo.Calls = nil

				ts.mockTokenSvc.ExpectedCalls = nil
				ts.mockTokenSvc.Calls = nil
			})
		})
	}
}

// XXX: Hashing logic copy-pasted from User microservice
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// func TestRefresh(t *testing.T) {
// 	ts := setup()
// }

// func TestLogout(t *testing.T) {
// 	ts := setup()
// }

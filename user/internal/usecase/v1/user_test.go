package v1_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
	"user/internal/entity"
	"user/internal/repository"
	"user/internal/usecase"
	v1 "user/internal/usecase/v1"
	"user/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testSetup struct {
	ctx         context.Context
	mockRepo    *mocks.UserRepository
	userUseCase usecase.UserUseCase
}

func setup() *testSetup {
	ctx := context.TODO()
	mockRepo := new(mocks.UserRepository)
	userUseCase := v1.NewUserUseCase(mockRepo)

	return &testSetup{
		ctx:         ctx,
		mockRepo:    mockRepo,
		userUseCase: userUseCase,
	}
}

func TestCreateUser(t *testing.T) {
	ts := setup()

	tests := []struct {
		name       string
		user       *entity.User
		mockRepoFn func()
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			user: &entity.User{
				Username:     "ValidUser",
				Email:        "valid@example.com",
				PasswordHash: "Password@123",
			},
			mockRepoFn: func() {
				ts.mockRepo.On("CreateUser", ts.ctx, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			user: &entity.User{
				Username:     "ValidUser",
				Email:        "invalid-email",
				PasswordHash: "Password@123",
			},
			mockRepoFn: func() {},
			wantErr:    true,
			errMsg:     "invalid email format",
		},
		{
			name: "invalid username - too short",
			user: &entity.User{
				Username:     "ba_D",
				Email:        "valid@example.com",
				PasswordHash: "Password@123",
			},
			mockRepoFn: func() {},
			wantErr:    true,
			errMsg:     "invalid username format",
		},
		{
			name: "invalid username - has special characters",
			user: &entity.User{
				Username:     "Usern@me",
				Email:        "valid@example.com",
				PasswordHash: "Password@123",
			},
			mockRepoFn: func() {},
			wantErr:    true,
			errMsg:     "invalid username format",
		},
		{
			name: "invalid password - too short",
			user: &entity.User{
				Username:     "ValidUser",
				Email:        "valid@example.com",
				PasswordHash: "p@Ss0",
			},
			mockRepoFn: func() {},
			wantErr:    true,
			errMsg:     "password must be between 8 and 32 characters long",
		},
		{
			name: "invalid password - too long",
			user: &entity.User{
				Username:     "ValidUser",
				Email:        "valid@example.com",
				PasswordHash: "p@ssw0rdUPopaBylaSobakaObYeyoLubilOnaSyelaKosokMyasaOnEeUbilINaMogileNapisalUPopaBylaSobakaOnYeyoLubil",
			},
			mockRepoFn: func() {},
			wantErr:    true,
			errMsg:     "password must be between 8 and 32 characters long",
		},
		{
			name: "invalid password - no uppercase letters",
			user: &entity.User{
				Username:     "ValidUser",
				Email:        "valid@example.com",
				PasswordHash: "p@ssw0rd",
			},
			mockRepoFn: func() {},
			wantErr:    true,
			errMsg:     "password must contain at least one uppercase letter",
		},
		{
			name: "invalid password - no lowercase letters",
			user: &entity.User{
				Username:     "ValidUser",
				Email:        "valid@example.com",
				PasswordHash: "P@SSW0RD",
			},
			mockRepoFn: func() {},
			wantErr:    true,
			errMsg:     "password must contain at least one lowercase letter",
		},
		{
			name: "invalid password - no digits",
			user: &entity.User{
				Username:     "ValidUser",
				Email:        "valid@example.com",
				PasswordHash: "P@ssworD",
			},
			mockRepoFn: func() {},
			wantErr:    true,
			errMsg:     "password must contain at least one digit",
		},
		{
			name: "invalid password - no special characters",
			user: &entity.User{
				Username:     "ValidUser",
				Email:        "valid@example.com",
				PasswordHash: "Passw0rD",
			},
			mockRepoFn: func() {},
			wantErr:    true,
			errMsg:     "password must contain at least one special character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn()

			err := ts.userUseCase.CreateUser(ts.ctx, *tt.user)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
				ts.mockRepo.AssertNotCalled(t, "CreateUser")
			} else {
				assert.Nil(t, err)
				ts.mockRepo.AssertCalled(t, "CreateUser", ts.ctx, mock.Anything)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	ts := setup()

	userID := uuid.New()
	user := &entity.User{
		ID:           userID,
		Username:     "ValidUser",
		Email:        "valid@example.com",
		PasswordHash: "P@ssw0rD",
	}

	tests := []struct {
		name       string
		userID     uuid.UUID
		user       *entity.User
		mockRepoFn func()
		wantErr    bool
		errMsg     string
	}{
		{
			name:   "success",
			userID: userID,
			user:   user,
			mockRepoFn: func() {
				ts.mockRepo.On("GetUserByID", ts.ctx, userID).Return(user, nil)
			},
			wantErr: false,
		},
		{
			name:   "fail",
			userID: uuid.New(),
			mockRepoFn: func() {
				ts.mockRepo.On("GetUserByID", ts.ctx, mock.Anything).Return(nil, errors.New(""))
			},
			wantErr: true,
			errMsg:  "failed to get user by id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn()

			user, err := ts.userUseCase.GetUserByID(ts.ctx, tt.userID)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
				ts.mockRepo.AssertNotCalled(t, "GetUserByID")
			} else {
				assert.Nil(t, err)
				assert.Equal(t, user, tt.user)
				ts.mockRepo.AssertCalled(t, "GetUserByID", ts.ctx, mock.Anything)
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	ts := setup()

	id := uuid.New()
	idStr := id.String()
	filterID := repository.UserFilter{
		ID: &idStr,
	}
	usersByID := []entity.User{
		{
			ID:       id,
			Username: "userByID",
			Email:    "email@byid.com",
		},
	}

	email := "email@example.com"
	filterEmail := repository.UserFilter{
		Email: &email,
	}
	usersByEmail := []entity.User{
		{
			ID:       uuid.New(),
			Username: "userByEmail",
			Email:    email,
		},
	}

	username := "username"
	filterUsername := repository.UserFilter{
		Username: &username,
	}
	usersByUsername := []entity.User{
		{
			ID:       uuid.New(),
			Username: username,
			Email:    "email@byusername.com",
		},
	}

	filterNoUsers := repository.UserFilter{
		ID:       &idStr,
		Email:    &email,
		Username: &username,
	}
	usersNoUsers := []entity.User{}

	tests := []struct {
		name       string
		filter     repository.UserFilter
		users      []entity.User
		mockRepoFn func()
		wantErr    bool
		errMsg     string
	}{
		{
			name:   "success - get users by id",
			filter: filterID,
			users:  usersByID,
			mockRepoFn: func() {
				ts.mockRepo.On("GetUsers", ts.ctx, filterID).Return(usersByID, nil)
			},
			wantErr: false,
		},
		{
			name:   "success - get users by email",
			filter: filterEmail,
			users:  usersByEmail,
			mockRepoFn: func() {
				ts.mockRepo.On("GetUsers", ts.ctx, filterEmail).Return(usersByEmail, nil)
			},
			wantErr: false,
		},
		{
			name:   "success - get users by username",
			filter: filterUsername,
			users:  usersByUsername,
			mockRepoFn: func() {
				ts.mockRepo.On("GetUsers", ts.ctx, filterUsername).Return(usersByUsername, nil)
			},
			wantErr: false,
		},
		{
			name:   "success - no users",
			filter: filterNoUsers,
			users:  usersNoUsers,
			mockRepoFn: func() {
				ts.mockRepo.On("GetUsers", ts.ctx, filterNoUsers).Return(usersNoUsers, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn()

			users, err := ts.userUseCase.GetUsers(ts.ctx, tt.filter)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
				ts.mockRepo.AssertNotCalled(t, "GetUsers")
			} else {
				assert.Nil(t, err)
				assert.Equal(t, users, tt.users)
				ts.mockRepo.AssertCalled(t, "GetUsers", ts.ctx, mock.Anything)
			}
		})
	}
}

func TestGetUsersBatch(t *testing.T) {
	ts := setup()

	users := []entity.User{
		{
			ID:       uuid.New(),
			Username: "Username1",
			Email:    "email1@test.com",
		},
		{
			ID:       uuid.New(),
			Username: "Username2",
			Email:    "email2@test.com",
		},
		{
			ID:       uuid.New(),
			Username: "Username3",
			Email:    "email3@test.com",
		},
		{
			ID:       uuid.New(),
			Username: "Username4",
			Email:    "email4@test.com",
		},
		{
			ID:       uuid.New(),
			Username: "Username5",
			Email:    "email5@test.com",
		},
	}

	mockRepoFnOk := func(limit, offset int, users []entity.User) {
		ts.mockRepo.On("GetUsersBatch", ts.ctx, limit, offset).Return(users, nil)
	}

	tests := []struct {
		name          string
		users         []entity.User
		limit, offset int
		mockRepoFn    func(int, int, []entity.User)
		wantErr       bool
		errMsg        string
	}{
		{
			name:       "success - get all users",
			users:      users,
			limit:      10,
			offset:     0,
			mockRepoFn: mockRepoFnOk,
			wantErr:    false,
		},
		{
			name:       "success - get [2:4] users slice (2nd and 3rd user)",
			users:      users[2:4],
			limit:      2,
			offset:     2,
			mockRepoFn: mockRepoFnOk,
			wantErr:    false,
		},
		{
			name:       "fail - negative limit",
			limit:      -1,
			offset:     0,
			mockRepoFn: func(int, int, []entity.User) {},
			wantErr:    true,
			errMsg:     "limit and offset can't be negative",
		},
		{
			name:       "fail - negative offset",
			limit:      1,
			offset:     -1,
			mockRepoFn: func(int, int, []entity.User) {},
			wantErr:    true,
			errMsg:     "limit and offset can't be negative",
		},
		{
			name:       "fail - zero limit",
			limit:      0,
			offset:     0,
			mockRepoFn: func(int, int, []entity.User) {},
			wantErr:    true,
			errMsg:     "limit should be greater than zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn(tt.limit, tt.offset, tt.users)

			users, err := ts.userUseCase.GetUsersBatch(ts.ctx, tt.limit, tt.offset)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
				ts.mockRepo.AssertNotCalled(t, "GetUsersBatch")
			} else {
				assert.Nil(t, err)
				assert.Equal(t, users, tt.users)
				ts.mockRepo.AssertCalled(t, "GetUsersBatch", ts.ctx, mock.Anything, mock.Anything)
			}
		})
	}
}

func TestGetNewUsers(t *testing.T) {
	ts := setup()

	users := []entity.User{
		{
			ID:       uuid.New(),
			Username: "Username1",
			Email:    "email1@test.com",
		},
		{
			ID:       uuid.New(),
			Username: "Username2",
			Email:    "email2@test.com",
		},
	}

	fromTime := time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)
	toTime := time.Date(2023, 9, 30, 23, 59, 59, 0, time.UTC)

	mockRepoFnOk := func(from, to time.Time, users []entity.User) {
		ts.mockRepo.On("GetNewUsers", ts.ctx, from, to).Return(users, nil)
	}

	tests := []struct {
		name       string
		users      []entity.User
		from, to   time.Time
		mockRepoFn func(from, to time.Time, users []entity.User)
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "success - get users within date range",
			users:      users,
			from:       fromTime,
			to:         toTime,
			mockRepoFn: mockRepoFnOk,
			wantErr:    false,
		},
		{
			name:       "fail - 'from' date is after 'to' date",
			from:       toTime,
			to:         fromTime,
			mockRepoFn: func(from, to time.Time, users []entity.User) {},
			wantErr:    true,
			errMsg:     v1.ErrInvalidFromTo.Error(),
		},
		{
			name:       "success - no users in date range",
			users:      []entity.User{},
			from:       fromTime,
			to:         toTime,
			mockRepoFn: mockRepoFnOk,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()
			tt.mockRepoFn(tt.from, tt.to, tt.users)

			fmt.Println(tt.name, tt.users)

			users, err := ts.userUseCase.GetNewUsers(ts.ctx, tt.from, tt.to)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
				ts.mockRepo.AssertNotCalled(t, "GetNewUsers")
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.users, users)
				ts.mockRepo.AssertCalled(t, "GetNewUsers", ts.ctx, mock.Anything, mock.Anything)
			}

			// TODO: Probably add this to some other tests and remove Parallel(),
			// because mockRepo.On(...) is hardcoded throughout the test suit, so
			// it will always expect A if mockRepo.On(...).Return(A) was called first.
			t.Cleanup(func() {
				ts.mockRepo.ExpectedCalls = nil
				ts.mockRepo.Calls = nil
			})
		})
	}
}

func TestUpdateUser(t *testing.T) {
	ts := setup()

	userExisting := &entity.User{
		ID:       uuid.New(),
		Username: "ExistingUser",
		Email:    "existing@user.com",
	}

	userNonExisting := &entity.User{
		ID:       uuid.New(),
		Username: "NonExistingUser",
		Email:    "nonexisting@user.com",
	}

	tests := []struct {
		name       string
		user       *entity.User
		mockRepoFn func()
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success - update existing user",
			user: userExisting,
			mockRepoFn: func() {
				ts.mockRepo.On("GetUserByID", ts.ctx, userExisting.ID).Return(userExisting, nil)
				ts.mockRepo.On("UpdateUser", ts.ctx, userExisting).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "fail - user does not exist",
			user: userNonExisting,
			mockRepoFn: func() {
				ts.mockRepo.On("GetUserByID", ts.ctx, userNonExisting.ID).Return(nil, errors.New(""))
				ts.mockRepo.On("UpdateUser", ts.ctx, userNonExisting).Return(errors.New(""))
			},
			wantErr: true,
			errMsg:  "user does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn()

			err := ts.userUseCase.UpdateUser(ts.ctx, tt.user)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
				ts.mockRepo.AssertNotCalled(t, "UpdateUser")
			} else {
				assert.Nil(t, err)
				ts.mockRepo.AssertCalled(t, "UpdateUser", ts.ctx, mock.Anything)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	ts := setup()

	userExistingID := uuid.New()
	userExisting := &entity.User{
		ID:       userExistingID,
		Username: "ExistingUser",
		Email:    "existing@user.com",
	}
	userNonExistingID := uuid.New()

	tests := []struct {
		name       string
		userID     uuid.UUID
		mockRepoFn func()
		wantErr    bool
		errMsg     string
	}{
		{
			name:   "success - delete existing user",
			userID: userExistingID,
			mockRepoFn: func() {
				ts.mockRepo.On("GetUserByID", ts.ctx, userExistingID).Return(userExisting, nil)
				ts.mockRepo.On("DeleteUser", ts.ctx, userExistingID).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "fail - user does not exist",
			userID: userNonExistingID,
			mockRepoFn: func() {
				ts.mockRepo.On("GetUserByID", ts.ctx, userNonExistingID).Return(nil, errors.New(""))
				ts.mockRepo.On("DeleteUser", ts.ctx, userNonExistingID).Return(errors.New(""))
			},
			wantErr: true,
			errMsg:  "user does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn()

			err := ts.userUseCase.DeleteUser(ts.ctx, tt.userID)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
				ts.mockRepo.AssertNotCalled(t, "DeleteUser")
			} else {
				assert.Nil(t, err)
				ts.mockRepo.AssertCalled(t, "DeleteUser", ts.ctx, mock.Anything)
			}
		})
	}
}

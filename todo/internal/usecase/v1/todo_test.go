package v1_test

import (
	"context"
	"errors"
	"testing"
	"time"
	"todo/internal/entity"
	"todo/internal/usecase"
	v1 "todo/internal/usecase/v1"
	"todo/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testSetup struct {
	ctx            context.Context
	mockBoardRepo  *mocks.BoardRepository
	mockColumnRepo *mocks.ColumnRepository
	mockCardRepo   *mocks.CardRepository
	todoUseCase    usecase.TodoUseCase
}

func setup() *testSetup {
	ctx := context.TODO()
	mockBoardRepo := new(mocks.BoardRepository)
	mockColumnRepo := new(mocks.ColumnRepository)
	mockCardRepo := new(mocks.CardRepository)
	todoUseCase := v1.NewTodoUseCase(mockBoardRepo, mockColumnRepo, mockCardRepo)

	return &testSetup{
		ctx:            ctx,
		mockBoardRepo:  mockBoardRepo,
		mockColumnRepo: mockColumnRepo,
		mockCardRepo:   mockCardRepo,
		todoUseCase:    todoUseCase,
	}
}

// CreateBoard(ctx context.Context, board *entity.Board) error
func TestCreateBoard(t *testing.T) {
	ts := setup()

	tests := []struct {
		name       string
		board      *entity.Board
		mockRepoFn func(board *entity.Board)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			board: &entity.Board{
				UserID: uuid.New(),
				Title:  "Title",
			},
			mockRepoFn: func(board *entity.Board) {
				ts.mockBoardRepo.On("CreateBoard", ts.ctx, board).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "board empty title",
			board: &entity.Board{
				UserID: uuid.New(),
			},
			mockRepoFn: func(board *entity.Board) {},
			wantErr:    true,
			errMsg:     v1.ErrBoardEmptyTitle.Error(),
		},
		{
			name: "board no user id",
			board: &entity.Board{
				Title: "Title",
			},
			mockRepoFn: func(board *entity.Board) {},
			wantErr:    true,
			errMsg:     v1.ErrBoardNoUserID.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn(tt.board)

			err := ts.todoUseCase.CreateBoard(ts.ctx, tt.board)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
				ts.mockBoardRepo.AssertNotCalled(t, "CreateBoard")
			} else {
				assert.Nil(t, err)
				ts.mockBoardRepo.AssertCalled(t, "CreateBoard", ts.ctx, mock.Anything)
			}
		})
	}
}

// GetBoardByID(ctx context.Context, id uuid.UUID) (*entity.Board, error)
func TestGetBoardByID(t *testing.T) {
	ts := setup()

	tests := []struct {
		name       string
		board      *entity.Board
		mockRepoFn func(board *entity.Board)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			board: &entity.Board{
				ID:     uuid.New(),
				UserID: uuid.New(),
				Title:  "Title",
			},
			mockRepoFn: func(board *entity.Board) {
				ts.mockBoardRepo.On("GetBoardByID", ts.ctx, board.ID).Return(board, nil)
			},
			wantErr: false,
		},
		{
			name: "failed to get board (not found for example)",
			board: &entity.Board{
				ID: uuid.New(),
			},
			mockRepoFn: func(board *entity.Board) {
				ts.mockBoardRepo.On("GetBoardByID", ts.ctx, mock.Anything).Return(nil, errors.New(""))
			},
			wantErr: true,
			errMsg:  v1.ErrGetBoardByID.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn(tt.board)

			board, err := ts.todoUseCase.GetBoardByID(ts.ctx, tt.board.ID)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, board, tt.board)
				ts.mockBoardRepo.AssertCalled(t, "GetBoardByID", ts.ctx, mock.Anything)
			}
		})
	}
}

// GetBoardsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.Board, error)
func TestGetBoardsByUser(t *testing.T) {
	ts := setup()

	userID := uuid.New()
	boards := []entity.Board{
		{
			ID:     uuid.New(),
			UserID: userID,
			Title:  "Board 1",
		},
		{
			ID:     uuid.New(),
			UserID: userID,
			Title:  "Board 2",
		},
		{
			ID:     uuid.New(),
			UserID: userID,
			Title:  "Board 3",
		},
		{
			ID:     uuid.New(),
			UserID: userID,
			Title:  "Board 4",
		},
		{
			ID:     uuid.New(),
			UserID: userID,
			Title:  "Board 5",
		},
	}

	tests := []struct {
		name          string
		userID        uuid.UUID
		limit, offset int
		boards        []entity.Board
		mockRepoFn    func(userID uuid.UUID, limit, offset int, boards []entity.Board)
		wantErr       bool
		errMsg        string
	}{
		{
			name:   "success, all",
			userID: userID,
			limit:  5,
			offset: 0,
			boards: boards[0:5],
			mockRepoFn: func(userID uuid.UUID, limit, offset int, boards []entity.Board) {
				ts.mockBoardRepo.On("GetBoardsByUser", ts.ctx, userID, limit, offset).Return(boards, nil)
			},
			wantErr: false,
		},
		{
			name:   "success, skip first",
			userID: userID,
			limit:  4,
			offset: 1,
			boards: boards[1:5],
			mockRepoFn: func(userID uuid.UUID, limit, offset int, boards []entity.Board) {
				ts.mockBoardRepo.On("GetBoardsByUser", ts.ctx, userID, limit, offset).Return(boards, nil)
			},
			wantErr: false,
		},
		{
			name:   "success, skip first and two last",
			userID: userID,
			limit:  2,
			offset: 1,
			boards: boards[1:3],
			mockRepoFn: func(userID uuid.UUID, limit, offset int, boards []entity.Board) {
				ts.mockBoardRepo.On("GetBoardsByUser", ts.ctx, userID, limit, offset).Return(boards, nil)
			},
			wantErr: false,
		},
		{
			name:       "negative limit",
			userID:     userID,
			limit:      -5,
			offset:     0,
			boards:     boards,
			mockRepoFn: func(boardID uuid.UUID, limit, offset int, boards []entity.Board) {},
			wantErr:    true,
			errMsg:     v1.ErrNegativeLimitOrOffset.Error(),
		},
		{
			name:       "negative offset",
			userID:     userID,
			limit:      5,
			offset:     -1,
			boards:     boards,
			mockRepoFn: func(boardID uuid.UUID, limit, offset int, boards []entity.Board) {},
			wantErr:    true,
			errMsg:     v1.ErrNegativeLimitOrOffset.Error(),
		},
		{
			name:       "zero limit",
			userID:     userID,
			limit:      0,
			offset:     0,
			boards:     boards,
			mockRepoFn: func(boardID uuid.UUID, limit, offset int, boards []entity.Board) {},
			wantErr:    true,
			errMsg:     v1.ErrZeroLimit.Error(),
		},
		{
			name:   "failed to get boards by user (not found for example)",
			userID: uuid.New(),
			limit:  5,
			offset: 0,
			boards: boards,
			mockRepoFn: func(userID uuid.UUID, limit, offset int, boards []entity.Board) {
				ts.mockBoardRepo.On("GetBoardsByUser", ts.ctx, userID, limit, offset).Return(nil, errors.New(""))
			},
			wantErr: true,
			errMsg:  v1.ErrGetBoardsByUser.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn(tt.userID, tt.limit, tt.offset, tt.boards)

			var err error
			boards, err := ts.todoUseCase.GetBoardsByUser(ts.ctx, tt.userID, tt.limit, tt.offset)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, boards, tt.boards)
				ts.mockBoardRepo.AssertCalled(t, "GetBoardsByUser", ts.ctx, mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}

// UpdateBoard(ctx context.Context, board *entity.Board) error
func TestUpdateBoard(t *testing.T) {
	ts := setup()

	tests := []struct {
		name       string
		board      *entity.Board
		mockRepoFn func(board *entity.Board)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			board: &entity.Board{
				UserID: uuid.New(),
				Title:  "Title",
			},
			mockRepoFn: func(board *entity.Board) {
				ts.mockBoardRepo.On("UpdateBoard", ts.ctx, board).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "board empty title",
			board: &entity.Board{
				UserID: uuid.New(),
			},
			mockRepoFn: func(board *entity.Board) {
				ts.mockBoardRepo.On("UpdateBoard", ts.ctx, board).Return(v1.ErrBoardEmptyTitle)
			},
			wantErr: true,
			errMsg:  v1.ErrBoardEmptyTitle.Error(),
		},
		{
			name: "board no user id",
			board: &entity.Board{
				Title: "Title",
			},
			mockRepoFn: func(board *entity.Board) {
				ts.mockBoardRepo.On("UpdateBoard", ts.ctx, board).Return(v1.ErrBoardNoUserID)
			},
			wantErr: true,
			errMsg:  v1.ErrBoardNoUserID.Error(),
		},
		{
			name: "failed to update board (not found for example)",
			board: &entity.Board{
				UserID: uuid.New(),
				Title:  "Title",
			},
			mockRepoFn: func(board *entity.Board) {
				ts.mockBoardRepo.On("UpdateBoard", ts.ctx, board).Return(v1.ErrUpdateBoard)
			},
			wantErr: true,
			errMsg:  v1.ErrUpdateBoard.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn(tt.board)

			err := ts.todoUseCase.UpdateBoard(ts.ctx, tt.board)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.Nil(t, err)
				ts.mockBoardRepo.AssertCalled(t, "UpdateBoard", ts.ctx, mock.Anything)
			}
		})
	}
}

// DeleteBoard(ctx context.Context, id uuid.UUID) error
func TestDeleteBoard(t *testing.T) {
	ts := setup()

	tests := []struct {
		name       string
		boardID    uuid.UUID
		mockRepoFn func(id uuid.UUID)
		wantErr    bool
		errMsg     string
	}{
		{
			name:    "success",
			boardID: uuid.New(),
			mockRepoFn: func(id uuid.UUID) {
				ts.mockBoardRepo.On("DeleteBoard", ts.ctx, id).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "failed to delete board (not found for example)",
			boardID: uuid.New(),
			mockRepoFn: func(id uuid.UUID) {
				ts.mockBoardRepo.On("DeleteBoard", ts.ctx, mock.Anything).Return(errors.New(""))
			},
			wantErr: true,
			errMsg:  v1.ErrDeleteBoard.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn(tt.boardID)

			err := ts.todoUseCase.DeleteBoard(ts.ctx, tt.boardID)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.Nil(t, err)
				ts.mockBoardRepo.AssertCalled(t, "DeleteBoard", ts.ctx, mock.Anything)
			}
		})
	}
}

// CreateColumn(ctx context.Context, column *entity.Column) error
func TestCreateColumn(t *testing.T) {
	ts := setup()

	tests := []struct {
		name       string
		column     *entity.Column
		mockRepoFn func(column *entity.Column)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			column: &entity.Column{
				Title:    "Title",
				UserID:   uuid.New(),
				BoardID:  uuid.New(),
				Position: 0,
			},
			mockRepoFn: func(column *entity.Column) {
				ts.mockColumnRepo.On("CreateColumn", ts.ctx, column).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "column empty title",
			column: &entity.Column{
				Title:    "",
				UserID:   uuid.New(),
				BoardID:  uuid.New(),
				Position: 0,
			},
			mockRepoFn: func(column *entity.Column) {},
			wantErr:    true,
			errMsg:     v1.ErrColumnEmptyTitle.Error(),
		},
		{
			name: "column no user id",
			column: &entity.Column{
				Title:    "Title",
				BoardID:  uuid.New(),
				Position: 0,
			},
			mockRepoFn: func(column *entity.Column) {},
			wantErr:    true,
			errMsg:     v1.ErrColumnNoUserID.Error(),
		},
		{
			name: "column no board id",
			column: &entity.Column{
				Title:    "Title",
				UserID:   uuid.New(),
				Position: 0,
			},
			mockRepoFn: func(column *entity.Column) {},
			wantErr:    true,
			errMsg:     v1.ErrColumnNoBoardID.Error(),
		},
		{
			name: "column negative position",
			column: &entity.Column{
				Title:    "Title",
				UserID:   uuid.New(),
				BoardID:  uuid.New(),
				Position: -1,
			},
			mockRepoFn: func(column *entity.Column) {},
			wantErr:    true,
			errMsg:     v1.ErrColumnNegativePosition.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn(tt.column)

			err := ts.todoUseCase.CreateColumn(ts.ctx, tt.column)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
				ts.mockColumnRepo.AssertNotCalled(t, "CreateColumn")
			} else {
				assert.Nil(t, err)
				ts.mockColumnRepo.AssertCalled(t, "CreateColumn", ts.ctx, mock.Anything)
			}
		})
	}
}

// GetColumnByID(ctx context.Context, id uuid.UUID) (*entity.Column, error)
func TestGetColumnByID(t *testing.T) {
	ts := setup()

	tests := []struct {
		name       string
		column     *entity.Column
		mockRepoFn func(column *entity.Column)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			column: &entity.Column{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				BoardID:  uuid.New(),
				Title:    "Title",
				Position: 1,
			},
			mockRepoFn: func(column *entity.Column) {
				ts.mockColumnRepo.On("GetColumnByID", ts.ctx, column.ID).Return(column, nil)
			},
			wantErr: false,
		},
		{
			name: "failed to get column (not found for example)",
			column: &entity.Column{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				BoardID:  uuid.New(),
				Title:    "Title",
				Position: 1,
			},
			mockRepoFn: func(column *entity.Column) {
				ts.mockColumnRepo.On("GetColumnByID", ts.ctx, mock.Anything).Return(nil, errors.New(""))
			},
			wantErr: true,
			errMsg:  v1.ErrGetColumnByID.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn(tt.column)

			column, err := ts.todoUseCase.GetColumnByID(ts.ctx, tt.column.ID)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, column, tt.column)
				ts.mockColumnRepo.AssertCalled(t, "GetColumnByID", ts.ctx, mock.Anything)
			}
		})
	}
}

// GetColumnsByBoard(ctx context.Context, boardID uuid.UUID, limit, offset int) ([]entity.Column, error)
func TestGetColumnsByBoard(t *testing.T) {
	ts := setup()

	boardID := uuid.New()
	columns := []entity.Column{
		{
			ID:       uuid.New(),
			UserID:   uuid.New(),
			BoardID:  boardID,
			Title:    "Column 1",
			Position: 0,
		},
		{
			ID:       uuid.New(),
			UserID:   uuid.New(),
			BoardID:  boardID,
			Title:    "Column 2",
			Position: 1,
		},
		{
			ID:       uuid.New(),
			UserID:   uuid.New(),
			BoardID:  boardID,
			Title:    "Column 3",
			Position: 2,
		},
		{
			ID:       uuid.New(),
			UserID:   uuid.New(),
			BoardID:  boardID,
			Title:    "Column 4",
			Position: 3,
		},
		{
			ID:       uuid.New(),
			UserID:   uuid.New(),
			BoardID:  boardID,
			Title:    "Column 5",
			Position: 4,
		},
	}

	tests := []struct {
		name          string
		boardID       uuid.UUID
		limit, offset int
		columns       []entity.Column
		mockRepoFn    func(boardID uuid.UUID, limit, offset int, columns []entity.Column)
		wantErr       bool
		errMsg        string
	}{
		{
			name:    "success, all",
			boardID: boardID,
			limit:   5,
			offset:  0,
			columns: columns[0:5],
			mockRepoFn: func(boardID uuid.UUID, limit, offset int, columns []entity.Column) {
				ts.mockColumnRepo.On("GetColumnsByBoard", ts.ctx, boardID, limit, offset).Return(columns, nil)
			},
			wantErr: false,
		},
		{
			name:    "success, skip first",
			boardID: boardID,
			limit:   4,
			offset:  1,
			columns: columns[1:5],
			mockRepoFn: func(boardID uuid.UUID, limit, offset int, columns []entity.Column) {
				ts.mockColumnRepo.On("GetColumnsByBoard", ts.ctx, boardID, limit, offset).Return(columns, nil)
			},
			wantErr: false,
		},
		{
			name:    "success, skip first and two last",
			boardID: boardID,
			limit:   2,
			offset:  1,
			columns: columns[1:3],
			mockRepoFn: func(boardID uuid.UUID, limit, offset int, columns []entity.Column) {
				ts.mockColumnRepo.On("GetColumnsByBoard", ts.ctx, boardID, limit, offset).Return(columns, nil)
			},
			wantErr: false,
		},
		{
			name:       "negative limit",
			boardID:    boardID,
			limit:      -5,
			offset:     0,
			columns:    columns,
			mockRepoFn: func(boardID uuid.UUID, limit, offset int, columns []entity.Column) {},
			wantErr:    true,
			errMsg:     v1.ErrNegativeLimitOrOffset.Error(),
		},
		{
			name:       "negative offset",
			boardID:    boardID,
			limit:      5,
			offset:     -1,
			columns:    columns,
			mockRepoFn: func(boardID uuid.UUID, limit, offset int, columns []entity.Column) {},
			wantErr:    true,
			errMsg:     v1.ErrNegativeLimitOrOffset.Error(),
		},
		{
			name:       "zero limit",
			boardID:    boardID,
			limit:      0,
			offset:     0,
			columns:    columns,
			mockRepoFn: func(boardID uuid.UUID, limit, offset int, columns []entity.Column) {},
			wantErr:    true,
			errMsg:     v1.ErrZeroLimit.Error(),
		},
		{
			name:    "failed to get columns by board (not found for example)",
			boardID: uuid.New(),
			limit:   5,
			offset:  0,
			columns: columns,
			mockRepoFn: func(boardID uuid.UUID, limit, offset int, columns []entity.Column) {
				ts.mockColumnRepo.On("GetColumnsByBoard", ts.ctx, boardID, limit, offset).Return(nil, errors.New(""))
			},
			wantErr: true,
			errMsg:  v1.ErrGetColumnsByBoard.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn(tt.boardID, tt.limit, tt.offset, tt.columns)

			var err error
			columns, err := ts.todoUseCase.GetColumnsByBoard(ts.ctx, tt.boardID, tt.limit, tt.offset)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, columns, tt.columns)
				ts.mockColumnRepo.AssertCalled(t, "GetColumnsByBoard", ts.ctx, mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}

// UpdateColumn(ctx context.Context, column *entity.Column) error
func TestUpdateColumn(t *testing.T) {
	ts := setup()

	tests := []struct {
		name       string
		column     *entity.Column
		mockRepoFn func(column *entity.Column)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			column: &entity.Column{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				BoardID:  uuid.New(),
				Title:    "Title",
				Position: 0,
			},
			mockRepoFn: func(column *entity.Column) {
				ts.mockColumnRepo.On("UpdateColumn", ts.ctx, column).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "empty title",
			column: &entity.Column{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				BoardID:  uuid.New(),
				Title:    "",
				Position: 0,
			},
			mockRepoFn: func(column *entity.Column) {},
			wantErr:    true,
			errMsg:     v1.ErrColumnEmptyTitle.Error(),
		},
		{
			name: "no user id",
			column: &entity.Column{
				ID:       uuid.New(),
				BoardID:  uuid.New(),
				Title:    "Title",
				Position: 0,
			},
			mockRepoFn: func(column *entity.Column) {},
			wantErr:    true,
			errMsg:     v1.ErrColumnNoUserID.Error(),
		},
		{
			name: "no board id",
			column: &entity.Column{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				Title:    "Title",
				Position: 0,
			},
			mockRepoFn: func(column *entity.Column) {},
			wantErr:    true,
			errMsg:     v1.ErrColumnNoBoardID.Error(),
		},
		{
			name: "negative position",
			column: &entity.Column{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				BoardID:  uuid.New(),
				Title:    "Title",
				Position: -1,
			},
			mockRepoFn: func(column *entity.Column) {},
			wantErr:    true,
			errMsg:     v1.ErrColumnNegativePosition.Error(),
		},
		{
			name: "failed to update column (not found for example)",
			column: &entity.Column{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				BoardID:  uuid.New(),
				Title:    "Title",
				Position: 0,
			},
			mockRepoFn: func(column *entity.Column) {
				ts.mockColumnRepo.On("UpdateColumn", ts.ctx, column).Return(errors.New(""))
			},
			wantErr: true,
			errMsg:  v1.ErrUpdateColumn.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn(tt.column)

			err := ts.todoUseCase.UpdateColumn(ts.ctx, tt.column)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.Nil(t, err)
				ts.mockColumnRepo.AssertCalled(t, "UpdateColumn", ts.ctx, mock.Anything)
			}
		})
	}
}

// DeleteColumn(ctx context.Context, id uuid.UUID) error
func TestDeleteColumn(t *testing.T) {
	ts := setup()

	tests := []struct {
		name       string
		columnID   uuid.UUID
		mockRepoFn func(id uuid.UUID)
		wantErr    bool
		errMsg     string
	}{
		{
			name:     "success",
			columnID: uuid.New(),
			mockRepoFn: func(id uuid.UUID) {
				ts.mockColumnRepo.On("DeleteColumn", ts.ctx, id).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "failed to delete column (not found for example)",
			columnID: uuid.New(),
			mockRepoFn: func(id uuid.UUID) {
				ts.mockColumnRepo.On("DeleteColumn", ts.ctx, id).Return(errors.New(""))
			},
			wantErr: true,
			errMsg:  v1.ErrDeleteColumn.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn(tt.columnID)

			err := ts.todoUseCase.DeleteColumn(ts.ctx, tt.columnID)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.Nil(t, err)
				ts.mockColumnRepo.AssertCalled(t, "DeleteColumn", ts.ctx, mock.Anything)
			}
		})
	}
}

// CreateCard(ctx context.Context, card *entity.Card) error
func TestCreateCard(t *testing.T) {
	ts := setup()

	tests := []struct {
		name       string
		card       *entity.Card
		mockRepoFn func(card *entity.Card)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			card: &entity.Card{
				UserID:      uuid.New(),
				ColumnID:    uuid.New(),
				Title:       "Title",
				Description: "Description",
				Position:    0,
			},
			mockRepoFn: func(card *entity.Card) {
				ts.mockCardRepo.On("CreateCard", ts.ctx, card).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "no user id",
			card: &entity.Card{
				ColumnID: uuid.New(),
				Title:    "Title",
				Position: 0,
			},
			mockRepoFn: func(card *entity.Card) {},
			wantErr:    true,
			errMsg:     v1.ErrCardNoUserID.Error(),
		},
		{
			name: "no column id",
			card: &entity.Card{
				UserID:   uuid.New(),
				Title:    "Title",
				Position: 0,
			},
			mockRepoFn: func(card *entity.Card) {},
			wantErr:    true,
			errMsg:     v1.ErrCardNoColumnID.Error(),
		},
		{
			name: "negative position",
			card: &entity.Card{
				UserID:   uuid.New(),
				ColumnID: uuid.New(),
				Title:    "Title",
				Position: -1,
			},
			mockRepoFn: func(card *entity.Card) {},
			wantErr:    true,
			errMsg:     v1.ErrCardNegativePosition.Error(),
		},
		{
			name: "empty title",
			card: &entity.Card{
				UserID:   uuid.New(),
				ColumnID: uuid.New(),
				Title:    "",
				Position: 0,
			},
			mockRepoFn: func(card *entity.Card) {},
			wantErr:    true,
			errMsg:     v1.ErrCardEmptyTitle.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn(tt.card)

			err := ts.todoUseCase.CreateCard(ts.ctx, tt.card)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
				ts.mockCardRepo.AssertNotCalled(t, "CreateCard")
			} else {
				assert.Nil(t, err)
				ts.mockCardRepo.AssertCalled(t, "CreateCard", ts.ctx, mock.Anything)
			}
		})
	}
}

// GetCardByID(ctx context.Context, id uuid.UUID) (*entity.Card, error)
func TestGetCardByID(t *testing.T) {
	ts := setup()

	tests := []struct {
		name       string
		card       *entity.Card
		mockRepoFn func(card *entity.Card)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			card: &entity.Card{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				ColumnID: uuid.New(),
				Title:    "Title",
				Position: 0,
			},
			mockRepoFn: func(card *entity.Card) {
				ts.mockCardRepo.On("GetCardByID", ts.ctx, card.ID).Return(card, nil)
			},
			wantErr: false,
		},
		{
			name: "failed to get card by id (not found for example)",
			card: &entity.Card{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				ColumnID: uuid.New(),
				Title:    "Title",
				Position: 0,
			},
			mockRepoFn: func(card *entity.Card) {
				ts.mockCardRepo.On("GetCardByID", ts.ctx, card.ID).Return(nil, errors.New(""))
			},
			wantErr: true,
			errMsg:  v1.ErrGetCardByID.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn(tt.card)

			card, err := ts.todoUseCase.GetCardByID(ts.ctx, tt.card.ID)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, card, tt.card)
				ts.mockCardRepo.AssertCalled(t, "GetCardByID", ts.ctx, mock.Anything)
			}
		})
	}
}

// GetCardsByColumn(ctx context.Context, columnID uuid.UUID, limit, offset int) ([]entity.Card, error)
func TestGetCardsByColumn(t *testing.T) {
	ts := setup()

	cards := []entity.Card{
		{
			ID:       uuid.New(),
			UserID:   uuid.New(),
			ColumnID: uuid.New(),
			Title:    "Card 1",
			Position: 0,
		},
		{
			ID:       uuid.New(),
			UserID:   uuid.New(),
			ColumnID: uuid.New(),
			Title:    "Card 2",
			Position: 1,
		},
		{
			ID:       uuid.New(),
			UserID:   uuid.New(),
			ColumnID: uuid.New(),
			Title:    "Card 3",
			Position: 2,
		},
		{
			ID:       uuid.New(),
			UserID:   uuid.New(),
			ColumnID: uuid.New(),
			Title:    "Card 4",
			Position: 3,
		},
		{
			ID:       uuid.New(),
			UserID:   uuid.New(),
			ColumnID: uuid.New(),
			Title:    "Card 5",
			Position: 4,
		},
	}

	tests := []struct {
		name          string
		columnID      uuid.UUID
		limit, offset int
		cards         []entity.Card
		mockRepoFn    func(columnID uuid.UUID, limit, offset int, cards []entity.Card)
		wantErr       bool
		errMsg        string
	}{
		{
			name:     "success, all",
			columnID: uuid.New(),
			limit:    5,
			offset:   0,
			cards:    cards,
			mockRepoFn: func(columnID uuid.UUID, limit, offset int, cards []entity.Card) {
				ts.mockCardRepo.On("GetCardsByColumn", ts.ctx, columnID, limit, offset).Return(cards, nil)
			},
			wantErr: false,
		},
		{
			name:     "success, skip first",
			columnID: uuid.New(),
			limit:    4,
			offset:   1,
			cards:    cards[1:5],
			mockRepoFn: func(columnID uuid.UUID, limit, offset int, cards []entity.Card) {
				ts.mockCardRepo.On("GetCardsByColumn", ts.ctx, columnID, limit, offset).Return(cards, nil)
			},
			wantErr: false,
		},
		{
			name:     "success, skip first and last two",
			columnID: uuid.New(),
			limit:    2,
			offset:   1,
			cards:    cards[1:3],
			mockRepoFn: func(columnID uuid.UUID, limit, offset int, cards []entity.Card) {
				ts.mockCardRepo.On("GetCardsByColumn", ts.ctx, columnID, limit, offset).Return(cards, nil)
			},
			wantErr: false,
		},
		{
			name:       "negative limit",
			columnID:   uuid.New(),
			limit:      -1,
			offset:     0,
			cards:      cards,
			mockRepoFn: func(columnID uuid.UUID, limit, offset int, cards []entity.Card) {},
			wantErr:    true,
			errMsg:     v1.ErrNegativeLimitOrOffset.Error(),
		},
		{
			name:       "negative offset",
			columnID:   uuid.New(),
			limit:      5,
			offset:     -1,
			cards:      cards,
			mockRepoFn: func(columnID uuid.UUID, limit, offset int, cards []entity.Card) {},
			wantErr:    true,
			errMsg:     v1.ErrNegativeLimitOrOffset.Error(),
		},
		{
			name:       "negative offset",
			columnID:   uuid.New(),
			limit:      0,
			offset:     0,
			cards:      cards,
			mockRepoFn: func(columnID uuid.UUID, limit, offset int, cards []entity.Card) {},
			wantErr:    true,
			errMsg:     v1.ErrZeroLimit.Error(),
		},
		{
			name:     "failed to get cards by column (for example not found)",
			columnID: uuid.New(),
			limit:    5,
			offset:   0,
			cards:    cards,
			mockRepoFn: func(columnID uuid.UUID, limit, offset int, cards []entity.Card) {
				ts.mockCardRepo.On("GetCardsByColumn", ts.ctx, columnID, limit, offset).Return(nil, errors.New(""))
			},
			wantErr: true,
			errMsg:  v1.ErrGetCardsByColumn.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn(tt.columnID, tt.limit, tt.offset, tt.cards)

			cards, err := ts.todoUseCase.GetCardsByColumn(ts.ctx, tt.columnID, tt.limit, tt.offset)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, cards, tt.cards)
				ts.mockCardRepo.AssertCalled(t, "GetCardsByColumn", ts.ctx, mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}

func TestGetNewCards(t *testing.T) {
	ts := setup()

	cards := []entity.Card{
		{
			ID:    uuid.New(),
			Title: "Card 1",
		},
		{
			ID:    uuid.New(),
			Title: "Card 2",
		},
		{
			ID:    uuid.New(),
			Title: "Card 3",
		},
	}

	fromTime := time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)
	toTime := time.Date(2023, 9, 30, 23, 59, 59, 0, time.UTC)

	mockRepoFnOk := func(from, to time.Time, cards []entity.Card) {
		ts.mockCardRepo.On("GetNewCards", ts.ctx, from, to).Return(cards, nil)
	}

	tests := []struct {
		name       string
		from, to   time.Time
		cards      []entity.Card
		mockRepoFn func(from, to time.Time, cards []entity.Card)
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "success, all",
			from:       fromTime,
			to:         toTime,
			cards:      cards,
			mockRepoFn: mockRepoFnOk,
			wantErr:    false,
		},
		{
			name:       "invalid time range (from greater than to)",
			from:       toTime,
			to:         fromTime,
			cards:      cards,
			mockRepoFn: func(from, to time.Time, cards []entity.Card) {},
			wantErr:    true,
			errMsg:     v1.ErrInvalidTimeRange.Error(),
		},
		{
			name:       "success, but no cards found",
			from:       fromTime,
			to:         toTime,
			cards:      []entity.Card{},
			mockRepoFn: mockRepoFnOk,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoFn(tt.from, tt.to, tt.cards)

			cards, err := ts.todoUseCase.GetNewCards(ts.ctx, tt.from, tt.to)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, cards, tt.cards)
				ts.mockCardRepo.AssertCalled(t, "GetNewCards", ts.ctx, mock.Anything, mock.Anything, mock.Anything)
			}

			t.Cleanup(func() {
				ts.mockCardRepo.ExpectedCalls = nil
				ts.mockCardRepo.Calls = nil
			})
		})
	}
}

// UpdateCard(ctx context.Context, card *entity.Card) error
func TestUpdateCard(t *testing.T) {
	ts := setup()

	tests := []struct {
		name       string
		card       *entity.Card
		mockRepoFn func(card *entity.Card)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			card: &entity.Card{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				ColumnID: uuid.New(),
				Title:    "Title",
				Position: 0,
			},
			mockRepoFn: func(card *entity.Card) {
				ts.mockCardRepo.On("UpdateCard", ts.ctx, card).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "no user id",
			card: &entity.Card{
				ID:       uuid.New(),
				ColumnID: uuid.New(),
				Title:    "Title",
				Position: 0,
			},
			mockRepoFn: func(card *entity.Card) {},
			wantErr:    true,
			errMsg:     v1.ErrCardNoUserID.Error(),
		},
		{
			name: "no column id",
			card: &entity.Card{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				Title:    "Title",
				Position: 0,
			},
			mockRepoFn: func(card *entity.Card) {},
			wantErr:    true,
			errMsg:     v1.ErrCardNoColumnID.Error(),
		},
		{
			name: "negative position",
			card: &entity.Card{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				ColumnID: uuid.New(),
				Title:    "Title",
				Position: -1,
			},
			mockRepoFn: func(card *entity.Card) {},
			wantErr:    true,
			errMsg:     v1.ErrCardNegativePosition.Error(),
		},
		{
			name: "empty title",
			card: &entity.Card{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				ColumnID: uuid.New(),
				Title:    "",
				Position: 0,
			},
			mockRepoFn: func(card *entity.Card) {},
			wantErr:    true,
			errMsg:     v1.ErrCardEmptyTitle.Error(),
		},
		{
			name: "failed to update card (not found for example)",
			card: &entity.Card{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				ColumnID: uuid.New(),
				Title:    "Title",
				Position: 0,
			},
			mockRepoFn: func(card *entity.Card) {
				ts.mockCardRepo.On("UpdateCard", ts.ctx, card).Return(errors.New(""))
			},
			wantErr: true,
			errMsg:  v1.ErrUpdateCard.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn(tt.card)

			err := ts.todoUseCase.UpdateCard(ts.ctx, tt.card)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.Nil(t, err)
				ts.mockCardRepo.AssertCalled(t, "UpdateCard", ts.ctx, mock.Anything)
			}
		})
	}
}

// DeleteCard(ctx context.Context, id uuid.UUID) error
func TestDeleteCard(t *testing.T) {
	ts := setup()

	tests := []struct {
		name       string
		cardID     uuid.UUID
		mockRepoFn func(id uuid.UUID)
		wantErr    bool
		errMsg     string
	}{
		{
			name:   "success",
			cardID: uuid.New(),
			mockRepoFn: func(id uuid.UUID) {
				ts.mockCardRepo.On("DeleteCard", ts.ctx, id).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "failed to delete card (not found for example)",
			cardID: uuid.New(),
			mockRepoFn: func(id uuid.UUID) {
				ts.mockCardRepo.On("DeleteCard", ts.ctx, id).Return(errors.New(""))
			},
			wantErr: true,
			errMsg:  v1.ErrDeleteCard.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockRepoFn(tt.cardID)

			err := ts.todoUseCase.DeleteCard(ts.ctx, tt.cardID)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.Nil(t, err)
				ts.mockCardRepo.AssertCalled(t, "DeleteCard", ts.ctx, mock.Anything)
			}
		})
	}
}

package v1_test

import (
	"aggregator/internal/dto"
	"aggregator/internal/entity"
	"aggregator/internal/usecase"
	v1 "aggregator/internal/usecase/v1"
	"aggregator/mocks"
	"context"
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testSetup struct {
	ctx         context.Context
	mockUserSvc *mocks.UserService
	mockAuthSvc *mocks.AuthService
	mockTodoSvc *mocks.TodoService
	uc          usecase.AggregatorUseCase
}

func setup() *testSetup {
	ctx := context.TODO()

	mockUserSvc := new(mocks.UserService)
	mockAuthSvc := new(mocks.AuthService)
	mockTodoSvc := new(mocks.TodoService)

	aggregatorUseCase := v1.NewAggregatorUseCase(mockUserSvc, mockAuthSvc, mockTodoSvc)

	return &testSetup{
		ctx:         ctx,
		mockUserSvc: mockUserSvc,
		mockAuthSvc: mockAuthSvc,
		mockTodoSvc: mockTodoSvc,
		uc:          aggregatorUseCase,
	}
}

func Date(yyyy, mm, dd int) time.Time {
	return time.Date(yyyy, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)
}

func TestGetStats(t *testing.T) {
	ts := setup()

	mockUserFnOk := func(from, to time.Time, users []dto.User) {
		ts.mockUserSvc.On("GetNewUsers", ts.ctx, from, to).Return(users, nil)
	}

	mockUserFnErr := func(from, to time.Time, users []dto.User) {
		ts.mockUserSvc.On("GetNewUsers", ts.ctx, from, to).Return(nil, errors.New(""))
	}

	mockTodoFnOk := func(from, to time.Time, cards []dto.Card) {
		ts.mockTodoSvc.On("GetNewCards", ts.ctx, from, to).Return(cards, nil)
	}

	mockTodoFnErr := func(from, to time.Time, cards []dto.Card) {
		ts.mockTodoSvc.On("GetNewCards", ts.ctx, from, to).Return(nil, errors.New(""))
	}

	dates := []time.Time{
		Date(2000, 1, 1),
		Date(2000, 1, 2),
		Date(2000, 1, 3),
		Date(2000, 1, 4),
	}

	nu := 4
	uuids := make([]uuid.UUID, nu)
	for i := 0; i < nu; i++ {
		uuids[i] = uuid.New()
	}

	users := make([]dto.User, nu)
	for i := 0; i < nu; i++ {
		users[i] = dto.User{
			ID: uuids[i],
		}
	}

	users[0].CreatedAt = dates[0]
	users[1].CreatedAt = dates[1]
	users[2].CreatedAt = dates[1]
	users[3].CreatedAt = dates[2]

	nc := 13
	cards := make([]dto.Card, nc)
	for i := 0; i < nc; i++ {
		cards[i].ID = uuid.New()
	}

	// User0:
	// - Day0:
	//   - Card0
	cards[0].UserID = uuids[0]
	cards[0].CreatedAt = dates[0]

	// - Day1:
	//   - Card1
	cards[1].UserID = uuids[0]
	cards[1].CreatedAt = dates[1]

	//   - Card2
	cards[2].UserID = uuids[0]
	cards[2].CreatedAt = dates[1]

	// - Day2:
	//   - Card6
	cards[6].UserID = uuids[0]
	cards[6].CreatedAt = dates[2]

	//   - Card7
	cards[7].UserID = uuids[0]
	cards[7].CreatedAt = dates[2]

	//   - Card8
	cards[8].UserID = uuids[0]
	cards[8].CreatedAt = dates[2]

	// - Day3:
	//   - Card11
	cards[11].UserID = uuids[0]
	cards[11].CreatedAt = dates[3]

	//   - Card12
	cards[12].UserID = uuids[0]
	cards[12].CreatedAt = dates[3]

	// User1:
	// - Day1:
	//   - Card3
	cards[3].UserID = uuids[1]
	cards[3].CreatedAt = dates[1]

	//   - Card4
	cards[4].UserID = uuids[1]
	cards[4].CreatedAt = dates[1]

	//   - Card5
	cards[5].UserID = uuids[1]
	cards[5].CreatedAt = dates[1]

	// User2:
	// - Day1:
	//   - Card9
	cards[9].UserID = uuids[2]
	cards[9].CreatedAt = dates[1]

	//   - Card10
	cards[10].UserID = uuids[2]
	cards[10].CreatedAt = dates[1]

	// User3:
	// - Day2:

	tests := []struct {
		name       string
		from, to   time.Time
		users      []dto.User
		cards      []dto.Card
		stats      []entity.NewUsersAndCardsStats
		mockUserFn func(from, to time.Time, users []dto.User)
		mockTodoFn func(from, to time.Time, cards []dto.Card)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success, first",
			from: dates[0],
			to:   dates[0],
			users: []dto.User{
				users[0],
			},
			cards: []dto.Card{
				cards[0],
			},
			stats: []entity.NewUsersAndCardsStats{
				{
					Date: dates[0],
					Users: dto.ToUserEntities([]dto.User{
						users[0],
					}),
					Cards: dto.ToCardEntities([]dto.Card{
						cards[0],
					}),
					NumCardsByNewUsers: 1,
				},
			},
			mockUserFn: mockUserFnOk,
			mockTodoFn: mockTodoFnOk,
			wantErr:    false,
		},
		{
			name:  "success, all",
			from:  dates[0],
			to:    dates[len(dates)-1],
			users: users,
			cards: cards,
			stats: []entity.NewUsersAndCardsStats{
				{
					Date: dates[0],
					Users: dto.ToUserEntities([]dto.User{
						users[0],
					}),
					Cards: dto.ToCardEntities([]dto.Card{
						cards[0],
					}),
					NumCardsByNewUsers: 1,
				},
				{
					Date: dates[1],
					Users: dto.ToUserEntities([]dto.User{
						users[1],
						users[2],
					}),
					Cards: dto.ToCardEntities([]dto.Card{
						cards[1],
						cards[2],
						cards[3],
						cards[4],
						cards[5],
						cards[9],
						cards[10],
					}),
					NumCardsByNewUsers: 5,
				},
				{
					Date: dates[2],
					Users: dto.ToUserEntities([]dto.User{
						users[3],
					}),
					Cards: dto.ToCardEntities([]dto.Card{
						cards[6],
						cards[7],
						cards[8],
					}),
					NumCardsByNewUsers: 0,
				},
				{
					Date:  dates[3],
					Users: dto.ToUserEntities([]dto.User{}),
					Cards: dto.ToCardEntities([]dto.Card{
						cards[11],
						cards[12],
					}),
					NumCardsByNewUsers: 0,
				},
			},
			mockUserFn: mockUserFnOk,
			mockTodoFn: mockTodoFnOk,
			wantErr:    false,
		},
		{
			name: "success, some",
			from: dates[1],
			to:   dates[len(dates)-2],
			users: []dto.User{
				users[1],
				users[2],
				users[3],
			},
			cards: []dto.Card{
				cards[1],
				cards[2],
				cards[3],
				cards[4],
				cards[5],
				cards[6],
				cards[7],
				cards[8],
				cards[9],
				cards[10],
			},
			stats: []entity.NewUsersAndCardsStats{
				{
					Date: dates[1],
					Users: dto.ToUserEntities([]dto.User{
						users[1],
						users[2],
					}),
					Cards: dto.ToCardEntities([]dto.Card{
						cards[1],
						cards[2],
						cards[3],
						cards[4],
						cards[5],
						cards[9],
						cards[10],
					}),
					NumCardsByNewUsers: 5,
				},
				{
					Date: dates[2],
					Users: dto.ToUserEntities([]dto.User{
						users[3],
					}),
					Cards: dto.ToCardEntities([]dto.Card{
						cards[6],
						cards[7],
						cards[8],
					}),
					NumCardsByNewUsers: 0,
				},
			},
			mockUserFn: mockUserFnOk,
			mockTodoFn: mockTodoFnOk,
			wantErr:    false,
		},
		{
			name:       "fail - <<from>> is greater than <<to>>",
			from:       dates[1],
			to:         dates[0],
			users:      users,
			cards:      cards,
			stats:      []entity.NewUsersAndCardsStats{},
			mockUserFn: func(from, to time.Time, users []dto.User) {},
			mockTodoFn: func(from, to time.Time, users []dto.Card) {},
			wantErr:    true,
			errMsg:     v1.ErrInvalidTimeRange.Error(),
		},
		{
			name:       "failed to get new users",
			from:       dates[0],
			to:         dates[1],
			users:      users,
			cards:      cards,
			stats:      []entity.NewUsersAndCardsStats{},
			mockUserFn: mockUserFnErr,
			mockTodoFn: mockTodoFnOk,
			wantErr:    true,
			errMsg:     v1.ErrGetNewUsers.Error(),
		},
		{
			name:       "failed to get new cards",
			from:       dates[0],
			to:         dates[1],
			users:      users,
			cards:      cards,
			stats:      []entity.NewUsersAndCardsStats{},
			mockUserFn: mockUserFnOk,
			mockTodoFn: mockTodoFnErr,
			wantErr:    true,
			errMsg:     v1.ErrGetNewCards.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUserFn(tt.from, tt.to, tt.users)
			tt.mockTodoFn(tt.from, tt.to, tt.cards)

			stats, err := ts.uc.GetStats(ts.ctx, tt.from, tt.to)

			csExpected := ToComparableStats(tt.stats)
			csActual := ToComparableStats(stats)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.Nil(t, err)
				// assert.Equal(t, tt.stats, stats)
				assert.Equal(t, csExpected, csActual)
				ts.mockUserSvc.AssertCalled(t, "GetNewUsers", ts.ctx, mock.Anything, mock.Anything, mock.Anything)
				ts.mockTodoSvc.AssertCalled(t, "GetNewCards", ts.ctx, mock.Anything, mock.Anything, mock.Anything)
			}

			t.Cleanup(func() {
				ts.mockUserSvc.ExpectedCalls = nil
				ts.mockUserSvc.Calls = nil

				ts.mockTodoSvc.ExpectedCalls = nil
				ts.mockTodoSvc.Calls = nil
			})
		})
	}
}

type ComparableStats struct {
	Date               time.Time
	NumUsers           int
	NumCards           int
	NumCardsByNewUsers int
}

func ToComparableStats(stats []entity.NewUsersAndCardsStats) []ComparableStats {
	cs := make([]ComparableStats, len(stats))

	for i, stat := range stats {
		cs[i] = ComparableStats{
			Date:               stat.Date,
			NumUsers:           len(stat.Users),
			NumCards:           len(stat.Cards),
			NumCardsByNewUsers: stat.NumCardsByNewUsers,
		}
	}

	sort.Slice(cs, func(i, j int) bool {
		return cs[i].Date.Unix() < cs[j].Date.Unix()
	})

	return cs
}

package user

import (
	"aggregator/internal/dto"
	"context"
	"time"
)

type UserService interface {
	GetNewUsers(ctx context.Context, from time.Time, to time.Time) ([]dto.User, error)
}

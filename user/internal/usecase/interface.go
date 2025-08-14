package usecase

import (
	"context"
	"time"
	"user/internal/entity"
	"user/internal/repository"

	"github.com/google/uuid"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, user entity.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetUsers(ctx context.Context, filter repository.UserFilter) ([]entity.User, error)
	GetUsersBatch(ctx context.Context, limit, offset int) ([]entity.User, error)
	GetNewUsers(ctx context.Context, from time.Time, to time.Time) ([]entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

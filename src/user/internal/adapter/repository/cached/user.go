package repository

import (
	"context"
	"fmt"
	"user/internal/common/cache"
	"user/internal/entity"
	r "user/internal/repository"

	"github.com/google/uuid"
)

type CachedUserRepository struct {
	repo  r.UserRepository
	cache cache.Cache
}

func NewCachedUserRepository(repo r.UserRepository, cache cache.Cache) *CachedUserRepository {
	return &CachedUserRepository{
		repo:  repo,
		cache: cache,
	}
}

func (r *CachedUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	cacheKey := id.String()

	cachedData, err := r.cache.GetOrSet(ctx, cacheKey, func() (interface{}, error) {
		return r.repo.GetUserByID(ctx, id)
	})

	if err != nil {
		return nil, err
	}

	user := cachedData.(*entity.User)
	return user, nil
}

// TODO: GetUsers (with filter)
// TODO: GetNewUsers

func (r *CachedUserRepository) GetUsersBatch(ctx context.Context, limit, offset int) ([]entity.User, error) {
	cacheKey := fmt.Sprintf("l%do%d", limit, offset)

	cachedData, err := r.cache.GetOrSet(ctx, cacheKey, func() (interface{}, error) {
		return r.repo.GetUsersBatch(ctx, limit, offset)
	})

	if err != nil {
		return nil, err
	}

	users := cachedData.([]entity.User)
	return users, nil
}

func (r *CachedUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	return r.repo.CreateUser(ctx, user)
}

func (r *CachedUserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	return r.repo.UpdateUser(ctx, user)
}

func (r *CachedUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.repo.DeleteUser(ctx, id)
}

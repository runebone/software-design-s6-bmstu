package cache

import "context"

type Cache interface {
	GetOrSet(ctx context.Context, key string, fetchFunc func() (interface{}, error)) (interface{}, error)
}

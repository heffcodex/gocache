package gocache

import (
	"context"
	"time"
)

type RememberInt64Func func() (value int64, err error)
type RememberFloat64Func func() (value float64, err error)
type RememberStringFunc func() (value string, err error)

type Cacher interface {
	GetScan(ctx context.Context, key string, v interface{}) (err error)
	GetString(ctx context.Context, key string) (value string, err error)
	GetInt64(ctx context.Context, key string) (value int64, err error)
	GetFloat64(ctx context.Context, key string) (value float64, err error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) (err error)
	Del(ctx context.Context, key string) (err error)
	RememberInt64(ctx context.Context, key string, f RememberInt64Func, ttl time.Duration) (value int64, err error)
	RememberFloat64(ctx context.Context, key string, f RememberFloat64Func, ttl time.Duration) (value float64, err error)
	RememberString(ctx context.Context, key string, f RememberStringFunc, ttl time.Duration) (value string, err error)
	WithPrefix(prefix string) Cacher
}

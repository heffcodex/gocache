package gocache

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

var _ Cacher = new(NopCache)
var notFoundError = errors.New("not found")

type NopCache struct {
	prefix string
}

func NewNopCache() Cacher {
	return &NopCache{}
}

func (c *NopCache) GetScan(ctx context.Context, key string, v interface{}) (err error) {
	return notFoundError
}

func (c *NopCache) GetInt64(ctx context.Context, key string) (value int64, err error) {
	return 0, notFoundError
}

func (c *NopCache) GetFloat64(ctx context.Context, key string) (value float64, err error) {
	return 0, notFoundError
}

func (c *NopCache) GetString(ctx context.Context, key string) (value string, err error) {
	return "", notFoundError
}

func (c *NopCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) (err error) {
	return nil
}

func (c *NopCache) Del(ctx context.Context, key string) (err error) {
	return nil
}

func (c *NopCache) RememberInt64(ctx context.Context, key string, f RememberInt64Func, ttl time.Duration) (value int64, err error) {
	return f()
}

func (c *NopCache) RememberFloat64(ctx context.Context, key string, f RememberFloat64Func, ttl time.Duration) (value float64, err error) {
	return f()
}

func (c *NopCache) RememberString(ctx context.Context, key string, f RememberStringFunc, ttl time.Duration) (value string, err error) {
	return f()
}

func (c *NopCache) WithPrefix(prefix string) Cacher {
	return &NopCache{
		prefix: c.prefix + prefix,
	}
}

func (c *NopCache) key(key string) string {
	return c.prefix + key
}

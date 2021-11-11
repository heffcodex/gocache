package gocache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var _ Cacher = new(RedisCache)

type RedisCache struct {
	c      *redis.Client
	prefix string
}

func NewRedisCache(c *redis.Client) Cacher {
	return &RedisCache{c: c}
}

func (c *RedisCache) GetScan(ctx context.Context, key string, v interface{}) (err error) {
	return c.wrapErr(c.c.Get(ctx, c.key(key)).Scan(v))
}

func (c *RedisCache) GetInt64(ctx context.Context, key string) (value int64, err error) {
	value, err = c.c.Get(ctx, c.key(key)).Int64()
	err = c.wrapErr(err)

	return
}

func (c *RedisCache) GetFloat64(ctx context.Context, key string) (value float64, err error) {
	value, err = c.c.Get(ctx, c.key(key)).Float64()
	err = c.wrapErr(err)

	return
}

func (c *RedisCache) GetString(ctx context.Context, key string) (value string, err error) {
	value, err = c.c.Get(ctx, c.key(key)).Result()
	err = c.wrapErr(err)

	return
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) (err error) {
	return c.c.Set(ctx, c.key(key), value, ttl).Err()
}

func (c *RedisCache) Del(ctx context.Context, key string) (err error) {
	return c.wrapErr(c.c.Del(ctx, key).Err())
}

func (c *RedisCache) RememberInt64(ctx context.Context, key string, f RememberInt64Func, ttl time.Duration) (value int64, err error) {
	value, err = c.c.Get(ctx, c.key(key)).Int64()
	if err != nil {
		value, err = f()
		if err != nil {
			return
		}

		err = c.Set(ctx, key, value, ttl)
		if err != nil {
			err = WrapRememberError(err, c.key(key))
		}
	}

	return
}

func (c *RedisCache) RememberFloat64(ctx context.Context, key string, f RememberFloat64Func, ttl time.Duration) (value float64, err error) {
	value, err = c.c.Get(ctx, c.key(key)).Float64()
	if err != nil {
		value, err = f()
		if err != nil {
			return
		}

		err = c.Set(ctx, key, value, ttl)
		if err != nil {
			err = WrapRememberError(err, c.key(key))
		}
	}

	return
}

func (c *RedisCache) RememberString(ctx context.Context, key string, f RememberStringFunc, ttl time.Duration) (value string, err error) {
	value, err = c.c.Get(ctx, c.key(key)).Result()
	if err != nil {
		value, err = f()
		if err != nil {
			return
		}

		err = c.Set(ctx, key, value, ttl)
		if err != nil {
			err = WrapRememberError(err, c.key(key))
		}
	}

	return
}

func (c *RedisCache) WithPrefix(prefix string) Cacher {
	return &RedisCache{
		c:      c.c,
		prefix: c.prefix + prefix,
	}
}

func (c *RedisCache) key(key string) string {
	return c.prefix + key
}

func (c *RedisCache) wrapErr(err error) error {
	if err == redis.Nil {
		err = &NilError{err: err}
	}

	return err
}

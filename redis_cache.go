package gocache

import (
	"time"

	"github.com/go-redis/redis/v7"
)

var _ Cacher = new(RedisCache)

type RedisCache struct {
	c      *redis.Client
	prefix string
}

func NewRedisCache(c *redis.Client) Cacher {
	return &RedisCache{c: c}
}

func (c *RedisCache) GetScan(key string, v interface{}) (err error) {
	return c.wrapErr(c.c.Get(c.key(key)).Scan(v))
}

func (c *RedisCache) GetInt64(key string) (value int64, err error) {
	value, err = c.c.Get(c.key(key)).Int64()
	err = c.wrapErr(err)

	return
}

func (c *RedisCache) GetString(key string) (value string, err error) {
	value, err = c.c.Get(c.key(key)).Result()
	err = c.wrapErr(err)

	return
}

func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) (err error) {
	return c.c.Set(c.key(key), value, ttl).Err()
}

func (c *RedisCache) Del(key string) (err error) {
	return c.wrapErr(c.c.Del(key).Err())
}

func (c *RedisCache) RememberInt64(key string, f RememberInt64Func, ttl time.Duration) (value int64, err error) {
	value, err = c.c.Get(c.key(key)).Int64()
	if err != nil {
		value, err = f()
		if err != nil {
			return
		}

		err = c.Set(key, value, ttl)
		if err != nil {
			err = WrapRememberError(err, c.key(key))
		}
	}

	return
}

func (c *RedisCache) RememberString(key string, f RememberStringFunc, ttl time.Duration) (value string, err error) {
	value, err = c.c.Get(c.key(key)).Result()
	if err != nil {
		value, err = f()
		if err != nil {
			return
		}

		err = c.Set(key, value, ttl)
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
		err = nil
	}

	return err
}

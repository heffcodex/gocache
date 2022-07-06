package gocache

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var _ Cacher = new(MemoryCache)

type memoryCacheValue struct {
	v   interface{}
	ttl int64
}

type MemoryCache struct {
	l        *sync.Mutex
	lastTick *int64
	keys     map[string]*memoryCacheValue
	prefix   string
}

func NewMemoryCache() *MemoryCache {
	now := time.Now().Unix()

	return &MemoryCache{
		l:        new(sync.Mutex),
		lastTick: &now,
		keys:     make(map[string]*memoryCacheValue),
	}
}

func (c *MemoryCache) GetScan(ctx context.Context, key string, v interface{}) (err error) {
	c.l.Lock()
	defer c.l.Unlock()

	c.tick(ctx)

	vv := reflect.ValueOf(v)
	if vv.Kind() != reflect.Ptr {
		return errors.New("value must be pointer type")
	}

	value, ok := c.keys[c.key(key)]
	if !ok {
		return new(NilError)
	}

	valueV := reflect.ValueOf(value.v)

	if !vv.Elem().Type().AssignableTo(valueV.Type()) {
		return errors.New("unassignable type")
	}

	vv.Elem().Set(valueV)

	return
}

func (c *MemoryCache) GetString(ctx context.Context, key string) (value string, err error) {
	err = c.GetScan(ctx, key, &value)
	return
}

func (c *MemoryCache) GetInt64(ctx context.Context, key string) (value int64, err error) {
	err = c.GetScan(ctx, key, &value)
	return
}

func (c *MemoryCache) GetFloat64(ctx context.Context, key string) (value float64, err error) {
	err = c.GetScan(ctx, key, &value)
	return
}

func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) (err error) {
	c.l.Lock()
	defer c.l.Unlock()

	c.tick(ctx)

	k := c.key(key)
	c.del(ctx, k)

	v := &memoryCacheValue{v: value, ttl: 0}
	if ttl > 0 {
		v.ttl = time.Now().Add(ttl).Unix()
	}

	c.keys[k] = v

	return
}

func (c *MemoryCache) Del(ctx context.Context, key string) (err error) {
	c.l.Lock()
	defer c.l.Unlock()

	c.tick(ctx)
	c.del(ctx, c.key(key))

	return
}

func (c *MemoryCache) del(_ context.Context, rawKey string) {
	_, ok := c.keys[rawKey]
	if ok {
		delete(c.keys, rawKey)
	}
}

func (c *MemoryCache) RememberInt64(ctx context.Context, key string, f RememberInt64Func, ttl time.Duration) (value int64, err error) {
	value, err = c.GetInt64(ctx, key)
	if _, ok := err.(*NilError); ok {
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

func (c *MemoryCache) RememberFloat64(ctx context.Context, key string, f RememberFloat64Func, ttl time.Duration) (value float64, err error) {
	value, err = c.GetFloat64(ctx, key)
	if _, ok := err.(*NilError); ok {
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

func (c *MemoryCache) RememberString(ctx context.Context, key string, f RememberStringFunc, ttl time.Duration) (value string, err error) {
	value, err = c.GetString(ctx, key)
	if _, ok := err.(*NilError); ok {
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

func (c *MemoryCache) WithPrefix(prefix string) Cacher {
	return &MemoryCache{
		l:        c.l,
		lastTick: c.lastTick,
		keys:     c.keys,
		prefix:   c.prefix + prefix,
	}
}

func (c *MemoryCache) tick(ctx context.Context) {
	now := time.Now().Unix()

	if now == *c.lastTick {
		return
	}

	c.lastTick = &now

	for k, v := range c.keys {
		if v.ttl > 0 && now > v.ttl {
			c.del(ctx, k)
		}
	}
}

func (c *MemoryCache) key(key string) string {
	return c.prefix + key
}

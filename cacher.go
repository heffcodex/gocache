package gocache

import (
	"time"
)

type RememberInt64Func func() (value int64, err error)
type RememberStringFunc func() (value string, err error)

type Cacher interface {
	GetScan(key string, v interface{}) (err error)
	GetString(key string) (value string, err error)
	GetInt64(key string) (value int64, err error)
	Set(key string, value interface{}, ttl time.Duration) (err error)
	Del(key string) (err error)
	RememberInt64(key string, f RememberInt64Func, ttl time.Duration) (value int64, err error)
	RememberString(key string, f RememberStringFunc, ttl time.Duration) (value string, err error)
	WithPrefix(prefix string) Cacher
}

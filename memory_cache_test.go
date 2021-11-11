package gocache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMemoryCache_WithPrefix(t *testing.T) {
	ctx := context.TODO()
	c := NewMemoryCache()
	cpref := c.WithPrefix("p:")

	require.NoError(t, c.Set(ctx, "k", "nopref", 0))
	require.NoError(t, cpref.Set(ctx, "k", "pref", 0))

	v, err := c.GetString(ctx, "k")
	require.NoError(t, err)
	require.Equal(t, "nopref", v)

	v, err = cpref.GetString(ctx, "k")
	require.NoError(t, err)
	require.Equal(t, "pref", v)
}

func TestMemoryCache_GetScan(t *testing.T) {
	ctx := context.TODO()
	c := NewMemoryCache()

	require.NoError(t, c.Set(ctx, "k", "v", 0))

	t.Run("positive", func(t *testing.T) {
		v := ""
		require.NoError(t, c.GetScan(ctx, "k", &v))
		require.Equal(t, "v", v)
	})

	t.Run("errors", func(t *testing.T) {
		v := ""
		require.Error(t, c.GetScan(ctx, "k", v))

		vv := 0
		require.Error(t, c.GetScan(ctx, "k", &vv))
	})
}

func TestMemoryCache_Set(t *testing.T) {
	ctx := context.TODO()
	c := NewMemoryCache()

	t.Run("ttl=2", func(t *testing.T) {
		require.NoError(t, c.Set(ctx, "k", "ttl2", 2))

		v, err := c.GetString(ctx, "k")
		require.NoError(t, err)
		require.Equal(t, "ttl2", v)

		time.Sleep(3 * time.Second)

		v, err = c.GetString(ctx, "k")
		require.EqualError(t, err, new(NilError).Error())
	})

	t.Run("ttl=0", func(t *testing.T) {
		require.NoError(t, c.Set(ctx, "k", "ttl0", 0))

		v, err := c.GetString(ctx, "k")
		require.NoError(t, err)
		require.Equal(t, "ttl0", v)

		time.Sleep(3 * time.Second)

		v, err = c.GetString(ctx, "k")
		require.NoError(t, err)
		require.Equal(t, "ttl0", v)
	})

}

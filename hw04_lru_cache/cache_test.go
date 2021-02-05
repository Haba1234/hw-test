package hw04_lru_cache //nolint:golint,stylecheck

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c, _ := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("capacity <= 0", func(t *testing.T) {
		c, err := NewCache(-1)
		require.Nil(t, c)
		require.Error(t, err)

		c, err = NewCache(0)
		require.Nil(t, c)
		require.Error(t, err)
	})

	t.Run("simple", func(t *testing.T) {
		c, _ := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		// Write me
		c, _ := NewCache(3)
		wasInCache := c.Set("Clear", 100)
		require.False(t, wasInCache)

		c.Clear()
		val, ok := c.Get("Clear")
		require.False(t, ok)
		require.Nil(t, val)

		for _, v := range [...]int{1, 2, 3} {
			wasInCache = c.Set(Key(fmt.Sprint(v)), v)
		} // [1, 2, 3]
		wasInCache = c.Set("4", 4) // [2, 3, 4]
		val, ok = c.Get("1")
		require.False(t, ok)
		require.Nil(t, val)

		rand.Seed(1)
		for _, v := range [...]int{2, 3, 3, 4, 3, 3, 4, 2, 4, 3} {
			if rand.Int() > 0 {
				wasInCache = c.Set(Key(fmt.Sprint(v)), v)
			} else {
				val, ok = c.Get(Key(fmt.Sprint(v)))
			}
		} // [3, 4, 2]
		wasInCache = c.Set("5", 5) // [5, 3, 4]
		val, ok = c.Get("2")
		require.False(t, ok)
		require.Nil(t, val)

	})
}

func TestCacheMultithreading(t *testing.T) {
	//t.Skip() // NeedRemove if task with asterisk completed

	c, _ := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

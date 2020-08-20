package ratelimit

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRateLimit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// возьмем небольшой лимит, чтобы время регенерации токенов было относительно большим
	limit := 3
	rd, _ := time.ParseDuration("2s")

	key1 := "key1"
	key2 := "key2"
	t.Run("simple case", func(t *testing.T) {
		rLimit := NewRateLimit(ctx, limit, rd)
		for i := 0; i < limit; i++ {
			require.True(t, rLimit.Allow(key1))
		}
		require.False(t, rLimit.Allow(key1))
		require.True(t, rLimit.Allow(key2))
	})
	t.Run("with regeneration", func(t *testing.T) {
		rLimit := NewRateLimit(ctx, limit, rd)
		for i := 0; i < limit; i++ {
			require.True(t, rLimit.Allow(key1))
			time.Sleep(500 * time.Millisecond)
		}
		time.Sleep(500 * time.Millisecond)
		require.True(t, rLimit.Allow(key1))
	})
	t.Run("clearing", func(t *testing.T) {
		rLimit := NewRateLimit(ctx, limit, rd)
		for i := 0; i < limit; i++ {
			require.True(t, rLimit.Allow(key1))
		}
		rLimit.Clear(key1)
		require.True(t, rLimit.Allow(key1))
	})
	t.Run("goroutine safe", func(t *testing.T) {
		rLimit := NewRateLimit(ctx, limit, rd)

		ch := make(chan bool)
		wg := sync.WaitGroup{}
		goroutineCnt := 2
		wg.Add(goroutineCnt)
		for j := 0; j < goroutineCnt; j++ {
			go func() {
				for i := 0; i < limit; i++ {
					ch <- rLimit.Allow(key1)
				}
				wg.Done()
			}()
		}
		go func() {
			wg.Wait()
			close(ch)
		}()
		var trueCnt int
		var falseCnt int
		for res := range ch {
			if res {
				trueCnt++
			} else {
				falseCnt++
			}
		}
		require.Equal(t, limit, trueCnt)
		require.Equal(t, limit, falseCnt)
	})
}

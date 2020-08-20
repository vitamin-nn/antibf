package ratelimit

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	cleanupInterval = 5 * time.Minute
)

type RateLimit struct {
	mutex *sync.Mutex
	limit int
	rd    time.Duration
	data  map[string]*RateLimiter
}

type RateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewRateLimit
// limit - maximum amount of request per rd interval
// rd - rate duration.
func NewRateLimit(ctx context.Context, limit int, rd time.Duration) *RateLimit {
	r := &RateLimit{
		mutex: new(sync.Mutex),
		limit: limit,
		rd:    rd,
		data:  make(map[string]*RateLimiter),
	}
	go r.cleanupOld(ctx)

	return r
}

func (rl *RateLimit) Allow(key string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	r := rl.getLimiter(key)
	r.lastSeen = time.Now()

	return r.limiter.Allow()
}

func (rl *RateLimit) Clear(key string) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	delete(rl.data, key)
}

func (rl *RateLimit) getLimiter(key string) *RateLimiter {
	r, exists := rl.data[key]
	if !exists {
		r = &RateLimiter{
			limiter: rate.NewLimiter(rate.Every(time.Duration(rl.rd.Nanoseconds()/int64(rl.limit))), rl.limit),
		}
		rl.data[key] = r
	}

	return r
}

func (rl *RateLimit) cleanupOld(ctx context.Context) {
	ticker := time.NewTicker(cleanupInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rl.mutex.Lock()
			for k, v := range rl.data {
				if time.Since(v.lastSeen) > cleanupInterval {
					delete(rl.data, k)
				}
			}
			rl.mutex.Unlock()
		}
	}
}

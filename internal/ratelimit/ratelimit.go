package ratelimit

import (
	"time"
)

type RateLimitRepo interface {
	GetRate(key string, startTimestamp int64) int
	Increment(key string, timestamp int64)
	Clear(key string)
}

type RateLimit struct {
	repo RateLimitRepo
	cap  int   // capacity
	rd   int64 // rate duration
}

func NewRateLimit(rlRepo RateLimitRepo, cap int, rd int64) *RateLimit {
	rl := new(RateLimit)
	rl.repo = rlRepo
	rl.cap = cap
	rl.rd = rd
	return rl
}

func (rl *RateLimit) IsAllow(key string) bool {
	t := time.Now()
	ts := t.Unix()
	rl.repo.Increment(key, ts)
	startTs := ts - rl.rd
	rate := rl.repo.GetRate(key, startTs)
	return rate < rl.cap
}

func (rl *RateLimit) Clear(key string) {
	rl.Clear(key)
}

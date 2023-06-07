package ratelimiter

import (
	"sync"
	"time"
)

// A rate limiter with a manually specified maximum; throttles requests when they exceed the maximum or when Backoff() is called
type CappedRateLimiter struct {
	maxPerMinute     int
	currentPerMinute int
	currentLimiter   *time.Ticker
	successStreak    int
	lock             *sync.RWMutex
}

func NewCappedRateLimiter(requestsPerMinute int) *CappedRateLimiter {
	return &CappedRateLimiter{
		maxPerMinute:     requestsPerMinute,
		currentPerMinute: requestsPerMinute,
		currentLimiter:   time.NewTicker(time.Minute / time.Duration(requestsPerMinute)),
		successStreak:    0,
		lock:             &sync.RWMutex{},
	}
}

func (rl *CappedRateLimiter) LimitRate() {
	<-rl.currentLimiter.C
}

func (rl *CappedRateLimiter) Success() {
	rl.lock.Lock()
	defer rl.lock.Unlock()
	// linear growth
	rl.successStreak++
	if (rl.successStreak % 10) == 0 {
		if rl.currentPerMinute < rl.maxPerMinute {
			rl.currentPerMinute++
			rl.currentLimiter.Reset(time.Minute / time.Duration(rl.currentPerMinute))
		}
	}
}

func (rl *CappedRateLimiter) Backoff() {
	rl.lock.Lock()
	defer rl.lock.Unlock()
	// exponential backoff
	rl.successStreak = 0
	rl.currentPerMinute = rl.currentPerMinute / 2
	if rl.currentPerMinute < 1 {
		rl.currentPerMinute = 1
	}
	rl.currentLimiter.Reset(time.Minute / time.Duration(rl.currentPerMinute))
}

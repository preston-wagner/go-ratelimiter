package ratelimiter

import (
	"math"
	"sync"
	"time"
)

// A rate limiter with a manually specified maximum; throttles when the maximum is exceeded or when Backoff() is called
type CappedRateLimiter struct {
	// set by contstructor args
	maxPerInterval int           // maximum number of requests per interval (ex: 7 for "7 requests per minute")
	interval       time.Duration // (ex: time.Minute for "7 requests per minute")
	streakLength   int           // number of successful requests that must be made before stepping up the current rate
	backoffFactor  float64       // amount to decrease current rate when Backoff() is called (ex: 0.75 for "reduce limit by 25%")
	// automatically set
	currentPerMinute int
	currentLimiter   *time.Ticker
	successStreak    int
	lock             *sync.RWMutex
}

func NewCappedRateLimiter(maxPerInterval int, interval time.Duration, streakLength int, backoffFactor float64) *CappedRateLimiter {
	if maxPerInterval <= 0 {
		panic("NewCappedRateLimiter maxPerInterval must be > 0")
	}
	if interval <= 0 {
		panic("NewCappedRateLimiter interval must be > 0")
	}
	if streakLength <= 0 {
		panic("NewCappedRateLimiter streakLength must be > 0")
	}
	if backoffFactor <= 0 {
		panic("NewCappedRateLimiter backoffFactor must be > 0")
	}
	if backoffFactor >= 1 {
		panic("NewCappedRateLimiter backoffFactor must be < 1")
	}

	return &CappedRateLimiter{
		maxPerInterval: maxPerInterval,
		interval:       interval,
		streakLength:   streakLength,
		backoffFactor:  backoffFactor,

		currentPerMinute: maxPerInterval,
		currentLimiter:   time.NewTicker(interval / time.Duration(maxPerInterval)),
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
	if (rl.successStreak % rl.streakLength) == 0 {
		if rl.currentPerMinute < rl.maxPerInterval {
			rl.currentPerMinute++
			rl.currentLimiter.Reset(rl.interval / time.Duration(rl.currentPerMinute))
		}
	}
}

func (rl *CappedRateLimiter) Backoff() {
	rl.lock.Lock()
	defer rl.lock.Unlock()
	// exponential backoff
	rl.successStreak = 0
	rl.currentPerMinute = int(math.Floor(float64(rl.currentPerMinute) * rl.backoffFactor))
	if rl.currentPerMinute < 1 {
		rl.currentPerMinute = 1
	}
	rl.currentLimiter.Reset(rl.interval / time.Duration(rl.currentPerMinute))
}

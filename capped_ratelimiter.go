package ratelimiter

import (
	"sync"
	"time"
)

type cappedRateLimiter struct {
	// set by constructor args
	maxPerInterval int           // maximum number of requests per interval (ex: 7 for "7 requests per minute")
	interval       time.Duration // (ex: time.Minute for "7 requests per minute")
	streakLength   int           // number of successful requests that must be made before stepping up the current rate
	backoffFactor  float64       // amount to decrease current rate when Backoff() is called (ex: 0.75 for "reduce limit by 25%")
	// automatically set
	currentPerInterval int
	currentLimiter     *time.Ticker
	successStreak      int
	lock               *sync.RWMutex
}

// A rate limiter with a manually specified maximum; throttles when the maximum is exceeded or when Backoff() is called
func NewCappedRateLimiter(maxPerInterval int, interval time.Duration, streakLength int, backoffFactor float64) RateLimiter {
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

	return &cappedRateLimiter{
		maxPerInterval: maxPerInterval,
		interval:       interval,
		streakLength:   streakLength,
		backoffFactor:  backoffFactor,

		currentPerInterval: maxPerInterval,
		currentLimiter:     time.NewTicker(interval / time.Duration(maxPerInterval)),
		successStreak:      0,
		lock:               &sync.RWMutex{},
	}
}

func (rl *cappedRateLimiter) LimitRate() {
	<-rl.currentLimiter.C
}

func (rl *cappedRateLimiter) Success() {
	rl.lock.Lock()
	defer rl.lock.Unlock()
	// linear growth
	rl.successStreak++
	if (rl.successStreak % rl.streakLength) == 0 {
		if rl.currentPerInterval < rl.maxPerInterval {
			rl.currentPerInterval++
			rl.currentLimiter.Reset(rl.interval / time.Duration(rl.currentPerInterval))
		}
	}
}

func (rl *cappedRateLimiter) Backoff() {
	rl.lock.Lock()
	defer rl.lock.Unlock()
	// exponential backoff
	rl.currentPerInterval = int(float64(rl.currentPerInterval) * rl.backoffFactor)
	if rl.currentPerInterval < 1 {
		rl.currentPerInterval = 1
	}
	rl.currentLimiter.Reset(rl.interval / time.Duration(rl.currentPerInterval))
	rl.successStreak = 0
}

func (rl *cappedRateLimiter) GetCurrentRate() *time.Duration {
	rl.lock.RLock()
	defer rl.lock.RUnlock()
	rate := rl.interval / time.Duration(rl.currentPerInterval)
	return &rate
}

package ratelimiter

import (
	"sync"
	"time"
)

// A rate limiter that "feels out" the limits of the api being accessed; begins throttling when Backoff() is called
type AutoRateLimiter struct {
	// set by constructor args
	streakLength  int     // number of successful requests that must be made before stepping up the current rate
	backoffFactor float64 // amount to decrease current rate when Backoff() is called (ex: 0.75 for "reduce limit by 25%")
	// automatically set
	startTime        time.Time
	currentPerMinute int
	currentLimiter   *time.Ticker
	successStreak    int
	lock             *sync.RWMutex
}

// A rate limiter that "feels out" the limits of the api being accessed; begins throttling when Backoff() is called
func NewAutoRateLimiter(streakLength int, backoffFactor float64) *AutoRateLimiter {
	if streakLength <= 0 {
		panic("NewAutoRateLimiter streakLength must be > 0")
	}
	if backoffFactor <= 0 {
		panic("NewAutoRateLimiter backoffFactor must be > 0")
	}
	if backoffFactor >= 1 {
		panic("NewAutoRateLimiter backoffFactor must be < 1")
	}

	return &AutoRateLimiter{
		streakLength:  streakLength,
		backoffFactor: backoffFactor,

		startTime:     time.Now(),
		successStreak: 0,
		lock:          &sync.RWMutex{},
	}
}

func (rl *AutoRateLimiter) LimitRate() {
	rl.lock.RLock()
	defer rl.lock.RUnlock()
	// if currentLimiter is nil, then Backoff() has not been called at all yet
	if rl.currentLimiter != nil {
		<-rl.currentLimiter.C
	}
}

func (rl *AutoRateLimiter) Success() {
	rl.lock.Lock()
	defer rl.lock.Unlock()
	// linear growth
	rl.successStreak++
	if rl.currentLimiter != nil {
		if (rl.successStreak % rl.streakLength) == 0 {
			rl.currentPerMinute++
			rl.currentLimiter.Reset(time.Minute / time.Duration(rl.currentPerMinute))
		}
	}
}

func (rl *AutoRateLimiter) Backoff() {
	rl.lock.Lock()
	defer rl.lock.Unlock()
	if rl.currentLimiter == nil {
		rl.currentLimiter = time.NewTicker(time.Minute) // duration doesn't matter since it's overwritten below anyways
		durationSinceStart := time.Since(rl.startTime)
		minutesSinceStart := float64(durationSinceStart) / float64(time.Minute)
		if minutesSinceStart > 0 {
			successesPerMinute := float64(rl.successStreak) / minutesSinceStart
			rl.currentPerMinute = int(successesPerMinute * rl.backoffFactor)
		} else {
			rl.currentPerMinute = 1
		}
	} else {
		rl.currentPerMinute = int(float64(rl.currentPerMinute) * rl.backoffFactor)
	}
	if rl.currentPerMinute < 1 {
		rl.currentPerMinute = 1
	}
	rl.currentLimiter.Reset(time.Minute / time.Duration(rl.currentPerMinute))
	rl.successStreak = 0
}

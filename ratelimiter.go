package ratelimiter

import "time"

type RateLimiter interface {
	LimitRate()
	Success()
	Backoff()
	GetCurrentRate() *time.Duration // 7 per minute = time.Minute / 7, unlimited = nil
}

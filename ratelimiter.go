package ratelimiter

type RateLimiter interface {
	LimitRate()
	Success()
	Backoff()
}

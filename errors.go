package ratelimiter

import "fmt"

type RateLimitExceeded struct {
	Err error
}

func (rle RateLimitExceeded) Error() string {
	return "Rate limit exceeded: " + rle.Err.Error()
}

func (rle RateLimitExceeded) Unwrap() error {
	return rle.Err
}

type RetriesExceeded struct {
	Retries int
}

func (re RetriesExceeded) Error() string {
	return fmt.Sprintf("Retries exceeded: %v", re.Retries)
}

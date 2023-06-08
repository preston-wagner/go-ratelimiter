package ratelimiter

type RateLimitExceeded struct {
	Err error
}

func (rle RateLimitExceeded) Error() string {
	return "Rate limit exceeded: " + rle.Err.Error()
}

func (rle RateLimitExceeded) Unwrap() error {
	return rle.Err
}

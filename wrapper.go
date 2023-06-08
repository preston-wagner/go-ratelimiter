package ratelimiter

import "errors"

func RateLimitedCall[ARG_TYPE any, RETURN_TYPE any](rl RateLimiter, wrapped func(arg ARG_TYPE) (RETURN_TYPE, error), arg ARG_TYPE) (RETURN_TYPE, error) {
	for {
		rl.LimitRate()
		result, err := wrapped(arg)
		if err == nil {
			rl.Success()
			return result, nil
		} else {
			var rateLimitErr RateLimitExceeded
			if errors.As(err, &rateLimitErr) {
				rl.Backoff()
			} else {
				return result, err
			}
		}
	}
}

package ratelimiter

import "errors"

type BasicLimitableCall[ARG_TYPE any, RETURN_TYPE any] func(arg ARG_TYPE) (RETURN_TYPE, error)

func RateLimitedCall[ARG_TYPE any, RETURN_TYPE any](rl RateLimiter, wrapped BasicLimitableCall[ARG_TYPE, RETURN_TYPE], arg ARG_TYPE) (RETURN_TYPE, error) {
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

func WrapWithLimit[ARG_TYPE any, RETURN_TYPE any](rl RateLimiter, wrapped BasicLimitableCall[ARG_TYPE, RETURN_TYPE]) BasicLimitableCall[ARG_TYPE, RETURN_TYPE] {
	return func(arg ARG_TYPE) (RETURN_TYPE, error) {
		return RateLimitedCall(rl, wrapped, arg)
	}
}

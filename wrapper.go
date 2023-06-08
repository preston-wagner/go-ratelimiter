package ratelimiter

import (
	"errors"

	"github.com/preston-wagner/unicycle"
)

type BasicLimitableCall[ARG_TYPE any, RETURN_TYPE any] func(arg ARG_TYPE) (RETURN_TYPE, error)

// RateLimitedCall uses the provided ratelimiter to slow the flow of requests, but does not retry on RateLimitExceeded
func RateLimitedCall[ARG_TYPE any, RETURN_TYPE any](rl RateLimiter, wrapped BasicLimitableCall[ARG_TYPE, RETURN_TYPE], arg ARG_TYPE) (RETURN_TYPE, error) {
	rl.LimitRate()
	result, err := wrapped(arg)
	var rateLimitErr RateLimitExceeded
	if (err == nil) || (!errors.As(err, &rateLimitErr)) {
		rl.Success()
	} else {
		rl.Backoff()
	}
	return result, err
}

// RateLimitedRetryCall is like RateLimitedCall, but on RateLimitExceeded attempts the request again, up to a given maximum of retries
func RateLimitedRetryCall[ARG_TYPE any, RETURN_TYPE any](rl RateLimiter, retries int, wrapped BasicLimitableCall[ARG_TYPE, RETURN_TYPE], arg ARG_TYPE) (RETURN_TYPE, error) {
	for i := 0; i < retries; i++ {
		result, err := RateLimitedCall(rl, wrapped, arg)
		var rateLimitErr RateLimitExceeded
		if (err == nil) || (!errors.As(err, &rateLimitErr)) {
			return result, err
		}
	}
	return unicycle.ZeroValue[RETURN_TYPE](), RetriesExceeded{Retries: retries}
}

func WrapWithLimit[ARG_TYPE any, RETURN_TYPE any](rl RateLimiter, wrapped BasicLimitableCall[ARG_TYPE, RETURN_TYPE]) BasicLimitableCall[ARG_TYPE, RETURN_TYPE] {
	return func(arg ARG_TYPE) (RETURN_TYPE, error) {
		return RateLimitedCall(rl, wrapped, arg)
	}
}

func WrapWithRetryLimit[ARG_TYPE any, RETURN_TYPE any](rl RateLimiter, retries int, wrapped BasicLimitableCall[ARG_TYPE, RETURN_TYPE]) BasicLimitableCall[ARG_TYPE, RETURN_TYPE] {
	return func(arg ARG_TYPE) (RETURN_TYPE, error) {
		return RateLimitedRetryCall(rl, retries, wrapped, arg)
	}
}

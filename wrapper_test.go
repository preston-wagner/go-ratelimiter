package ratelimiter

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type rateLimitTester struct {
	lastCall     *time.Time
	minFrequency time.Duration
}

func (rlt *rateLimitTester) callWithSpeedLimit(key string) (int, error) {
	now := time.Now()
	if rlt.lastCall != nil {
		if rlt.minFrequency > time.Since(*rlt.lastCall) {
			return 0, RateLimitExceeded{Err: errors.New("going too fast!")}
		}
	}
	rlt.lastCall = &now
	return 7, nil
}

func (rlt *rateLimitTester) failWithSpeedLimit(key string) (int, error) {
	now := time.Now()
	if rlt.lastCall != nil {
		if rlt.lastCall.Add(rlt.minFrequency).After(now) {
			return 0, RateLimitExceeded{Err: errors.New("going too fast!")}
		}
	}
	rlt.lastCall = &now
	return 0, errors.New("A different error!")
}

func TestRateLimitedRetryCall(t *testing.T) {
	rlt := rateLimitTester{minFrequency: time.Second}

	limiter := NewCappedRateLimiter(120, time.Minute, 1, 0.5)

	for i := 0; i < 5; i++ {
		val, err := RateLimitedRetryCall(limiter, 5, rlt.callWithSpeedLimit, "lorem ipsum")
		if err != nil {
			t.Error(err)
		}
		if val != 7 {
			t.Error("unexpected value returned")
		}
	}

	for i := 0; i < 5; i++ {
		_, err := RateLimitedRetryCall(limiter, 5, rlt.failWithSpeedLimit, "dolor sit amet")
		if err == nil {
			t.Error("expected error not returned!")
		}
		assert.Equal(t, err.Error(), "A different error!")
	}
}

func TestWrapWithRetryLimit(t *testing.T) {
	rlt := rateLimitTester{minFrequency: time.Second}

	limiter := NewCappedRateLimiter(120, time.Minute, 1, 0.5)
	wrappedCallWithSpeedLimit := WrapWithRetryLimit(limiter, 5, rlt.callWithSpeedLimit)

	for i := 0; i < 5; i++ {
		val, err := wrappedCallWithSpeedLimit("lorem ipsum")
		if err != nil {
			t.Error(err)
		}
		if val != 7 {
			t.Error("unexpected value returned")
		}
	}

	wrappedFailWithSpeedLimit := WrapWithRetryLimit(limiter, 5, rlt.failWithSpeedLimit)
	for i := 0; i < 5; i++ {
		_, err := wrappedFailWithSpeedLimit("dolor sit amet")
		if err == nil {
			t.Error("expected error not returned!")
		}
		assert.Equal(t, err.Error(), "A different error!")
	}
}

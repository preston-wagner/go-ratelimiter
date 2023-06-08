package ratelimiter

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type rateLimitTester struct {
	lastCall *time.Time
	maxRate  time.Duration
}

func (rlt *rateLimitTester) callWithSpeedLimit(key string) (int, error) {
	now := time.Now()
	if rlt.lastCall != nil {
		if rlt.lastCall.Add(rlt.maxRate).After(now) {
			return 0, RateLimitExceeded{Err: errors.New("going too fast!")}
		}
	}
	rlt.lastCall = &now
	return 7, nil
}

func (rlt *rateLimitTester) failWithSpeedLimit(key string) (int, error) {
	now := time.Now()
	if rlt.lastCall != nil {
		if rlt.lastCall.Add(rlt.maxRate).After(now) {
			return 0, RateLimitExceeded{Err: errors.New("going too fast!")}
		}
	}
	rlt.lastCall = &now
	return 0, errors.New("A different error!")
}

func TestRateLimitedCall(t *testing.T) {
	rlt := rateLimitTester{maxRate: time.Second}

	limiter := NewCappedRateLimiter(120, time.Minute, 1, 0.5)

	for i := 0; i < 5; i++ {
		val, err := RateLimitedCall(limiter, rlt.callWithSpeedLimit, "lorem ipsum")
		if err != nil {
			t.Error(err)
		}
		if val != 7 {
			t.Error("unexpected value returned")
		}
	}

	for i := 0; i < 5; i++ {
		_, err := RateLimitedCall(limiter, rlt.failWithSpeedLimit, "dolor sit amet")
		if err == nil {
			t.Error("expected error not returned!")
		}
		assert.Equal(t, err.Error(), "A different error!")
	}
}

func TestWrapWithLimit(t *testing.T) {
	rlt := rateLimitTester{maxRate: time.Second}

	limiter := NewCappedRateLimiter(120, time.Minute, 1, 0.5)
	wrappedCallWithSpeedLimit := WrapWithLimit(limiter, rlt.callWithSpeedLimit)

	for i := 0; i < 5; i++ {
		val, err := wrappedCallWithSpeedLimit("lorem ipsum")
		if err != nil {
			t.Error(err)
		}
		if val != 7 {
			t.Error("unexpected value returned")
		}
	}

	wrappedFailWithSpeedLimit := WrapWithLimit(limiter, rlt.failWithSpeedLimit)
	for i := 0; i < 5; i++ {
		_, err := wrappedFailWithSpeedLimit("dolor sit amet")
		if err == nil {
			t.Error("expected error not returned!")
		}
		assert.Equal(t, err.Error(), "A different error!")
	}
}

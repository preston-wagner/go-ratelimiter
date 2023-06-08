package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAutoRateLimit(t *testing.T) {
	const secondsToTest = 15

	rl := NewAutoRateLimiter(2, .9)
	rlt := rateLimitTester{minFrequency: time.Second}

	keepLooping := true
	go func() {
		for keepLooping {
			RateLimitedRetryCall(rl, 5, rlt.callWithSpeedLimit, "lorem ipsum")
			time.Sleep(time.Second / 100)
		}
	}()

	time.Sleep(time.Second * secondsToTest)
	keepLooping = false

	assert.GreaterOrEqual(t, rl.currentPerMinute, 45)
	assert.LessOrEqual(t, rl.currentPerMinute, 61)
}

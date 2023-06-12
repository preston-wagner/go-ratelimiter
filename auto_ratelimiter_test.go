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

	currentRate := *rl.GetCurrentRate()
	currentPerMinute := float64(time.Minute / currentRate)
	assert.GreaterOrEqual(t, currentPerMinute, 45.0)
	assert.LessOrEqual(t, currentPerMinute, 61.0)
}

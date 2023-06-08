package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCappedRateLimit(t *testing.T) {
	const secondsToTest = 8
	rl := NewCappedRateLimiter(60, time.Minute, 10, .5)

	tickCount := 0
	keepLooping := true

	tick := func() {
		rl.LimitRate()
		tickCount++
	}

	for i := 0; i < 5; i++ {
		go func() {
			for keepLooping {
				tick()
			}
		}()
	}

	time.Sleep(time.Second * secondsToTest)
	keepLooping = false

	assert.GreaterOrEqual(t, tickCount, secondsToTest-1)
	assert.LessOrEqual(t, tickCount, secondsToTest+1)
}

func TestCappedRateLimitBackOff(t *testing.T) {
	const secondsToTest = 8
	rl := NewCappedRateLimiter(60, time.Minute, 10, .5)

	tickCount := 0
	keepLooping := true

	tick := func() {
		rl.LimitRate()
		tickCount++
		if (tickCount % 5) == 0 {
			rl.Backoff()
		} else {
			rl.Success()
		}
	}

	for i := 0; i < 5; i++ {
		go func() {
			for keepLooping {
				tick()
			}
		}()
	}

	time.Sleep(time.Second * secondsToTest)
	keepLooping = false

	assert.GreaterOrEqual(t, tickCount, secondsToTest/4)
	assert.LessOrEqual(t, tickCount, secondsToTest*3/4)
}

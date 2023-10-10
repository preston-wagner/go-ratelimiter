package ratelimiter

// import (
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// )

// type fakeQueryType string

// const (
// 	postType    fakeQueryType = "post"
// 	commentType fakeQueryType = "comment"
// 	replyType   fakeQueryType = "reply"
// )

// var queryTypes = []fakeQueryType{
// 	postType,
// 	commentType,
// 	replyType,
// }

// func TestRateLimitWithTypes(t *testing.T) {
// 	const secondsToTest = 8
// 	rlwt := NewRateLimiterWithTypes[fakeQueryType](60)

// 	tickCount := 0
// 	keepLooping := true

// 	looper := func(queryType fakeQueryType) {
// 		for keepLooping {
// 			rlwt.LimitRate(queryType)
// 			tickCount++
// 		}
// 	}

// 	for _, queryType := range queryTypes {
// 		go looper(queryType)
// 	}

// 	time.Sleep(time.Second * secondsToTest)
// 	keepLooping = false

// 	assert.GreaterOrEqual(t, tickCount, secondsToTest-1)
// 	assert.LessOrEqual(t, tickCount, secondsToTest+1)
// }

// func TestRateLimitWithTypesBackOff(t *testing.T) {
// 	const secondsToTest = 8
// 	rlwt := NewRateLimiterWithTypes[fakeQueryType](60)

// 	tickCount := 0
// 	keepLooping := true

// 	looper := func(queryType fakeQueryType) {
// 		for keepLooping {
// 			rlwt.LimitRate(queryType)
// 			tickCount++
// 			rlwt.Backoff(queryType)
// 		}
// 	}

// 	for _, queryType := range queryTypes {
// 		go looper(queryType)
// 	}

// 	time.Sleep(time.Second * secondsToTest)
// 	keepLooping = false

// 	assert.GreaterOrEqual(t, tickCount, secondsToTest*5/8)
// 	assert.LessOrEqual(t, tickCount, secondsToTest*7/8)
// }

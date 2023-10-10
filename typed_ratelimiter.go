package ratelimiter

// import (
// 	"log"
// 	"sync"
// )

// type RateLimiterWithTypes[REQUEST_TYPE comparable] struct {
// 	maxPerMinute int
// 	// cumulativeLimiter  *RateLimiter
// 	individualLimiters map[REQUEST_TYPE]*RateLimiter
// 	lock               *sync.RWMutex
// }

// func NewRateLimiterWithTypes[REQUEST_TYPE comparable](requestsPerMinute int) *RateLimiterWithTypes[REQUEST_TYPE] {
// 	return &RateLimiterWithTypes[REQUEST_TYPE]{
// 		maxPerMinute: requestsPerMinute,
// 		// cumulativeLimiter:  NewRateLimiter(requestsPerMinute),
// 		individualLimiters: map[REQUEST_TYPE]*RateLimiter{},
// 		lock:               &sync.RWMutex{},
// 	}
// }

// func (rlwt *RateLimiterWithTypes[REQUEST_TYPE]) createIndividualLimiter(key REQUEST_TYPE) {
// 	rlwt.lock.Lock()
// 	defer rlwt.lock.Unlock()
// 	_, ok := rlwt.individualLimiters[key]
// 	if !ok {
// 		rlwt.individualLimiters[key] = NewRateLimiter(rlwt.maxPerMinute)
// 	}
// }

// func (rlwt *RateLimiterWithTypes[REQUEST_TYPE]) getIndividualLimiter(key REQUEST_TYPE) (*RateLimiter, bool) {
// 	rlwt.lock.RLock()
// 	defer rlwt.lock.RUnlock()
// 	limiter, ok := rlwt.individualLimiters[key]
// 	return limiter, ok
// }

// func (rlwt *RateLimiterWithTypes[REQUEST_TYPE]) getOrCreateIndividualLimiter(key REQUEST_TYPE) *RateLimiter {
// 	limiter, ok := rlwt.getIndividualLimiter(key)
// 	if !ok {
// 		rlwt.createIndividualLimiter(key)
// 		limiter, _ = rlwt.getIndividualLimiter(key)
// 	}
// 	return limiter
// }

// func (rlwt *RateLimiterWithTypes[REQUEST_TYPE]) LimitRate(key REQUEST_TYPE) {
// 	// rlwt.cumulativeLimiter.LimitRate()
// 	rlwt.getOrCreateIndividualLimiter(key).LimitRate()
// }

// func (rlwt *RateLimiterWithTypes[REQUEST_TYPE]) Success(key REQUEST_TYPE) {
// 	rlwt.getOrCreateIndividualLimiter(key).Success()
// }

// func (rlwt *RateLimiterWithTypes[REQUEST_TYPE]) Backoff(key REQUEST_TYPE) {
// 	rl := rlwt.getOrCreateIndividualLimiter(key)
// 	rl.Backoff()
// 	log.Println("Backoff(", key, "); currentPerMinute:", rl.currentPerMinute)
// }

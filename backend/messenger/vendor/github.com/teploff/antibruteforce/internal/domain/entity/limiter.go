package entity

import (
	"sync"
	"time"
)

// Limiter was inspired by https://github.com/golang/go/wiki/RateLimiting.
// However, the go example is not good for setting high qps limits because
// it will cause the ticker to fire too often. Also, the ticker will continue
// to fire when the system is idle. This new Limiter achieves the same thing,
// but by using just counters with no tickers or channels.
type Limiter struct {
	maxCount int
	interval time.Duration

	mu       sync.Mutex
	curCount int
	lastTime time.Time
	lastSeen time.Time
}

// NewLimiter creates a new RateLimiter. maxCount is the max burst allowed
// while interval specifies the duration for a burst. The effective rate limit is
// equal to maxCount/interval. For example, if you want to a max QPS of 5000,
// and want to limit bursts to no more than 500, you'd specify a maxCount of 500
// and an interval of 100*time.Millilsecond.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	return &Limiter{
		maxCount: maxCount,
		interval: interval,
		lastSeen: time.Now(),
	}
}

// LastSeen has last event by current bucket.
func (l *Limiter) LastSeen() time.Time {
	return l.lastSeen
}

// Allow returns true if a request is within the rate limit norms.
// Otherwise, it returns false.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if time.Since(l.lastTime) < l.interval {
		if l.curCount > 0 {
			l.curCount--
			return true
		}

		return false
	}

	l.curCount = l.maxCount - 1
	l.lastTime = time.Now()
	l.lastSeen = time.Now()

	return true
}

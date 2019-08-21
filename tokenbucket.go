package tokenbucket

import (
	"sync"
	"time"
)

type TokenBucket struct {
	tokenCount, capacity     uint64
	refreshRate              time.Duration
	lastRefresh, nextRefresh time.Time
	sync.Mutex
}

// NewTokenBucket takes an initial capacity for the rate limiter and a refreshRate
// which determines the time between adding single units of capacity to the limiter
func New(initialCapacity, totalCapacity uint64, refreshRate time.Duration) *TokenBucket {
	t := time.Now()
	return &TokenBucket{
		tokenCount:  initialCapacity,
		capacity:    totalCapacity,
		refreshRate: refreshRate,
		lastRefresh: t,
		nextRefresh: t.Add(refreshRate),
	}
}

// Take will adjust the current number of tokens availble in the bucket
// and if less than 1 token is available will return false, otherwise Take
// will remove one token from the bucket and return true.
func (b *TokenBucket) Take() bool {
	b.Lock()
	defer b.Unlock()
	now := time.Now()
	if now.After(b.nextRefresh) {
		newTokens := uint64(time.Since(b.lastRefresh).Nanoseconds() / b.refreshRate.Nanoseconds())
		if newTokens+b.tokenCount <= b.capacity {
			b.tokenCount += newTokens
		} else {
			b.tokenCount = b.capacity
		}
		b.nextRefresh = now.Add(b.refreshRate)
		b.lastRefresh = now
	}

	if b.tokenCount < 1 {
		return false
	}
	b.tokenCount--
	return true
}

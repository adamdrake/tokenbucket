package tokenbucket

import (
	"sync"
	"testing"
	"time"
)

func TestTokenBucket_Take(t *testing.T) {
	type fields struct {
		tokenCount  uint64
		capacity    uint64
		refreshRate time.Duration
		lastRefresh time.Time
		nextRefresh time.Time
		Mutex       sync.Mutex
	}
	tests := []struct {
		name          string
		fields        fields
		want          bool
		expectedCount uint64
	}{
		{
			"tokens available now",
			fields{
				tokenCount:  10,
				capacity:    10,
				refreshRate: 1 * time.Second,
				lastRefresh: time.Now(),
				nextRefresh: time.Now(),
			},
			true,
			9,
		},
		{
			"no tokens available now",
			fields{
				tokenCount:  0,
				capacity:    10,
				refreshRate: 1 * time.Second,
				lastRefresh: time.Now(),
				nextRefresh: time.Now().Add(10 * time.Second),
			},
			false,
			0,
		},
		{
			"tokens available after adjustment",
			fields{
				tokenCount:  0,
				capacity:    10,
				refreshRate: 1 * time.Second,
				lastRefresh: time.Now().Add(-10 * time.Minute),
				nextRefresh: time.Now().Add(-1 * time.Second),
			},
			true,
			9,
		},
		{
			"tokens do not refresh above capacity",
			fields{
				tokenCount:  10,
				capacity:    10,
				refreshRate: 1 * time.Nanosecond,
				lastRefresh: time.Now().Add(-2 * time.Minute),
				nextRefresh: time.Now().Add(-1 * time.Second),
			},
			true,
			9,
		},
		{
			"refresh 5 tokens",
			fields{
				tokenCount:  0,
				capacity:    10,
				refreshRate: 1 * time.Second,
				lastRefresh: time.Now().Add(-5 * time.Second),
				nextRefresh: time.Now().Add(-5 * time.Second),
			},
			true,
			4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &TokenBucket{
				tokenCount:  tt.fields.tokenCount,
				capacity:    tt.fields.capacity,
				refreshRate: tt.fields.refreshRate,
				lastRefresh: tt.fields.lastRefresh,
				nextRefresh: tt.fields.nextRefresh,
				Mutex:       tt.fields.Mutex,
			}
			if got := b.Take(); got != tt.want {
				t.Errorf("TokenBucket.Take() = %v, want %v", got, tt.want)
			}
			if count := b.tokenCount; count != tt.expectedCount {
				t.Errorf("Token count incorrect.  Got %v, want %v", count, tt.expectedCount)
			}
			if b.tokenCount > b.capacity {
				t.Errorf("Capacity is %v but current count is %v", b.capacity, b.tokenCount)
			}
		})
	}
}

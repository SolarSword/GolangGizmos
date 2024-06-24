package ratelimiter

import (
	"sync"
	"time"
)

type BucketLimiter struct {
	mu   sync.Mutex
	size int
	// the current token in the bucket
	count int
	// the time interval to fill a new token in the bucket
	fillInterval time.Duration
	// the last timestamp to access the component
	lastAccessTime time.Time
}

func NewBucketLimiter(interval time.Duration, size int) *BucketLimiter {
	return &BucketLimiter{
		fillInterval: interval,
		size:         size,
		count:        size,
	}
}

func (l *BucketLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.allow()
}

func (l *BucketLimiter) getFillNumber() int {
	if l.count >= l.size {
		return 0
	}
	// no need to fill the token when the first access to the component
	if !l.lastAccessTime.IsZero() {
		d := time.Now().Sub(l.lastAccessTime)
		count := int(d / l.fillInterval)
		if l.size-l.count < count {
			return l.size - l.count
		}
		return count
	}
	return 0
}

func (l *BucketLimiter) fillToken() {
	l.count += l.getFillNumber()
}

func (l *BucketLimiter) allow() bool {
	// to fill the token in the bucket first
	l.fillToken()
	if l.count > 0 {
		l.count--
		l.lastAccessTime = time.Now()
		return true
	}
	return false
}

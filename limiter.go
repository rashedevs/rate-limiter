package main

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu     sync.Mutex
	bucket map[string]*Bucket
}

type Bucket struct {
	Count     int
	ExpiresAt time.Time
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		bucket: make(map[string]*Bucket),
	}
}

func (rl *RateLimiter) Check(key string, limit, window int) (bool, int, int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, exists := rl.bucket[key]

	if !exists || now.After(b.ExpiresAt) {
		rl.bucket[key] = &Bucket{
			Count:     1,
			ExpiresAt: now.Add(time.Duration(window) * time.Second),
		}
		return true, limit - 1, window
	}

	if b.Count < limit {
		b.Count++
		return true, limit - b.Count, int(b.ExpiresAt.Sub(now).Seconds())
	}

	return false, 0, int(b.ExpiresAt.Sub(now).Seconds())
}

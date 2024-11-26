package rl

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu            sync.Mutex
	clientData    map[string]*ClientData
	requestWindow time.Duration
	maxRequests   int
}

type ClientData struct {
	Ip             string
	Count          int
	LastAccessTime time.Time
}

func New(requestWindow time.Duration, maxRequests int) *RateLimiter {
	return &RateLimiter{
		clientData:    make(map[string]*ClientData),
		requestWindow: requestWindow,
		maxRequests:   maxRequests,
	}
}

func (r *RateLimiter) IsAllowed(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry, ok := r.clientData[ip]
	if !ok || time.Since(entry.LastAccessTime) > r.requestWindow {
		r.clientData[ip].LastAccessTime = time.Now()
		r.clientData[ip].Count = 1
		return true
	} else if entry.Count <= r.maxRequests {
		return true
	}

	return false
}

package rl

import (
	"sync"
	"time"
)

type RateLimiter struct {
	store         Store
	requestWindow time.Duration
	maxRequests   int
}

func New(store Store, requestWindow time.Duration, maxRequests int) *RateLimiter {
	return &RateLimiter{
		store:         store,
		requestWindow: requestWindow,
		maxRequests:   maxRequests,
	}
}

func (r *RateLimiter) IsAllowed(ip string) bool {
	entry, ok := r.store.GetClientData(ip)
	if !ok || entry.WindowExpiresAt.Before(time.Now()) {
		r.store.InitClientData(ip, r.requestWindow)
	}

	if entry.Count < r.maxRequests {
		r.store.Increment(ip)
		return true
	}

	return false
}

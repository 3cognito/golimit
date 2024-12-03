package main

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
	count, ok := r.store.ClientCount(ip)
	if !ok {
		r.store.InitClientData(ip, r.requestWindow)
		return true
	}

	if count < r.maxRequests {
		r.store.Increment(ip)
		return true
	}

	return false
}

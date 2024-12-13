package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRateLimiterWithInMemoryStore(t *testing.T) {
	store := NewInMemoryStore()
	rateLimiter := New(store, time.Second, 3)

	ip := "192.168.1.1"

	for i := 1; i <= 3; i++ {
		if !rateLimiter.IsAllowed(ip) {
			t.Errorf("InMemoryStore: Request %d: Expected request to be allowed, but it was denied", i)
		}
	}

	if rateLimiter.IsAllowed(ip) {
		t.Errorf("InMemoryStore: Expected fourth request to be denied, but it was allowed")
	}

	time.Sleep(time.Second)
	if !rateLimiter.IsAllowed(ip) {
		t.Errorf("InMemoryStore: Expected request to be allowed after window expiration, but it was denied")
	}
}

func TestRateLimiterWithRedisStore(t *testing.T) {
	store, err := getRedisStore()
	if err != nil {
		t.Skip("Skipping Redis tests: Redis not configured")
	}

	rateLimiter := New(store, time.Second, 3)
	ip := "192.168.1.1"

	for i := 1; i <= 3; i++ {
		if !rateLimiter.IsAllowed(ip) {
			t.Errorf("RedisStore: Request %d: Expected request to be allowed, but it was denied", i)
		}
	}

	if rateLimiter.IsAllowed(ip) {
		t.Errorf("RedisStore: Expected fourth request to be denied, but it was allowed")
	}

	time.Sleep(time.Second)
	if !rateLimiter.IsAllowed(ip) {
		t.Errorf("RedisStore: Expected request to be allowed after window expiration, but it was denied")
	}
}

func TestRateLimiterWithMultipleIPs(t *testing.T) {
	store := NewInMemoryStore()
	rateLimiter := New(store, time.Second, 2)

	ip1 := "192.168.1.1"
	ip2 := "192.168.1.2"

	for i := 1; i <= 2; i++ {
		if !rateLimiter.IsAllowed(ip1) {
			t.Errorf("InMemoryStore: IP1 Request %d: Expected request to be allowed, but it was denied", i)
		}
		if !rateLimiter.IsAllowed(ip2) {
			t.Errorf("InMemoryStore: IP2 Request %d: Expected request to be allowed, but it was denied", i)
		}
	}

	if rateLimiter.IsAllowed(ip1) {
		t.Errorf("InMemoryStore: Expected IP1 third request to be denied, but it was allowed")
	}

	if rateLimiter.IsAllowed(ip2) {
		t.Errorf("InMemoryStore: Expected IP2 third request to be denied, but it was allowed")
	}
}

func TestRateLimiter_WindowExpiration(t *testing.T) {
	store := NewInMemoryStore()
	rateLimiter := New(store, 2*time.Second, 2)

	ip := "192.168.1.1"

	for i := 1; i <= 2; i++ {
		if !rateLimiter.IsAllowed(ip) {
			t.Errorf("InMemoryStore: Request %d: Expected request to be allowed, but it was denied", i)
		}
	}

	if rateLimiter.IsAllowed(ip) {
		t.Errorf("InMemoryStore: Expected third request to be denied, but it was allowed")
	}

	time.Sleep(1 * time.Second)
	if rateLimiter.IsAllowed(ip) {
		t.Errorf("InMemoryStore: Expected request within same window to be denied, but it was allowed")
	}

	time.Sleep(1 * time.Second)
	if !rateLimiter.IsAllowed(ip) {
		t.Errorf("InMemoryStore: Expected request to be allowed after window expiration, but it was denied")
	}
}

func TestRateLimiterWithConcurrency(t *testing.T) {
	store := NewInMemoryStore()
	rateLimiter := New(store, time.Second, 5)

	ip := "192.168.1.1"
	var wg sync.WaitGroup
	requests := 20
	allowed := 0
	mu := sync.Mutex{}

	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if rateLimiter.IsAllowed(ip) {
				mu.Lock()
				allowed++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if allowed != 5 {
		t.Errorf("Concurrency Test: Expected 5 requests to be allowed, but got %d", allowed)
	}

	time.Sleep(time.Second)

	allowed = 0
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if rateLimiter.IsAllowed(ip) {
				mu.Lock()
				allowed++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	fmt.Printf("Concurrency Test After Window Expiry: Allowed %d requests\n", allowed)
	if allowed != 5 {
		t.Errorf("Concurrency Test After Window Expiry: Expected 5 requests to be allowed, but got %d", allowed)
	}
}

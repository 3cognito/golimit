package main

import (
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestInMemoryStore_InitClientData(t *testing.T) {
	store := NewInMemoryStore()
	ip := "192.168.1.1"
	duration := 2 * time.Second

	store.InitClientData(ip, duration)

	count, ok := store.ClientCount(ip)
	if !ok || count != 1 {
		t.Errorf("Expected count=1 and ok=true, got count=%d and ok=%v", count, ok)
	}

	time.Sleep(duration)
	count, ok = store.ClientCount(ip)
	if ok {
		t.Errorf("Expected data to expire, but got count=%d and ok=%v", count, ok)
	}
}

func TestInMemoryStore_Increment(t *testing.T) {
	store := NewInMemoryStore()
	ip := "192.168.1.1"
	duration := 2 * time.Second

	store.InitClientData(ip, duration)
	store.Increment(ip)
	store.Increment(ip)

	count, ok := store.ClientCount(ip)
	if !ok || count != 3 {
		t.Errorf("Expected count=3 and ok=true, got count=%d and ok=%v", count, ok)
	}
}

func getRedisStore() (*RedisStore, error) {
	redisURL := ""
	if redisURL == "" {
		return nil, redis.ErrClosed
	}
	return NewRedisStore(redisURL), nil
}

func TestRedisStore_InitClientData(t *testing.T) {
	store, err := getRedisStore()
	if err != nil {
		t.Skip("Skipping Redis tests: Redis not configured")
	}

	ip := "192.168.1.1"
	duration := 2 * time.Second

	store.InitClientData(ip, duration)

	count, ok := store.ClientCount(ip)
	if !ok || count != 1 {
		t.Errorf("Expected count=1 and ok=true, got count=%d and ok=%v", count, ok)
	}

	time.Sleep(duration)
	count, ok = store.ClientCount(ip)
	if ok {
		t.Errorf("Expected data to expire, but got count=%d and ok=%v", count, ok)
	}
}

func TestRedisStore_Increment(t *testing.T) {
	store, err := getRedisStore()
	if err != nil {
		t.Skip("Skipping Redis tests: Redis not configured")
	}

	ip := "192.168.1.1"
	duration := 2 * time.Second

	store.InitClientData(ip, duration)
	store.Increment(ip)
	store.Increment(ip)

	count, ok := store.ClientCount(ip)
	if !ok || count != 3 {
		t.Errorf("Expected count=3 and ok=true, got count=%d and ok=%v", count, ok)
	}
}

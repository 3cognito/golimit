package main

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type ClientData struct {
	Ip              string
	Count           int
	WindowExpiresAt time.Time
}

type Store interface {
	Increment(ip string)
	ClientCount(ip string) (int, bool)
	InitClientData(ip string, windowDuration time.Duration)
}

type InMemoryStore struct {
	data map[string]ClientData
	mu   sync.Mutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string]ClientData),
	}
}

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStore(addr string) *RedisStore {
	opt, _ := redis.ParseURL(addr)
	client := redis.NewClient(opt)

	return &RedisStore{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisStore) Increment(ip string) {
	r.client.Incr(r.ctx, ip)
}

func (r *RedisStore) InitClientData(ip string, windowDuration time.Duration) {
	r.client.Set(r.ctx, ip, 1, windowDuration)
}

func (r *RedisStore) ClientCount(ip string) (int, bool) {
	val, err := r.client.Get(r.ctx, ip).Result()
	if err == redis.Nil {
		return 0, false
	}

	count, _ := strconv.Atoi(val)
	return count, true
}

func (i *InMemoryStore) ClientCount(ip string) (int, bool) {
	i.mu.Lock()
	defer i.mu.Unlock()

	data, ok := i.data[ip]
	if !ok || data.WindowExpiresAt.Before(time.Now()) {
		return 0, false
	}

	return data.Count, true
}

func (i *InMemoryStore) Increment(ip string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	data, ok := i.data[ip]
	if !ok {
		return
	} else {
		data.Count++
	}

	i.data[ip] = data
}

func (i *InMemoryStore) InitClientData(ip string, windowDuration time.Duration) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.data[ip] = ClientData{
		Ip:              ip,
		Count:           1,
		WindowExpiresAt: time.Now().Add(windowDuration),
	}
}

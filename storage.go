package main

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// errors
var (
	ErrClientDataNotFound = errors.New("client data not found")
)

type ClientData struct {
	Ip              string
	Count           int
	WindowExpiresAt time.Time
}

type Store interface {
	Increment(ip string)
	GetClientData(ip string) (ClientData, bool)
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
	opt, _ := redis.ParseURL("rediss://default:AVDFAAIjcDFkZjllNzQxY2NlNmY0MjAzODJmOWUxNDZlNjI4YmQyOHAxMA@guiding-emu-20677.upstash.io:6379")
	client := redis.NewClient(opt)

	return &RedisStore{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisStore) Increment(ip string) {

}

func (r *RedisStore) InitClientData(ip string, count int, windowDuration time.Duration) {

}

func (r *RedisStore) GetClientData(ip string) (ClientData, bool) {
	return ClientData{}, false
}

func (i *InMemoryStore) GetClientData(ip string) (ClientData, bool) {
	i.mu.Lock()
	defer i.mu.Unlock()

	data, ok := i.data[ip]
	return data, ok
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

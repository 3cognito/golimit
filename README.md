# Rate Limiter

A simple rate limiter implementation in Go, which supports both in-memory and Redis-based stores.

## Features
- Limits the number of requests from a client (identified by IP) within a defined time window.
- Supports in-memory or Redis as the data store.

## Usage

### Create a Rate Limiter

```go
// For in-memory store
store := NewInMemoryStore()
rateLimiter := New(store, time.Minute, 100) // 100 requests per minute

// For Redis store
store := NewRedisStore("redis://localhost:6379")
rateLimiter := New(store, time.Minute, 100)
```

### Check if a request is allowed

```go
if rateLimiter.IsAllowed(clientIP) {
    // Proceed with the request
} else {
    // Reject the request (rate limit exceeded)
}
```

## Store Interfaces

The rate limiter uses the `Store` interface for storing request counts. You can implement your own store or use:
- `InMemoryStore` for local memory storage.
- `RedisStore` for Redis-backed storage.

## License

MIT License.
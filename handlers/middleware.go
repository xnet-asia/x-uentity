package handlers

import (
	"fmt"
	"sync"
	"time"
)

// Example middleware implementations

// ValidationMiddleware validates request data
type ValidationMiddleware[T any] struct{}

func (m *ValidationMiddleware[T]) Ingress(auth *ClientAuth, req *EntityRequest[T]) error {
	if req.Action == "" {
		return fmt.Errorf("action is required")
	}
	return nil
}

func (m *ValidationMiddleware[T]) Egress(auth *ClientAuth, resp *EntityResponse[T]) error {
	return nil
}

// LoggingMiddleware logs all requests/responses
type LoggingMiddleware[T any] struct {
	logger func(msg string)
}

func NewLoggingMiddleware[T any](logger func(msg string)) *LoggingMiddleware[T] {
	return &LoggingMiddleware[T]{logger: logger}
}

func (m *LoggingMiddleware[T]) Ingress(auth *ClientAuth, req *EntityRequest[T]) error {
	if m.logger != nil {
		m.logger(fmt.Sprintf("Request: client=%s action=%s", auth.ID, req.Action))
	}
	return nil
}

func (m *LoggingMiddleware[T]) Egress(auth *ClientAuth, resp *EntityResponse[T]) error {
	if m.logger != nil {
		m.logger(fmt.Sprintf("Response: client=%s code=%d success=%v", auth.ID, resp.Code, resp.Success))
	}
	return nil
}

// RateLimitMiddleware enforces rate limits
type RateLimitMiddleware[T any] struct {
	mu     sync.Mutex
	limits map[string]rateLimit
	maxReq int
	window time.Duration
}

type rateLimit struct {
	count   int
	resetAt time.Time
}

func NewRateLimitMiddleware[T any](maxReqPerWindow int, windows ...time.Duration) *RateLimitMiddleware[T] {
	window := time.Minute
	if len(windows) > 0 && windows[0] > 0 {
		window = windows[0]
	}
	return &RateLimitMiddleware[T]{
		limits: make(map[string]rateLimit),
		maxReq: maxReqPerWindow,
		window: window,
	}
}

func (m *RateLimitMiddleware[T]) Ingress(auth *ClientAuth, req *EntityRequest[T]) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	limit := m.limits[auth.ID]
	if limit.resetAt.IsZero() || !now.Before(limit.resetAt) {
		limit = rateLimit{resetAt: now.Add(m.window)}
	}
	if limit.count >= m.maxReq {
		return fmt.Errorf("rate limit exceeded")
	}
	limit.count++
	m.limits[auth.ID] = limit
	return nil
}

func (m *RateLimitMiddleware[T]) Egress(auth *ClientAuth, resp *EntityResponse[T]) error {
	return nil
}

// CachingMiddleware caches GET responses
type CachingMiddleware[T any] struct {
	mu    sync.RWMutex
	cache map[string]*EntityResponse[T]
}

func NewCachingMiddleware[T any]() *CachingMiddleware[T] {
	return &CachingMiddleware[T]{
		cache: make(map[string]*EntityResponse[T]),
	}
}

func (m *CachingMiddleware[T]) Ingress(auth *ClientAuth, req *EntityRequest[T]) error {
	return nil
}

func (m *CachingMiddleware[T]) Egress(auth *ClientAuth, resp *EntityResponse[T]) error {
	if resp.Success && resp.Single != nil {
		m.mu.Lock()
		copy := *resp
		m.cache[auth.ID] = &copy
		m.mu.Unlock()
	}
	return nil
}

// Get returns the last successful single-entity response for a client.
func (m *CachingMiddleware[T]) Get(clientID string) (*EntityResponse[T], bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	resp, ok := m.cache[clientID]
	if !ok {
		return nil, false
	}
	copy := *resp
	return &copy, true
}

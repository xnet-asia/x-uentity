package handlers

import (
	"fmt"
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
	m.logger(fmt.Sprintf("Request: %s action=%s from %s", auth.ID, req.Action, auth.Token))
	return nil
}

func (m *LoggingMiddleware[T]) Egress(auth *ClientAuth, resp *EntityResponse[T]) error {
	m.logger(fmt.Sprintf("Response: %s code=%d success=%v", auth.ID, resp.Code, resp.Success))
	return nil
}

// RateLimitMiddleware enforces rate limits
type RateLimitMiddleware[T any] struct {
	limits map[string]int // clientID -> request count
	maxReq int
}

func NewRateLimitMiddleware[T any](maxReqPerWindow int) *RateLimitMiddleware[T] {
	return &RateLimitMiddleware[T]{
		limits: make(map[string]int),
		maxReq: maxReqPerWindow,
	}
}

func (m *RateLimitMiddleware[T]) Ingress(auth *ClientAuth, req *EntityRequest[T]) error {
	if m.limits[auth.ID] >= m.maxReq {
		return fmt.Errorf("rate limit exceeded")
	}
	m.limits[auth.ID]++
	return nil
}

func (m *RateLimitMiddleware[T]) Egress(auth *ClientAuth, resp *EntityResponse[T]) error {
	return nil
}

// CachingMiddleware caches GET responses
type CachingMiddleware[T any] struct {
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
		key := fmt.Sprintf("get:%s", auth.ID)
		m.cache[key] = resp
	}
	return nil
}

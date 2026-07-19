package interfaces

import (
	"context"
	"github.com/xnetltd/x-uentity/repositories"
)

// ExternalService defines the interface for external service calls
type ExternalService[T any] interface {
	Create(ctx context.Context, entity T) error
	GetByID(ctx context.Context, id string) (T, error)
	Query(ctx context.Context, filter repositories.QueryFilter[T]) ([]T, error)
	Update(ctx context.Context, entity T) error
	Delete(ctx context.Context, id string) error
}

// ServiceResponse wraps responses from external services
type ServiceResponse[T any] struct {
	Success bool
	Data    T
	Error   string
	Code    int
}

// ServiceRequest wraps requests to external services
type ServiceRequest[T any] struct {
	Action  string
	Payload T
	Headers map[string]string
}

// MessageBroker defines interface for event-driven communication
type MessageBroker interface {
	Publish(topic string, message interface{}) error
	Subscribe(topic string, handler func(message interface{})) error
	Unsubscribe(topic string) error
}

// CacheLayer defines interface for caching
type CacheLayer[T any] interface {
	Set(key string, value T, ttl int) error
	Get(key string) (T, error)
	Delete(key string) error
	Flush() error
}

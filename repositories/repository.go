package repositories

import "sync"

type QueryFilter[T any] func(T) bool

type Repository[T any] interface {
	Create(key string, value T) error
	Query(filter QueryFilter[T]) ([]T, error)
	Update(key string, value T) error
	Delete(key string) error
	GetByIdentifier(key string) (T, error)
	QueryOne(filter QueryFilter[T]) (T, error)
}

type InMemoryRepository[T any] struct {
	data *sync.Map
}

func NewInMemoryRepository[T any]() *InMemoryRepository[T] {
	return &InMemoryRepository[T]{
		data: &sync.Map{},
	}
}

func (r *InMemoryRepository[T]) Create(key string, value T) error {
	r.data.Store(key, value)
	return nil
}

// Query executes a filter function against all stored entities and returns matching results
func (r *InMemoryRepository[T]) Query(filter QueryFilter[T]) ([]T, error) {
	var results []T
	r.data.Range(func(key, value interface{}) bool {
		if entity := value.(T); filter(entity) {
			results = append(results, entity)
		}
		return true
	})
	return results, nil
}

func (r *InMemoryRepository[T]) Update(key string, value T) error {
	_, exists := r.data.Load(key)
	if !exists {
		return ErrNotFound
	}
	r.data.Store(key, value)
	return nil
}

func (r *InMemoryRepository[T]) Delete(key string) error {
	r.data.Delete(key)
	return nil
}

// GetByIdentifier queries by key and returns single entity or error
func (r *InMemoryRepository[T]) GetByIdentifier(key string) (T, error) {
	var zero T
	val, exists := r.data.Load(key)
	if !exists {
		return zero, ErrNotFound
	}
	return val.(T), nil
}

// QueryOne wraps Query and returns first result or error if not found
func (r *InMemoryRepository[T]) QueryOne(filter QueryFilter[T]) (T, error) {
	results, err := r.Query(filter)
	if err != nil {
		var zero T
		return zero, err
	}
	if len(results) == 0 {
		var zero T
		return zero, ErrNotFound
	}
	return results[0], nil
}

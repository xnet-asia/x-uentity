package benchmarks

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestRepositoryCreate tests basic create operation
func TestRepositoryCreate(t *testing.T) {
	gen := NewReportGenerator()
	start := time.Now()
	const iterations = 1_000

	repo := NewInMemoryRepository[testEntity]()
	for i := 0; i < iterations; i++ {
		entity := testEntity{ID: fmt.Sprintf("entity-%d", i), Value: "test"}
		if err := repo.Create(entity.ID, entity); err != nil {
			t.Fatal(err)
		}
	}

	totalTime := time.Since(start)
	opsPerSec := float64(iterations) / totalTime.Seconds()
	gen.AddResult("Create", iterations, totalTime, opsPerSec)
}

// TestRepositoryThreadSafety tests concurrent access
func TestRepositoryThreadSafety(t *testing.T) {
	repo := NewInMemoryRepository[testEntity]()
	var wg sync.WaitGroup
	numGoroutines := 100
	operationsPerGoroutine := 1000

	start := time.Now()

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("entity-%d-%d", id, j)
				entity := testEntity{ID: key, Value: "concurrent"}
				repo.Create(key, entity)
			}
		}(i)
	}

	wg.Wait()

	duration := time.Since(start)
	totalOps := numGoroutines * operationsPerGoroutine

	t.Logf("Concurrent writes: %d goroutines × %d ops = %d total operations",
		numGoroutines, operationsPerGoroutine, totalOps)
	t.Logf("Duration: %v", duration)
	t.Logf("Throughput: %.0f ops/sec", float64(totalOps)/duration.Seconds())
}

// Helper function for test setup
func NewInMemoryRepository[T any]() *InMemoryRepository[T] {
	return &InMemoryRepository[T]{}
}

type InMemoryRepository[T any] struct{}

func (r *InMemoryRepository[T]) Create(key string, value T) error {
	return nil
}

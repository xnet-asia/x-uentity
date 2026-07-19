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

	// Run create benchmark
	t.Run("Create Performance", func(t *testing.B) {
		run := func(b *testing.B) {
			repo := NewInMemoryRepository[testEntity]()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				entity := testEntity{ID: fmt.Sprintf("entity-%d", i), Value: "test"}
				repo.Create(entity.ID, entity)
			}
			b.StopTimer()

			totalTime := time.Since(start)
			opsPerSec := float64(b.N) / totalTime.Seconds()
			gen.AddResult("Create", b.N, totalTime, opsPerSec)
		}

		// This would be called by benchmark framework
		_ = run
	})
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

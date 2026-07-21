package benchmarks

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/xnetltd/x-uentity/repositories"
)

// BenchmarkResult holds the results of a performance test
type BenchmarkResult struct {
	Operation      string
	Iterations     int
	TotalTime      time.Duration
	AvgTime        time.Duration
	OpsPerSecond   float64
	Allocations    int
	BytesAllocated int
}

// PerformanceReport holds multiple benchmark results
type PerformanceReport struct {
	Timestamp  time.Time
	Results    []BenchmarkResult
	TotalTime  time.Duration
	MemBefore  uint64
	MemAfter   uint64
	MemAlloced uint64
}

// BenchmarkRepositoryCreate benchmarks the Create operation
func BenchmarkRepositoryCreate(b *testing.B) {
	repo := repositories.NewInMemoryRepository[testEntity]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entity := testEntity{ID: fmt.Sprintf("entity-%d", i), Value: "test"}
		repo.Create(entity.ID, entity)
	}
}

// BenchmarkRepositoryQuery benchmarks the Query operation
func BenchmarkRepositoryQuery(b *testing.B) {
	repo := repositories.NewInMemoryRepository[testEntity]()

	// Pre-populate with 1000 entities
	for i := 0; i < 1000; i++ {
		entity := testEntity{ID: fmt.Sprintf("entity-%d", i), Value: "test"}
		repo.Create(entity.ID, entity)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.Query(func(e testEntity) bool {
			return e.Value == "test"
		})
	}
}

// BenchmarkRepositoryGetByIdentifier benchmarks the GetByIdentifier operation
func BenchmarkRepositoryGetByIdentifier(b *testing.B) {
	repo := repositories.NewInMemoryRepository[testEntity]()

	// Pre-populate
	for i := 0; i < 1000; i++ {
		entity := testEntity{ID: fmt.Sprintf("entity-%d", i), Value: "test"}
		repo.Create(entity.ID, entity)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.GetByIdentifier(fmt.Sprintf("entity-%d", i%1000))
	}
}

// BenchmarkRepositoryUpdate benchmarks the Update operation
func BenchmarkRepositoryUpdate(b *testing.B) {
	repo := repositories.NewInMemoryRepository[testEntity]()

	// Pre-populate
	for i := 0; i < 1000; i++ {
		entity := testEntity{ID: fmt.Sprintf("entity-%d", i), Value: "test"}
		repo.Create(entity.ID, entity)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		updatedEntity := testEntity{ID: fmt.Sprintf("entity-%d", i%1000), Value: "updated"}
		repo.Update(updatedEntity.ID, updatedEntity)
	}
}

// BenchmarkRepositoryDelete benchmarks the Delete operation
func BenchmarkRepositoryDelete(b *testing.B) {
	repo := repositories.NewInMemoryRepository[testEntity]()

	// Create entities for each iteration
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entity := testEntity{ID: fmt.Sprintf("entity-%d", i), Value: "test"}
		repo.Create(entity.ID, entity)
	}

	b.StopTimer()
	// Now delete them
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		repo.Delete(fmt.Sprintf("entity-%d", i))
	}
}

// BenchmarkConcurrentOperations benchmarks concurrent read/write operations
func BenchmarkConcurrentOperations(b *testing.B) {
	repo := repositories.NewInMemoryRepository[testEntity]()

	// Pre-populate
	for i := 0; i < 100; i++ {
		entity := testEntity{ID: fmt.Sprintf("entity-%d", i), Value: "test"}
		repo.Create(entity.ID, entity)
	}

	b.ResetTimer()
	var wg sync.WaitGroup
	numGoroutines := 10
	operationsPerGoroutine := b.N / numGoroutines

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for i := 0; i < operationsPerGoroutine; i++ {
				if i%3 == 0 {
					// Create
					entity := testEntity{ID: fmt.Sprintf("entity-%d-%d", goroutineID, i), Value: "concurrent"}
					repo.Create(entity.ID, entity)
				} else if i%3 == 1 {
					// Read
					repo.GetByIdentifier(fmt.Sprintf("entity-%d", i%100))
				} else {
					// Query
					repo.Query(func(e testEntity) bool {
						return e.Value == "test"
					})
				}
			}
		}(g)
	}

	wg.Wait()
}

// testEntity is a simple test entity
type testEntity struct {
	ID    string
	Value string
}

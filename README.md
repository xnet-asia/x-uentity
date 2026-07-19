# x-uentity

A universal entity management library for Go providing clean architecture layers with thread-safe storage and query-based filtering.

## Features

- 🏗️ **Clean Architecture Layers**: Domain, Repository, Usecase, and Handler layers
- 🔒 **Thread-Safe**: Built on `sync.Map` for concurrent access
- 🔍 **Query Pattern**: Flexible filtering with `QueryFilter[T]` functions
- 📦 **Generic Types**: Reusable for any entity type
- 🎯 **External Interface Layer**: Ready to receive calls from external services
- 🏭 **Factory Pattern**: Easy instantiation of components

## Project Structure

```
x-uentity/
├── domain/              # Domain models
├── repositories/        # Repository layer
├── usecases/           # Business logic layer
├── handlers/           # HTTP handlers
├── interfaces/         # External service interfaces
├── factories/          # Factory patterns
├── examples/           # Usage examples
└── go.mod
```

## Quick Start

```go
import "github.com/xnetltd/x-uentity/repositories"

// Create repository
repo := repositories.NewInMemoryRepository[MyEntity]()

// Query
results, _ := repo.Query(func(e MyEntity) bool {
    return e.Status == "active"
})
```

## License

MIT

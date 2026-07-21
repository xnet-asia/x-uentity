# x-uentity

A minimal, clean entity management library for Go with P2P-first architecture supporting injectable middleware.

## Features

- 🔌 **P2P-First**: Primary P2P support with HTTP fallback
- 🧬 **Generic Types**: Reusable for any entity type
- 🔒 **Thread-Safe**: Built on `sync.Map` for concurrent access
- 🎯 **Middleware Injection**: Extensible ingress/egress pipeline
- 🏭 **Factory Pattern**: Easy instantiation with middleware composition
- 🔓 **Auth & Anonymous**: Support for authenticated and anonymous clients

## Quick Start

```go
package main

import (
	"context"
	"github.com/xnetltd/x-uentity/factories"
	"github.com/xnetltd/x-uentity/handlers"
)

type User struct {
	ID   string
	Name string
}

func main() {
	// Create factory with middleware
	factory := factories.NewFactory[User]().
		WithMiddleware(
			&handlers.ValidationMiddleware[User]{},
			handlers.NewLoggingMiddleware[User](println),
		)

	// Register clients (auth and anon)
	factory.GetAuthHandler().Register("client-1", "token123")

	// Get P2P server
	server := factory.GetP2PServer()
	
	// Handle requests
	resp, _ := server.HandleHTTPFallback("token123", &handlers.EntityRequest[User]{
		Action: "query",
		Filter: func(u User) bool { return u.Name == "John" },
		Context: context.Background(),
	})
}
```

## Architecture

```
x-uentity/
├── domain/              # Domain models
├── repositories/        # Repository layer
├── usecases/           # Business logic
├── handlers/           # P2P handlers + middleware
│   ├── entity_handler.go     # Entity CRUD handler
│   ├── auth_handler.go       # Auth/token validation
│   ├── p2p_server.go         # P2P with HTTP fallback
│   └── middleware_examples.go # Example middleware
├── interfaces/         # Service interfaces
├── factories/          # Factory with middleware injection
└── go.mod
```

## Middleware

Implement `Middleware[T]` interface for ingress/egress processing:

```go
type Middleware[T any] interface {
	Ingress(auth *ClientAuth, req *EntityRequest[T]) error
	Egress(auth *ClientAuth, resp *EntityResponse[T]) error
}
```

### Built-in Examples

- `ValidationMiddleware` - Request validation
- `LoggingMiddleware` - Request/response logging
- `RateLimitMiddleware` - Rate limiting per client
- `CachingMiddleware` - Response caching

## License

MIT

package factories

import (
	"github.com/xnetltd/x-uentity/handlers"
	"github.com/xnetltd/x-uentity/repositories"
	"github.com/xnetltd/x-uentity/usecases"
)

// Factory creates instances with injectable middleware
type Factory[T any] struct {
	repo            repositories.Repository[T]
	usecase         usecases.Usecase[T]
	entityHandler   *handlers.EntityHandler[T]
	authHandler     handlers.AuthHandler
	p2pServer       *handlers.P2PServer[T]
	middlewareChain *handlers.MiddlewareChain[T]
}

func NewFactory[T any]() *Factory[T] {
	repo := repositories.NewInMemoryRepository[T]()
	usecase := usecases.NewBaseUsecase(repo)
	authHandler := handlers.NewSimpleAuthHandler()

	// Default empty middleware chain
	middlewareChain := handlers.NewMiddlewareChain[T]()
	entityHandler := handlers.NewEntityHandler(repo, middlewareChain)
	p2pServer := handlers.NewP2PServer(entityHandler, authHandler)

	return &Factory[T]{
		repo:            repo,
		usecase:         usecase,
		entityHandler:   entityHandler,
		authHandler:     authHandler,
		p2pServer:       p2pServer,
		middlewareChain: middlewareChain,
	}
}

// WithMiddleware injects middleware that participates in both phases.
func (f *Factory[T]) WithMiddleware(middlewares ...handlers.Middleware[T]) *Factory[T] {
	f.middlewareChain = handlers.NewMiddlewareChain(middlewares...)
	f.rebuildHandlers()
	return f
}

// WithPipeline injects ingress and egress middleware independently.
func (f *Factory[T]) WithPipeline(
	ingress []handlers.IngressMiddleware[T],
	egress []handlers.EgressMiddleware[T],
) *Factory[T] {
	f.middlewareChain = handlers.NewMiddlewarePipeline(ingress, egress)
	f.rebuildHandlers()
	return f
}

func (f *Factory[T]) rebuildHandlers() {
	f.entityHandler = handlers.NewEntityHandler(f.repo, f.middlewareChain)
	f.p2pServer = handlers.NewP2PServer(f.entityHandler, f.authHandler)
}

func (f *Factory[T]) GetRepository() repositories.Repository[T] {
	return f.repo
}

func (f *Factory[T]) GetUsecase() usecases.Usecase[T] {
	return f.usecase
}

func (f *Factory[T]) GetEntityHandler() *handlers.EntityHandler[T] {
	return f.entityHandler
}

func (f *Factory[T]) GetAuthHandler() handlers.AuthHandler {
	return f.authHandler
}

func (f *Factory[T]) GetP2PServer() *handlers.P2PServer[T] {
	return f.p2pServer
}

package factories

import (
	"github.com/xnetltd/x-uentity/repositories"
	"github.com/xnetltd/x-uentity/usecases"
	"github.com/xnetltd/x-uentity/handlers"
	"github.com/gin-gonic/gin"
)

// Factory creates instances of repositories, usecases, and handlers
type Factory[T any] struct {
	repo    repositories.Repository[T]
	usecase usecases.Usecase[T]
	handler handlers.Handler
}

func NewFactory[T any]() *Factory[T] {
	repo := repositories.NewInMemoryRepository[T]()
	usecase := usecases.NewBaseUsecase(repo)

	return &Factory[T]{
		repo:    repo,
		usecase: usecase,
	}
}

func (f *Factory[T]) GetRepository() repositories.Repository[T] {
	return f.repo
}

func (f *Factory[T]) GetUsecase() usecases.Usecase[T] {
	return f.usecase
}

func (f *Factory[T]) SetHandler(handler handlers.Handler) {
	f.handler = handler
}

func (f *Factory[T]) GetHandler() handlers.Handler {
	return f.handler
}

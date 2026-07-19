package usecases

import "github.com/xnetltd/x-uentity/repositories"

// Usecase defines the interface for business logic layer
type Usecase[T any] interface {
	Create(entity T) error
	GetByID(id string) (T, error)
	GetAll(filter repositories.QueryFilter[T]) ([]T, error)
	Update(entity T) error
	Delete(id string) error
}

// BaseUsecase provides common usecase implementations
type BaseUsecase[T any] struct {
	repo repositories.Repository[T]
}

func NewBaseUsecase[T any](repo repositories.Repository[T]) *BaseUsecase[T] {
	return &BaseUsecase[T]{repo: repo}
}

func (u *BaseUsecase[T]) Create(entity T) error {
	// To be overridden in specific implementations
	return nil
}

func (u *BaseUsecase[T]) GetByID(id string) (T, error) {
	return u.repo.GetByIdentifier(id)
}

func (u *BaseUsecase[T]) GetAll(filter repositories.QueryFilter[T]) ([]T, error) {
	return u.repo.Query(filter)
}

func (u *BaseUsecase[T]) Update(entity T) error {
	// To be overridden in specific implementations
	return nil
}

func (u *BaseUsecase[T]) Delete(id string) error {
	return u.repo.Delete(id)
}

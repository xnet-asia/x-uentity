package handlers

import (
	"context"
	"github.com/xnetltd/x-uentity/repositories"
)

// ClientAuth represents client authentication state
type ClientAuth struct {
	ID     string
	Token  string
	IsAuth bool
}

// EntityRequest is a P2P/HTTP request
type EntityRequest[T any] struct {
	Action  string
	ID      string
	Data    T
	Filter  repositories.QueryFilter[T]
	Context context.Context
}

// EntityResponse wraps the response
type EntityResponse[T any] struct {
	Success bool
	Data    []T
	Single  *T
	Error   string
	Code    int
}

// Middleware allows ingress/egress processing
type Middleware[T any] interface {
	// Ingress processes request before handler
	Ingress(auth *ClientAuth, req *EntityRequest[T]) error
	
	// Egress processes response after handler
	Egress(auth *ClientAuth, resp *EntityResponse[T]) error
}

// MiddlewareChain chains multiple middleware
type MiddlewareChain[T any] struct {
	middlewares []Middleware[T]
}

func NewMiddlewareChain[T any](middlewares ...Middleware[T]) *MiddlewareChain[T] {
	return &MiddlewareChain[T]{middlewares: middlewares}
}

func (mc *MiddlewareChain[T]) Ingress(auth *ClientAuth, req *EntityRequest[T]) error {
	for _, mw := range mc.middlewares {
		if err := mw.Ingress(auth, req); err != nil {
			return err
		}
	}
	return nil
}

func (mc *MiddlewareChain[T]) Egress(auth *ClientAuth, resp *EntityResponse[T]) error {
	// Run in reverse order
	for i := len(mc.middlewares) - 1; i >= 0; i-- {
		if err := mc.middlewares[i].Egress(auth, resp); err != nil {
			return err
		}
	}
	return nil
}

// EntityHandler handles entity operations with middleware support
type EntityHandler[T any] struct {
	repo       repositories.Repository[T]
	middleware *MiddlewareChain[T]
}

func NewEntityHandler[T any](repo repositories.Repository[T], middleware *MiddlewareChain[T]) *EntityHandler[T] {
	return &EntityHandler[T]{
		repo:       repo,
		middleware: middleware,
	}
}

func (h *EntityHandler[T]) Handle(auth *ClientAuth, req *EntityRequest[T]) (*EntityResponse[T], error) {
	// Ingress middleware
	if h.middleware != nil {
		if err := h.middleware.Ingress(auth, req); err != nil {
			return &EntityResponse[T]{Success: false, Error: err.Error(), Code: 400}, nil
		}
	}

	var resp *EntityResponse[T]

	switch req.Action {
	case "create":
		resp = h.handleCreate(auth, req)
	case "get":
		resp = h.handleGet(auth, req)
	case "query":
		resp = h.handleQuery(auth, req)
	case "update":
		resp = h.handleUpdate(auth, req)
	case "delete":
		resp = h.handleDelete(auth, req)
	default:
		resp = &EntityResponse[T]{Success: false, Error: "unknown action", Code: 400}
	}

	// Egress middleware
	if h.middleware != nil {
		if err := h.middleware.Egress(auth, resp); err != nil {
			return &EntityResponse[T]{Success: false, Error: err.Error(), Code: 500}, nil
		}
	}

	return resp, nil
}

func (h *EntityHandler[T]) handleCreate(auth *ClientAuth, req *EntityRequest[T]) *EntityResponse[T] {
	if !auth.IsAuth {
		return &EntityResponse[T]{Success: false, Error: "auth required", Code: 403}
	}
	if err := h.repo.Create(req.Context, req.Data); err != nil {
		return &EntityResponse[T]{Success: false, Error: err.Error(), Code: 400}
	}
	return &EntityResponse[T]{Success: true, Single: &req.Data, Code: 201}
}

func (h *EntityHandler[T]) handleGet(auth *ClientAuth, req *EntityRequest[T]) *EntityResponse[T] {
	entity, err := h.repo.GetByID(req.Context, req.ID)
	if err != nil {
		return &EntityResponse[T]{Success: false, Error: "not found", Code: 404}
	}
	return &EntityResponse[T]{Success: true, Single: &entity, Code: 200}
}

func (h *EntityHandler[T]) handleQuery(auth *ClientAuth, req *EntityRequest[T]) *EntityResponse[T] {
	results, err := h.repo.Query(req.Context, req.Filter)
	if err != nil {
		return &EntityResponse[T]{Success: false, Error: err.Error(), Code: 400}
	}
	return &EntityResponse[T]{Success: true, Data: results, Code: 200}
}

func (h *EntityHandler[T]) handleUpdate(auth *ClientAuth, req *EntityRequest[T]) *EntityResponse[T] {
	if !auth.IsAuth {
		return &EntityResponse[T]{Success: false, Error: "auth required", Code: 403}
	}
	if err := h.repo.Update(req.Context, req.Data); err != nil {
		return &EntityResponse[T]{Success: false, Error: err.Error(), Code: 400}
	}
	return &EntityResponse[T]{Success: true, Single: &req.Data, Code: 200}
}

func (h *EntityHandler[T]) handleDelete(auth *ClientAuth, req *EntityRequest[T]) *EntityResponse[T] {
	if !auth.IsAuth {
		return &EntityResponse[T]{Success: false, Error: "auth required", Code: 403}
	}
	if err := h.repo.Delete(req.Context, req.ID); err != nil {
		return &EntityResponse[T]{Success: false, Error: err.Error(), Code: 400}
	}
	return &EntityResponse[T]{Success: true, Code: 200}
}

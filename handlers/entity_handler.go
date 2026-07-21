package handlers

import (
	"context"
	"errors"

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
	Action  string                      `json:"action"`
	ID      string                      `json:"id,omitempty"`
	Data    T                           `json:"data,omitempty"`
	Filter  repositories.QueryFilter[T] `json:"-"`
	Context context.Context             `json:"-"`
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

// IngressMiddleware processes a request before it reaches the entity handler.
type IngressMiddleware[T any] interface {
	Ingress(auth *ClientAuth, req *EntityRequest[T]) error
}

// EgressMiddleware processes a response after it leaves the entity handler.
type EgressMiddleware[T any] interface {
	Egress(auth *ClientAuth, resp *EntityResponse[T]) error
}

// MiddlewareChain chains multiple middleware
type MiddlewareChain[T any] struct {
	ingress []IngressMiddleware[T]
	egress  []EgressMiddleware[T]
}

func NewMiddlewareChain[T any](middlewares ...Middleware[T]) *MiddlewareChain[T] {
	ingress := make([]IngressMiddleware[T], len(middlewares))
	egress := make([]EgressMiddleware[T], len(middlewares))
	for i, middleware := range middlewares {
		ingress[i] = middleware
		egress[i] = middleware
	}
	return NewMiddlewarePipeline(ingress, egress)
}

// NewMiddlewarePipeline allows ingress and egress middleware to be injected
// independently. Egress middleware executes in reverse registration order.
func NewMiddlewarePipeline[T any](ingress []IngressMiddleware[T], egress []EgressMiddleware[T]) *MiddlewareChain[T] {
	return &MiddlewareChain[T]{
		ingress: append([]IngressMiddleware[T](nil), ingress...),
		egress:  append([]EgressMiddleware[T](nil), egress...),
	}
}

func (mc *MiddlewareChain[T]) Ingress(auth *ClientAuth, req *EntityRequest[T]) error {
	for _, mw := range mc.ingress {
		if err := mw.Ingress(auth, req); err != nil {
			return err
		}
	}
	return nil
}

func (mc *MiddlewareChain[T]) Egress(auth *ClientAuth, resp *EntityResponse[T]) error {
	// Run in reverse order
	for i := len(mc.egress) - 1; i >= 0; i-- {
		if err := mc.egress[i].Egress(auth, resp); err != nil {
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
	if req == nil {
		return &EntityResponse[T]{Success: false, Error: "request is required", Code: 400}, nil
	}
	if auth == nil {
		auth = &ClientAuth{ID: "anonymous"}
	}
	if req.Context == nil {
		req.Context = context.Background()
	}
	if err := req.Context.Err(); err != nil {
		return nil, err
	}

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
	if auth == nil || !auth.IsAuth {
		return &EntityResponse[T]{Success: false, Error: "auth required", Code: 403}
	}
	if req.ID == "" {
		return &EntityResponse[T]{Success: false, Error: "id is required", Code: 400}
	}
	if err := h.repo.Create(req.ID, req.Data); err != nil {
		return &EntityResponse[T]{Success: false, Error: err.Error(), Code: 400}
	}
	return &EntityResponse[T]{Success: true, Single: &req.Data, Code: 201}
}

func (h *EntityHandler[T]) handleGet(auth *ClientAuth, req *EntityRequest[T]) *EntityResponse[T] {
	entity, err := h.repo.GetByIdentifier(req.ID)
	if err != nil {
		return repositoryErrorResponse[T](err)
	}
	return &EntityResponse[T]{Success: true, Single: &entity, Code: 200}
}

func (h *EntityHandler[T]) handleQuery(auth *ClientAuth, req *EntityRequest[T]) *EntityResponse[T] {
	filter := req.Filter
	if filter == nil {
		filter = func(T) bool { return true }
	}
	results, err := h.repo.Query(filter)
	if err != nil {
		return &EntityResponse[T]{Success: false, Error: err.Error(), Code: 400}
	}
	return &EntityResponse[T]{Success: true, Data: results, Code: 200}
}

func (h *EntityHandler[T]) handleUpdate(auth *ClientAuth, req *EntityRequest[T]) *EntityResponse[T] {
	if auth == nil || !auth.IsAuth {
		return &EntityResponse[T]{Success: false, Error: "auth required", Code: 403}
	}
	if req.ID == "" {
		return &EntityResponse[T]{Success: false, Error: "id is required", Code: 400}
	}
	if err := h.repo.Update(req.ID, req.Data); err != nil {
		return repositoryErrorResponse[T](err)
	}
	return &EntityResponse[T]{Success: true, Single: &req.Data, Code: 200}
}

func (h *EntityHandler[T]) handleDelete(auth *ClientAuth, req *EntityRequest[T]) *EntityResponse[T] {
	if auth == nil || !auth.IsAuth {
		return &EntityResponse[T]{Success: false, Error: "auth required", Code: 403}
	}
	if req.ID == "" {
		return &EntityResponse[T]{Success: false, Error: "id is required", Code: 400}
	}
	if err := h.repo.Delete(req.ID); err != nil {
		return &EntityResponse[T]{Success: false, Error: err.Error(), Code: 400}
	}
	return &EntityResponse[T]{Success: true, Code: 200}
}

func repositoryErrorResponse[T any](err error) *EntityResponse[T] {
	if errors.Is(err, repositories.ErrNotFound) {
		return &EntityResponse[T]{Success: false, Error: err.Error(), Code: 404}
	}
	return &EntityResponse[T]{Success: false, Error: err.Error(), Code: 400}
}

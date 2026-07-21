package handlers

import (
	"errors"
	"sync"
)

var ErrInvalidCredentials = errors.New("client id and token are required")

// AuthHandler handles client authentication
type AuthHandler interface {
	// Authenticate validates token and returns client auth state
	Authenticate(token string) (*ClientAuth, error)

	// Register creates new client auth
	Register(id, token string) (*ClientAuth, error)

	// Revoke removes client auth
	Revoke(id string) error
}

// SimpleAuthHandler implements AuthHandler
type SimpleAuthHandler struct {
	mu     sync.RWMutex
	tokens map[string]*ClientAuth // token -> ClientAuth
}

func NewSimpleAuthHandler() *SimpleAuthHandler {
	return &SimpleAuthHandler{
		tokens: make(map[string]*ClientAuth),
	}
}

func (h *SimpleAuthHandler) Authenticate(token string) (*ClientAuth, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if auth, exists := h.tokens[token]; exists {
		copy := *auth
		return &copy, nil
	}
	return newAnonymousAuth(), nil
}

func (h *SimpleAuthHandler) Register(id, token string) (*ClientAuth, error) {
	if id == "" || token == "" {
		return nil, ErrInvalidCredentials
	}

	auth := &ClientAuth{
		ID:     id,
		Token:  token,
		IsAuth: true,
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.tokens[token] = auth
	copy := *auth
	return &copy, nil
}

func (h *SimpleAuthHandler) Revoke(id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	for token, auth := range h.tokens {
		if auth.ID == id {
			delete(h.tokens, token)
		}
	}
	return nil
}

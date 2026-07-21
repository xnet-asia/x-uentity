package handlers

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
	tokens map[string]*ClientAuth // token -> ClientAuth
}

func NewSimpleAuthHandler() *SimpleAuthHandler {
	return &SimpleAuthHandler{
		tokens: make(map[string]*ClientAuth),
	}
}

func (h *SimpleAuthHandler) Authenticate(token string) (*ClientAuth, error) {
	if auth, exists := h.tokens[token]; exists {
		return auth, nil
	}
	return &ClientAuth{IsAuth: false}, nil // Anonymous
}

func (h *SimpleAuthHandler) Register(id, token string) (*ClientAuth, error) {
	auth := &ClientAuth{
		ID:     id,
		Token:  token,
		IsAuth: true,
	}
	h.tokens[token] = auth
	return auth, nil
}

func (h *SimpleAuthHandler) Revoke(id string) error {
	for token, auth := range h.tokens {
		if auth.ID == id {
			delete(h.tokens, token)
		}
	}
	return nil
}

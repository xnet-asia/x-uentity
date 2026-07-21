package handlers

import (
	"encoding/json"
	"io"
	"net"
)

// P2PServer handles P2P connections with HTTP fallback
type P2PServer[T any] struct {
	entityHandler *EntityHandler[T]
	authHandler   AuthHandler
}

func NewP2PServer[T any](entityHandler *EntityHandler[T], authHandler AuthHandler) *P2PServer[T] {
	return &P2PServer[T]{
		entityHandler: entityHandler,
		authHandler:   authHandler,
	}
}

// HandleP2PConnection handles a P2P client connection
func (s *P2PServer[T]) HandleP2PConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var msg struct {
			Token   string
			Request *EntityRequest[T]
		}

		if err := decoder.Decode(&msg); err != nil {
			if err != io.EOF {
				break
			}
		}

		// Authenticate
		auth, _ := s.authHandler.Authenticate(msg.Token)

		// Handle entity request
		resp, _ := s.entityHandler.Handle(auth, msg.Request)

		// Send response back
		encoder.Encode(resp)
	}
}

// HandleHTTPFallback handles HTTP request as fallback
func (s *P2PServer[T]) HandleHTTPFallback(token string, req *EntityRequest[T]) (*EntityResponse[T], error) {
	// Authenticate
	auth, _ := s.authHandler.Authenticate(token)

	// Handle entity request
	return s.entityHandler.Handle(auth, req)
}

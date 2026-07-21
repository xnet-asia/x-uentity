package handlers

import (
	"encoding/json"
	"errors"
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
	_ = s.ServeConnection(conn)
}

// ServeConnection handles newline-delimited JSON requests until the peer
// disconnects. It returns transport and authentication errors to callers that
// need lifecycle control.
func (s *P2PServer[T]) ServeConnection(conn net.Conn) error {
	if conn == nil {
		return errors.New("connection is required")
	}
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var msg struct {
			Token   string            `json:"token"`
			Request *EntityRequest[T] `json:"request"`
		}

		if err := decoder.Decode(&msg); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		// Authenticate
		auth, err := s.authHandler.Authenticate(msg.Token)
		if err != nil {
			return err
		}

		// Handle entity request
		resp, err := s.entityHandler.Handle(auth, msg.Request)
		if err != nil {
			return err
		}

		// Send response back
		if err := encoder.Encode(resp); err != nil {
			return err
		}
	}
}

// HandleHTTPFallback handles HTTP request as fallback
func (s *P2PServer[T]) HandleHTTPFallback(token string, req *EntityRequest[T]) (*EntityResponse[T], error) {
	// Authenticate
	auth, err := s.authHandler.Authenticate(token)
	if err != nil {
		return nil, err
	}

	// Handle entity request
	return s.entityHandler.Handle(auth, req)
}

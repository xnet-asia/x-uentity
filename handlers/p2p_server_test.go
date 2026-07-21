package handlers

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/xnetltd/x-uentity/repositories"
)

func TestP2PServerHandlesAnonymousRequest(t *testing.T) {
	auth := NewSimpleAuthHandler()
	entityHandler := NewEntityHandler(
		repositories.NewInMemoryRepository[handlerTestEntity](),
		NewMiddlewareChain[handlerTestEntity](),
	)
	server := NewP2PServer(entityHandler, auth)
	serverConn, clientConn := net.Pipe()

	done := make(chan error, 1)
	go func() { done <- server.ServeConnection(serverConn) }()

	request := map[string]any{
		"request": map[string]any{"action": "query"},
	}
	if err := json.NewEncoder(clientConn).Encode(request); err != nil {
		t.Fatal(err)
	}

	var response EntityResponse[handlerTestEntity]
	if err := json.NewDecoder(clientConn).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if !response.Success || response.Code != 200 {
		t.Fatalf("P2P response = %+v", response)
	}

	if err := clientConn.Close(); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

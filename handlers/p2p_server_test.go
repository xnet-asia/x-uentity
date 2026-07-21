package handlers

import (
	"encoding/json"
	"net"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/xnetltd/x-uentity/repositories"
)

type p2pAuthRecordingMiddleware struct {
	mu      sync.Mutex
	sources map[string]string
}

func (m *p2pAuthRecordingMiddleware) Ingress(auth *ClientAuth, req *EntityRequest[handlerTestEntity]) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sources[req.ID] = auth.Source
	return nil
}

func (m *p2pAuthRecordingMiddleware) Egress(_ *ClientAuth, _ *EntityResponse[handlerTestEntity]) error {
	return nil
}

func (m *p2pAuthRecordingMiddleware) source(id string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sources[id]
}

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

func TestP2PServerKeepsAnonymousSourcePerConnection(t *testing.T) {
	recorder := &p2pAuthRecordingMiddleware{sources: make(map[string]string)}
	entityHandler := NewEntityHandler(
		repositories.NewInMemoryRepository[handlerTestEntity](),
		NewMiddlewareChain[handlerTestEntity](recorder),
	)
	server := NewP2PServer(entityHandler, NewSimpleAuthHandler())

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	acceptDone := make(chan error, 1)
	serveDone := make(chan error, 2)
	go func() {
		for i := 0; i < 2; i++ {
			conn, err := listener.Accept()
			if err != nil {
				acceptDone <- err
				return
			}
			go func() { serveDone <- server.ServeConnection(conn) }()
		}
		acceptDone <- nil
	}()

	first, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer first.Close()
	second, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	if err := <-acceptDone; err != nil {
		t.Fatal(err)
	}

	sendAnonymousQuery := func(conn net.Conn, requestID string) {
		t.Helper()
		message := map[string]any{
			"request": map[string]any{"action": "query", "id": requestID},
		}
		if err := json.NewEncoder(conn).Encode(message); err != nil {
			t.Fatal(err)
		}
		var response EntityResponse[handlerTestEntity]
		if err := json.NewDecoder(conn).Decode(&response); err != nil {
			t.Fatal(err)
		}
		if !response.Success || response.Code != 200 {
			t.Fatalf("P2P response = %+v", response)
		}
	}

	sendAnonymousQuery(first, "first-1")
	sendAnonymousQuery(first, "first-2")
	sendAnonymousQuery(second, "second-1")
	sendAnonymousQuery(second, "second-2")
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	if err := second.Close(); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := <-serveDone; err != nil {
			t.Fatal(err)
		}
	}

	firstSource := recorder.source("first-1")
	secondSource := recorder.source("second-1")
	if _, err := uuid.Parse(firstSource); err != nil {
		t.Fatalf("first connection source %q is not a valid UUID: %v", firstSource, err)
	}
	if _, err := uuid.Parse(secondSource); err != nil {
		t.Fatalf("second connection source %q is not a valid UUID: %v", secondSource, err)
	}
	if recorder.source("first-2") != firstSource {
		t.Fatal("anonymous source changed within the first connection")
	}
	if recorder.source("second-2") != secondSource {
		t.Fatal("anonymous source changed within the second connection")
	}
	if firstSource == secondSource {
		t.Fatalf("anonymous connections share source %q", firstSource)
	}
	t.Logf("anonymous connection sources: first=%s second=%s", firstSource, secondSource)
}

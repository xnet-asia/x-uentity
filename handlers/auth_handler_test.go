package handlers

import (
	"testing"

	"github.com/google/uuid"
)

func TestSimpleAuthHandler(t *testing.T) {
	handler := NewSimpleAuthHandler()

	anon, err := handler.Authenticate("")
	if err != nil {
		t.Fatal(err)
	}
	if anon.IsAuth || anon.Source == "" {
		t.Fatalf("anonymous auth = %+v", anon)
	}
	if _, err := uuid.Parse(anon.Source); err != nil {
		t.Fatalf("anonymous source %q is not a valid UUID: %v", anon.Source, err)
	}

	secondAnon, err := handler.Authenticate("")
	if err != nil {
		t.Fatal(err)
	}
	if secondAnon.Source == anon.Source {
		t.Fatalf("anonymous sources are not unique: %q", anon.Source)
	}

	if _, err := handler.Register("client-1", "token-1"); err != nil {
		t.Fatal(err)
	}
	auth, err := handler.Authenticate("token-1")
	if err != nil {
		t.Fatal(err)
	}
	if !auth.IsAuth || auth.ID != "client-1" {
		t.Fatalf("registered auth = %+v", auth)
	}

	if err := handler.Revoke("client-1"); err != nil {
		t.Fatal(err)
	}
	auth, err = handler.Authenticate("token-1")
	if err != nil {
		t.Fatal(err)
	}
	if auth.IsAuth {
		t.Fatal("revoked token is still authenticated")
	}
	if auth.Source == "" {
		t.Fatal("revoked token did not receive an anonymous source")
	}
	if _, err := uuid.Parse(auth.Source); err != nil {
		t.Fatalf("anonymous source %q is not a valid UUID: %v", auth.Source, err)
	}
}

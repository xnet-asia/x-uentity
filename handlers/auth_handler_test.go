package handlers

import "testing"

func TestSimpleAuthHandler(t *testing.T) {
	handler := NewSimpleAuthHandler()

	anon, err := handler.Authenticate("")
	if err != nil {
		t.Fatal(err)
	}
	if anon.IsAuth || anon.ID != "anonymous" {
		t.Fatalf("anonymous auth = %+v", anon)
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
}

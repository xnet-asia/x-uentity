package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/xnetltd/x-uentity/factories"
	"github.com/xnetltd/x-uentity/handlers"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type entityHTTPMessage struct {
	Token   string                         `json:"token"`
	Request *handlers.EntityRequest[User] `json:"request"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	factory := factories.NewFactory[User]().WithMiddleware(
		&handlers.ValidationMiddleware[User]{},
		handlers.NewLoggingMiddleware[User](func(message string) {
			log.Printf("middleware %s", message)
		}),
	)

	clientID := os.Getenv("X_UENTITY_CLIENT_ID")
	clientToken := os.Getenv("X_UENTITY_CLIENT_TOKEN")
	if clientID != "" || clientToken != "" {
		if clientID == "" || clientToken == "" {
			log.Fatal("X_UENTITY_CLIENT_ID and X_UENTITY_CLIENT_TOKEN must be set together")
		}
		if _, err := factory.GetAuthHandler().Register(clientID, clientToken); err != nil {
			log.Fatalf("register auth client: %v", err)
		}
	}

	server := factory.GetP2PServer()
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/entities", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		defer r.Body.Close()

		var msg entityHTTPMessage
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		if msg.Request == nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "request is required"})
			return
		}
		msg.Request.Context = r.Context()

		resp, err := server.HandleHTTPFallback(msg.Token, msg.Request)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		status := resp.Code
		if status == 0 {
			status = http.StatusOK
		}
		writeJSON(w, status, resp)
	})

	addr := ":" + port
	log.Printf("x-uentity server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write response: %v", err)
	}
}

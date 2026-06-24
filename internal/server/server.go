// Package server wires the calc-api HTTP routes together. NewRouter is the
// single place routes are registered; endpoint subs extend it as they add the
// arithmetic handlers (/add, /subtract, /multiply, /divide).
package server

import (
	"net/http"

	"github.com/strategylippo/calc-api/internal/httpx"
)

// NewRouter constructs the application's HTTP handler with all routes
// registered. It returns an http.Handler so callers (main, tests) stay
// decoupled from the concrete mux implementation.
func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	return mux
}

// handleHealth reports that the service is up. It always returns
// HTTP 200 with the body {"status":"ok"}.
func handleHealth(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

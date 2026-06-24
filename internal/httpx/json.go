// Package httpx holds small shared helpers for writing JSON HTTP responses.
// Every calc-api endpoint encodes its output through here so success and
// error responses keep one consistent shape and Content-Type.
package httpx

import (
	"encoding/json"
	"net/http"
)

// WriteJSON marshals v and writes it as the response body with the given
// status code and an application/json Content-Type. It writes the exact
// JSON bytes (no trailing newline) so response bodies are predictable for
// clients and tests.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	body, err := json.Marshal(v)
	if err != nil {
		// Our own response types are expected to always marshal; if one
		// somehow doesn't, fail loudly with a 500 rather than emit a
		// half-written body with a success status.
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

// ErrorResponse is the JSON shape returned for every error condition:
// {"error": "..."}.
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteJSONError writes a JSON error body with the given status code. The
// endpoints use this for input validation and arithmetic failures so all
// error responses share the ErrorResponse shape.
func WriteJSONError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrorResponse{Error: message})
}

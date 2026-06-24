package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthEndpoint exercises GET /health through the real router using
// httptest, asserting the status code, Content-Type, and exact body required
// by the acceptance criteria.
func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	NewRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}
	if got := rec.Body.String(); got != `{"status":"ok"}` {
		t.Errorf("body = %q, want %q", got, `{"status":"ok"}`)
	}

	// The body must also decode to the documented shape.
	var payload map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("body is not valid JSON: %v", err)
	}
	if payload["status"] != "ok" {
		t.Errorf("status field = %q, want %q", payload["status"], "ok")
	}
}

// TestHealthRejectsNonGET confirms the method-scoped route does not answer
// other verbs, so later subs can register sibling methods on the same path
// without surprise.
func TestHealthRejectsNonGET(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rec := httptest.NewRecorder()

	NewRouter().ServeHTTP(rec, req)

	if rec.Code == http.StatusOK {
		t.Errorf("POST /health returned 200, want a non-200 (method not allowed)")
	}
}

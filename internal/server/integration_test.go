package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestEndToEnd boots the production router behind a real HTTP listener
// (httptest.NewServer) and exercises every endpoint over real network calls,
// proving the full path: request -> routing -> handler -> calc-lib -> JSON
// response. The in-process route tests (server_test.go, arithmetic_routes_test.go)
// use httptest.NewRequest and never cross a socket; this one does, so it is the
// assembled-service proof rather than a wiring check.
func TestEndToEnd(t *testing.T) {
	srv := httptest.NewServer(NewRouter())
	defer srv.Close()

	t.Run("success responses", func(t *testing.T) {
		cases := []struct {
			name     string
			path     string
			wantBody string
		}{
			{"health", "/health", `{"status":"ok"}`},
			{"add", "/add?a=2&b=3", `{"result":5}`},
			{"subtract", "/subtract?a=5&b=3", `{"result":2}`},
			{"multiply", "/multiply?a=4&b=3", `{"result":12}`},
			{"divide", "/divide?a=6&b=3", `{"result":2}`},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				status, ctype, body := get(t, srv.URL+tc.path)

				if status != http.StatusOK {
					t.Fatalf("status = %d, want %d (body %q)", status, http.StatusOK, body)
				}
				if ctype != "application/json" {
					t.Errorf("Content-Type = %q, want application/json", ctype)
				}
				if body != tc.wantBody {
					t.Errorf("body = %q, want %q", body, tc.wantBody)
				}
			})
		}
	})

	t.Run("error responses", func(t *testing.T) {
		cases := []struct {
			name string
			path string
		}{
			{"divide by zero", "/divide?a=1&b=0"},
			{"non-numeric operand", "/add?a=x&b=3"},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				status, ctype, body := get(t, srv.URL+tc.path)

				if status != http.StatusBadRequest {
					t.Fatalf("status = %d, want %d (body %q)", status, http.StatusBadRequest, body)
				}
				if ctype != "application/json" {
					t.Errorf("Content-Type = %q, want application/json", ctype)
				}

				// The body must decode to the documented {"error":"..."}
				// shape with a non-empty message.
				var payload struct {
					Error string `json:"error"`
				}
				if err := json.Unmarshal([]byte(body), &payload); err != nil {
					t.Fatalf("error body is not valid JSON: %v (body %q)", err, body)
				}
				if payload.Error == "" {
					t.Errorf("error message is empty (body %q)", body)
				}
			})
		}
	})
}

// get performs a real HTTP GET against url and returns the status code, the
// Content-Type header, and the response body as a string.
func get(t *testing.T, url string) (status int, contentType, body string) {
	t.Helper()

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body from %s: %v", url, err)
	}
	return resp.StatusCode, resp.Header.Get("Content-Type"), string(b)
}

package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestArithmeticRoutes confirms the four arithmetic endpoints are registered on
// the router from NewRouter and reach their handlers end to end. Handler-level
// edge cases (invalid params, divide-by-zero) live in internal/handlers; this
// test is about route wiring.
func TestArithmeticRoutes(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		wantBody string
	}{
		{"add", "/add?a=2&b=3", `{"result":5}`},
		{"subtract", "/subtract?a=5&b=3", `{"result":2}`},
		{"multiply", "/multiply?a=4&b=3", `{"result":12}`},
		{"divide", "/divide?a=6&b=3", `{"result":2}`},
	}

	router := NewRouter()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.target, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
			}
			if got := rec.Body.String(); got != tc.wantBody {
				t.Errorf("body = %q, want %q", got, tc.wantBody)
			}
		})
	}
}

// TestArithmeticRoutesRejectNonGET confirms the arithmetic routes are
// method-scoped to GET, consistent with /health.
func TestArithmeticRoutesRejectNonGET(t *testing.T) {
	router := NewRouter()
	for _, path := range []string{"/add", "/subtract", "/multiply", "/divide"} {
		req := httptest.NewRequest(http.MethodPost, path+"?a=1&b=2", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code == http.StatusOK {
			t.Errorf("POST %s returned 200, want a non-200 (method not allowed)", path)
		}
	}
}

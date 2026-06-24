package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/strategylippo/calc-lib/calc"
)

// TestArithmeticHandlers exercises each handler directly through httptest,
// asserting status, Content-Type, and exact JSON body for success, invalid /
// missing parameters, and divide-by-zero. The divide-by-zero expectation is
// built from calc.ErrDivideByZero so the test fails if the error stops
// originating from calc-lib.
func TestArithmeticHandlers(t *testing.T) {
	divZeroBody := `{"error":"` + calc.ErrDivideByZero.Error() + `"}`

	tests := []struct {
		name       string
		handler    http.HandlerFunc
		query      string
		wantStatus int
		wantBody   string
	}{
		// Success cases from the acceptance criteria.
		{"add success", Add, "a=2&b=3", http.StatusOK, `{"result":5}`},
		{"subtract success", Subtract, "a=5&b=3", http.StatusOK, `{"result":2}`},
		{"multiply success", Multiply, "a=4&b=3", http.StatusOK, `{"result":12}`},
		{"divide success", Divide, "a=6&b=3", http.StatusOK, `{"result":2}`},

		// Non-integer operands and results must round-trip cleanly.
		{"add floats", Add, "a=2.5&b=0.5", http.StatusOK, `{"result":3}`},
		{"divide fractional", Divide, "a=7&b=2", http.StatusOK, `{"result":3.5}`},
		{"subtract negative result", Subtract, "a=3&b=5", http.StatusOK, `{"result":-2}`},

		// Invalid parameter: non-numeric value.
		{"add invalid a", Add, "a=x&b=3", http.StatusBadRequest, `{"error":"query parameter a must be a number"}`},
		{"multiply invalid b", Multiply, "a=4&b=nope", http.StatusBadRequest, `{"error":"query parameter b must be a number"}`},

		// Missing parameters.
		{"add missing b", Add, "a=2", http.StatusBadRequest, `{"error":"missing required query parameter b"}`},
		{"subtract missing a", Subtract, "b=3", http.StatusBadRequest, `{"error":"missing required query parameter a"}`},

		// Non-finite operands: strconv.ParseFloat accepts Inf/NaN but json.Marshal
		// rejects them, so they must be caught at parse time as a 400 (not a 500).
		{"add inf operand", Add, "a=Inf&b=1", http.StatusBadRequest, `{"error":"query parameter a must be a finite number"}`},
		{"subtract nan operand", Subtract, "a=1&b=NaN", http.StatusBadRequest, `{"error":"query parameter b must be a finite number"}`},

		// Finite operands whose result overflows to +Inf must also map to 400,
		// never a 500 with a non-JSON body.
		{"multiply overflow result", Multiply, "a=1e308&b=1e308", http.StatusBadRequest, `{"error":"result is not a finite number"}`},

		// Divide-by-zero maps to 400 with calc-lib's error.
		{"divide by zero", Divide, "a=1&b=0", http.StatusBadRequest, divZeroBody},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/?"+tc.query, nil)
			rec := httptest.NewRecorder()

			tc.handler(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want %q", ct, "application/json")
			}
			if got := rec.Body.String(); got != tc.wantBody {
				t.Errorf("body = %q, want %q", got, tc.wantBody)
			}
		})
	}
}

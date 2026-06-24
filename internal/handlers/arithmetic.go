// Package handlers holds the arithmetic HTTP handlers (/add, /subtract,
// /multiply, /divide). Each handler parses the `a` and `b` query parameters as
// numbers, delegates the actual arithmetic to github.com/strategylippo/calc-lib
// (no local +,-,*,/ on operands), and writes the result as JSON {"result":N}.
// Invalid or missing parameters, and calc-lib's ErrDivideByZero, map to HTTP 400
// JSON error bodies via the shared httpx helper.
package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/strategylippo/calc-api/internal/httpx"
	"github.com/strategylippo/calc-lib/calc"
)

// result is the success response shape: {"result":N}.
type result struct {
	Result float64 `json:"result"`
}

// Add handles GET /add: returns a+b.
func Add(w http.ResponseWriter, r *http.Request) {
	a, b, ok := parseOperands(w, r)
	if !ok {
		return
	}
	writeResult(w, calc.Add(a, b))
}

// Subtract handles GET /subtract: returns a-b.
func Subtract(w http.ResponseWriter, r *http.Request) {
	a, b, ok := parseOperands(w, r)
	if !ok {
		return
	}
	writeResult(w, calc.Subtract(a, b))
}

// Multiply handles GET /multiply: returns a*b.
func Multiply(w http.ResponseWriter, r *http.Request) {
	a, b, ok := parseOperands(w, r)
	if !ok {
		return
	}
	writeResult(w, calc.Multiply(a, b))
}

// Divide handles GET /divide: returns a/b, or a 400 JSON error when calc-lib
// reports a divide-by-zero. The error provenance is calc.ErrDivideByZero, not a
// locally-invented check.
func Divide(w http.ResponseWriter, r *http.Request) {
	a, b, ok := parseOperands(w, r)
	if !ok {
		return
	}
	q, err := calc.Divide(a, b)
	if err != nil {
		httpx.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeResult(w, q)
}

// writeResult writes a successful {"result":N} response, but first guards
// against non-finite results. json.Marshal cannot encode NaN or ±Inf, so a
// result that overflowed (e.g. two large finite operands multiplying to +Inf)
// would otherwise surface as a 500 with a non-JSON body. Mapping it to a 400
// JSON error keeps the failure client-visible and consistent in shape with the
// other validation errors.
func writeResult(w http.ResponseWriter, v float64) {
	if math.IsInf(v, 0) || math.IsNaN(v) {
		httpx.WriteJSONError(w, http.StatusBadRequest, "result is not a finite number")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, result{Result: v})
}

// parseOperands reads the required `a` and `b` query parameters and parses them
// as float64. On the first missing or non-numeric value it writes a 400 JSON
// error via the shared helper and returns ok=false; callers must return
// immediately without writing a second response.
func parseOperands(w http.ResponseWriter, r *http.Request) (a, b float64, ok bool) {
	a, ok = parseParam(w, r, "a")
	if !ok {
		return 0, 0, false
	}
	b, ok = parseParam(w, r, "b")
	if !ok {
		return 0, 0, false
	}
	return a, b, true
}

// parseParam reads a single required numeric query parameter. A missing value,
// a non-numeric value, and a non-finite value (NaN/±Inf, which strconv accepts
// but json.Marshal rejects) each produce a distinct 400 JSON error message.
func parseParam(w http.ResponseWriter, r *http.Request, name string) (float64, bool) {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		httpx.WriteJSONError(w, http.StatusBadRequest, "missing required query parameter "+name)
		return 0, false
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		httpx.WriteJSONError(w, http.StatusBadRequest, "query parameter "+name+" must be a number")
		return 0, false
	}
	if math.IsInf(v, 0) || math.IsNaN(v) {
		httpx.WriteJSONError(w, http.StatusBadRequest, "query parameter "+name+" must be a finite number")
		return 0, false
	}
	return v, true
}

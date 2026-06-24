package main

import "testing"

// TestPortDefault covers the PORT-unset case: the server must listen on 8080.
// An empty PORT exercises the same fallback branch as a fully unset variable.
func TestPortDefault(t *testing.T) {
	t.Setenv("PORT", "")
	if got := port(); got != "8080" {
		t.Errorf("port() with PORT empty/unset = %q, want %q", got, "8080")
	}
}

// TestPortFromEnv covers PORT=9090: the server must listen on the supplied
// port.
func TestPortFromEnv(t *testing.T) {
	t.Setenv("PORT", "9090")
	if got := port(); got != "9090" {
		t.Errorf("port() with PORT=9090 = %q, want %q", got, "9090")
	}
}

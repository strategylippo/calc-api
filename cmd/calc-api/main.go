// Command calc-api starts the calc-api HTTP server: it reads the listen port
// from the environment, builds the router, and serves until the process is
// stopped.
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/strategylippo/calc-api/internal/server"
)

func main() {
	addr := ":" + port()
	// Explicit timeouts bound how long a single connection can tie up server
	// resources, closing the door on slow-client (Slowloris) exhaustion that a
	// bare http.ListenAndServe leaves wide open.
	srv := &http.Server{
		Addr:              addr,
		Handler:           server.NewRouter(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("calc-api listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("calc-api: server stopped: %v", err)
	}
}

// port returns the TCP port to listen on. It uses the PORT environment
// variable and falls back to 8080 when PORT is unset or empty.
func port() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}

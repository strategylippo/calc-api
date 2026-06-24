// Command calc-api starts the calc-api HTTP server: it reads the listen port
// from the environment, builds the router, and serves until the process is
// stopped.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/strategylippo/calc-api/internal/server"
)

func main() {
	addr := ":" + port()
	log.Printf("calc-api listening on %s", addr)
	if err := http.ListenAndServe(addr, server.NewRouter()); err != nil {
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

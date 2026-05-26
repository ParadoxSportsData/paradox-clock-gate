package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ParadoxSportsData/paradox-clock-gate/internal/server"
)

// runServe parses serve-mode flags and starts the HTTP server.
// It is called when os.Args[1] == "serve".
func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.Int("port", 8080, "port to listen on")
	dataDir := fs.String("data", "./testdata", "directory containing game JSON files")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "serve: %v\n", err)
		os.Exit(1)
	}

	// Validate data directory exists before starting server.
	if _, err := os.Stat(*dataDir); err != nil {
		fmt.Fprintf(os.Stderr, "serve: data directory %q not found: %v\n", *dataDir, err)
		os.Exit(1)
	}

	cache := server.NewGameCache(*dataDir)
	mux := server.NewServeMux(cache)
	handler := corsMiddleware(mux)

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("clock-gate serve listening on %s\n", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

// corsMiddleware wraps an http.Handler to add CORS headers on every response
// and handle OPTIONS preflight requests with a 204 No Content.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

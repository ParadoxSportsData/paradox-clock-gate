package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ParadoxSportsData/paradox-clock-gate/internal/server"
)

// runServe parses serve-mode flags and starts the HTTP server.
// It is called when os.Args[1] == "serve".
func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.Int("port", 8080, "port to listen on")
	dataDir := fs.String("data", "./data/raw", "directory containing game JSON files")
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

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("clock-gate serve listening on %s\n", addr)
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ParadoxSportsData/paradox-clock-gate/internal/server"
)

// newTestMux returns a ServeMux backed by a GameCache with no data directory,
// sufficient for testing HTTP-level concerns like CORS headers and routing.
func newTestMux() http.Handler {
	cache := server.NewGameCache(".")
	return server.NewServeMux(cache)
}

// TestCORSMiddleware_SetsAllowOriginHeader verifies that the CORS middleware
// (applied inside server.NewServeMux) injects Access-Control-Allow-Origin: *
// on every response.
func TestCORSMiddleware_SetsAllowOriginHeader(t *testing.T) {
	handler := newTestMux()

	req := httptest.NewRequest(http.MethodGet, "/games", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	got := w.Header().Get("Access-Control-Allow-Origin")
	if got != "*" {
		t.Errorf("expected Access-Control-Allow-Origin: *, got %q", got)
	}
}

// TestCORSMiddleware_OptionsReturns204 verifies OPTIONS preflight returns 204
// without forwarding to inner handlers.
func TestCORSMiddleware_OptionsReturns204(t *testing.T) {
	handler := newTestMux()

	req := httptest.NewRequest(http.MethodOptions, "/games", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204 for OPTIONS, got %d", w.Code)
	}
}

// TestCORSMiddleware_PassesThroughNonOptions verifies that GET requests are
// forwarded to the inner handler (non-OPTIONS path is not short-circuited).
func TestCORSMiddleware_PassesThroughNonOptions(t *testing.T) {
	handler := newTestMux()

	req := httptest.NewRequest(http.MethodGet, "/games", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// The handler should respond (200 or any non-options code) — not 204.
	if w.Code == http.StatusNoContent {
		t.Error("GET request must not be treated as OPTIONS preflight")
	}
}

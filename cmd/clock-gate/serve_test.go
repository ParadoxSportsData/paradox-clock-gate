package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCORSMiddleware_SetsAllowOriginHeader verifies corsMiddleware injects
// Access-Control-Allow-Origin: * on every response.
func TestCORSMiddleware_SetsAllowOriginHeader(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := corsMiddleware(inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	got := w.Header().Get("Access-Control-Allow-Origin")
	if got != "*" {
		t.Errorf("expected Access-Control-Allow-Origin: *, got %q", got)
	}
}

// TestCORSMiddleware_OptionsReturns204 verifies OPTIONS preflight returns 204
// without calling the inner handler.
func TestCORSMiddleware_OptionsReturns204(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If this runs the test should fail — preflight must be short-circuited.
		w.WriteHeader(http.StatusInternalServerError)
	})
	handler := corsMiddleware(inner)

	req := httptest.NewRequest(http.MethodOptions, "/games", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204 for OPTIONS, got %d", w.Code)
	}
}

// TestCORSMiddleware_PassesThroughNonOptions verifies that non-OPTIONS requests
// are forwarded to the inner handler.
func TestCORSMiddleware_PassesThroughNonOptions(t *testing.T) {
	called := false
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	handler := corsMiddleware(inner)

	req := httptest.NewRequest(http.MethodGet, "/games", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if !called {
		t.Error("expected inner handler to be called for GET request")
	}
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

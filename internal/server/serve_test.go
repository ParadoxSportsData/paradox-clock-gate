package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// testServer returns an httptest.Server backed by the real GameCache+Mux.
// Caller must call server.Close() via defer.
func testServer(t *testing.T) *httptest.Server {
	t.Helper()
	cache := NewGameCache("../../testdata")
	mux := NewServeMux(cache)
	return httptest.NewServer(mux)
}

// TestServeListGames verifies GET /games returns a non-empty list of game summaries.
func TestServeListGames(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/games")
	if err != nil {
		t.Fatalf("GET /games: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var games []GameSummary
	if err := json.NewDecoder(resp.Body).Decode(&games); err != nil {
		t.Fatalf("decode GameSummary slice: %v", err)
	}
	if len(games) < 1 {
		t.Errorf("expected at least 1 game, got %d", len(games))
	}
}

// TestServeStateAtTick verifies GET /games/2011_01_NO_GB/state?tick=1800 returns quarter 3.
func TestServeStateAtTick(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/games/2011_01_NO_GB/state?tick=1800")
	if err != nil {
		t.Fatalf("GET /games/2011_01_NO_GB/state: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var gs GameStateResponse
	if err := json.NewDecoder(resp.Body).Decode(&gs); err != nil {
		t.Fatalf("decode GameStateResponse: %v", err)
	}
	if gs.Quarter != 3 {
		t.Errorf("expected quarter 3 at tick 1800, got %d", gs.Quarter)
	}
}

// TestServeStateTickTooHigh verifies GET /games/2011_01_NO_GB/state?tick=999999
// returns 422 Unprocessable Entity with a non-nil max_tick in the error body.
func TestServeStateTickTooHigh(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/games/2011_01_NO_GB/state?tick=999999")
	if err != nil {
		t.Fatalf("GET /games/2011_01_NO_GB/state?tick=999999: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", resp.StatusCode)
	}

	var er ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&er); err != nil {
		t.Fatalf("decode ErrorResponse: %v", err)
	}
	if er.MaxTick == nil {
		t.Error("expected max_tick to be non-nil in 422 response")
	}
}

// TestServeTimeline verifies GET /games/2011_01_NO_GB/timeline returns > 50 plays.
func TestServeTimeline(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/games/2011_01_NO_GB/timeline")
	if err != nil {
		t.Fatalf("GET /games/2011_01_NO_GB/timeline: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var timeline GameTimelineResponse
	if err := json.NewDecoder(resp.Body).Decode(&timeline); err != nil {
		t.Fatalf("decode GameTimelineResponse: %v", err)
	}
	if len(timeline.Plays) <= 50 {
		t.Errorf("expected > 50 plays, got %d", len(timeline.Plays))
	}
}

// TestServeCORSHeader verifies every response carries Access-Control-Allow-Origin: *.
func TestServeCORSHeader(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/games")
	if err != nil {
		t.Fatalf("GET /games: %v", err)
	}
	defer resp.Body.Close()

	got := resp.Header.Get("Access-Control-Allow-Origin")
	if got != "*" {
		t.Errorf("expected Access-Control-Allow-Origin: *, got %q", got)
	}
}

// TestServePathTraversal verifies that a path-traversal game ID returns 404.
func TestServePathTraversal(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	// Use the raw URL string; net/http client will keep the encoded form.
	resp, err := http.Get(srv.URL + "/games/..%2F..%2Fetc%2Fpasswd/state?tick=0")
	if err != nil {
		t.Fatalf("GET path traversal: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 for path traversal game ID, got %d", resp.StatusCode)
	}
}

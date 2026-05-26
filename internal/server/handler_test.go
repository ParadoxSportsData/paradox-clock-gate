package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// testCache creates a GameCache pointing at the real testdata directory.
func testCache(t *testing.T) *GameCache {
	t.Helper()
	return NewGameCache("../../testdata")
}

// TestHandleGames_Returns200WithGames verifies GET /games returns HTTP 200 and
// a non-empty JSON array of GameSummary objects.
func TestHandleGames_Returns200WithGames(t *testing.T) {
	mux := NewServeMux(testCache(t))
	req := httptest.NewRequest(http.MethodGet, "/games", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var summaries []GameSummary
	if err := json.Unmarshal(w.Body.Bytes(), &summaries); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(summaries) < 1 {
		t.Errorf("expected ≥1 game, got %d", len(summaries))
	}
}

// TestHandleState_Returns200WithCorrectQuarter verifies GET /games/{id}/state?tick=1800
// returns HTTP 200 and quarter=3 for 2011_01_NO_GB at tick 1800.
func TestHandleState_Returns200WithCorrectQuarter(t *testing.T) {
	mux := NewServeMux(testCache(t))
	req := httptest.NewRequest(http.MethodGet, "/games/2011_01_NO_GB/state?tick=1800", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: body=%s", w.Code, w.Body.String())
	}
	var resp GameStateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.Quarter != 3 {
		t.Errorf("expected quarter=3 at tick 1800, got %d", resp.Quarter)
	}
	if !resp.HasState {
		t.Error("expected has_state=true at tick 1800")
	}
}

// TestHandleState_Returns422WhenTickTooHigh verifies that requesting a tick beyond
// maxTick returns HTTP 422 with a max_tick field in the error body.
func TestHandleState_Returns422WhenTickTooHigh(t *testing.T) {
	mux := NewServeMux(testCache(t))
	req := httptest.NewRequest(http.MethodGet, "/games/2011_01_NO_GB/state?tick=999999", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
	var errResp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if errResp.MaxTick == nil {
		t.Error("expected max_tick field in 422 response, got nil")
	}
}

// TestHandleState_Returns400WhenTickMissing verifies that omitting the tick
// query parameter returns HTTP 400.
func TestHandleState_Returns400WhenTickMissing(t *testing.T) {
	mux := NewServeMux(testCache(t))
	req := httptest.NewRequest(http.MethodGet, "/games/2011_01_NO_GB/state", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

// TestHandleState_Returns400WhenTickNotInt verifies that a non-integer tick
// parameter returns HTTP 400.
func TestHandleState_Returns400WhenTickNotInt(t *testing.T) {
	mux := NewServeMux(testCache(t))
	req := httptest.NewRequest(http.MethodGet, "/games/2011_01_NO_GB/state?tick=abc", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

// TestHandleState_Returns404ForUnknownGame verifies that requesting state for a
// game that does not exist returns HTTP 404.
func TestHandleState_Returns404ForUnknownGame(t *testing.T) {
	mux := NewServeMux(testCache(t))
	req := httptest.NewRequest(http.MethodGet, "/games/nonexistent_game/state?tick=0", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// TestHandleTimeline_Returns200WithPlays verifies GET /games/{id}/timeline
// returns HTTP 200 and a non-empty plays array.
func TestHandleTimeline_Returns200WithPlays(t *testing.T) {
	mux := NewServeMux(testCache(t))
	req := httptest.NewRequest(http.MethodGet, "/games/2011_01_NO_GB/timeline", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: body=%s", w.Code, w.Body.String())
	}
	var resp GameTimelineResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.MaxTick == 0 {
		t.Error("expected max_tick > 0")
	}
	if len(resp.Plays) < 50 {
		t.Errorf("expected ≥50 plays, got %d", len(resp.Plays))
	}
}

// TestHandleResponse_HasCORSHeader verifies that all responses include the
// Access-Control-Allow-Origin: * header.
func TestHandleResponse_HasCORSHeader(t *testing.T) {
	mux := NewServeMux(testCache(t))
	req := httptest.NewRequest(http.MethodGet, "/games", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	cors := w.Header().Get("Access-Control-Allow-Origin")
	if cors != "*" {
		t.Errorf("expected Access-Control-Allow-Origin: *, got %q", cors)
	}
}

// TestHandleState_ContentTypeIsJSON verifies that state responses have
// Content-Type: application/json.
func TestHandleState_ContentTypeIsJSON(t *testing.T) {
	mux := NewServeMux(testCache(t))
	req := httptest.NewRequest(http.MethodGet, "/games/2011_01_NO_GB/state?tick=0", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type: application/json, got %q", ct)
	}
}

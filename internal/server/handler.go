// Package server implements the HTTP REST API for clock-gate, exposing game
// state and timeline endpoints backed by an in-memory GameCache.
package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ParadoxSportsData/paradox-clock-gate/internal/matrix"
)

// NewServeMux wires the three REST endpoints onto a new ServeMux.
// All routes have CORS headers injected via a thin middleware wrapper.
func NewServeMux(cache *GameCache) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/games", corsMiddlewareHandler(handleGames(cache)))
	// /games/ catches /games/{id}/state and /games/{id}/timeline.
	mux.Handle("/games/", corsMiddlewareHandler(routeGameSub(cache)))
	return mux
}

// corsMiddlewareHandler sets CORS headers on every response, handles OPTIONS preflight.
func corsMiddlewareHandler(next http.Handler) http.Handler {
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

// routeGameSub dispatches /games/{id}/state and /games/{id}/timeline.
func routeGameSub(cache *GameCache) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path // e.g. /games/2011_01_NO_GB/state
		switch {
		case strings.HasSuffix(path, "/state"):
			handleState(cache).ServeHTTP(w, r)
		case strings.HasSuffix(path, "/timeline"):
			handleTimeline(cache).ServeHTTP(w, r)
		default:
			writeError(w, http.StatusNotFound, "not found")
		}
	})
}

// handleGames lists all games from the cache.
func handleGames(cache *GameCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		summaries, err := cache.ListGames()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, summaries)
	}
}

// handleState returns the O(1) StateMatrix lookup for a single tick.
func handleState(cache *GameCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := extractGameID(r.URL.Path, "/state")
		if gameID == "" {
			writeError(w, http.StatusNotFound, "game not found")
			return
		}

		tickStr := r.URL.Query().Get("tick")
		if tickStr == "" {
			writeError(w, http.StatusBadRequest, "tick query parameter required")
			return
		}
		tick, err := strconv.Atoi(tickStr)
		if err != nil || tick < 0 {
			writeError(w, http.StatusBadRequest, "tick must be a non-negative integer")
			return
		}

		sm, err := cache.Load(gameID)
		if err != nil {
			writeError(w, http.StatusNotFound, "game not found: "+gameID)
			return
		}

		if uint16(tick) > sm.Meta.MaxTick {
			maxT := sm.Meta.MaxTick
			writeJSON(w, http.StatusUnprocessableEntity, ErrorResponse{
				Error:   "tick exceeds game duration",
				MaxTick: &maxT,
			})
			return
		}

		gs := sm.States[tick]
		resp := gameStateToResponse(gs, sm)
		writeJSON(w, http.StatusOK, resp)
	}
}

// handleTimeline returns all HasState ticks as a sorted PlaySnapshot slice.
func handleTimeline(cache *GameCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := extractGameID(r.URL.Path, "/timeline")
		if gameID == "" {
			writeError(w, http.StatusNotFound, "game not found")
			return
		}

		sm, err := cache.Load(gameID)
		if err != nil {
			writeError(w, http.StatusNotFound, "game not found: "+gameID)
			return
		}

		plays := make([]PlaySnapshot, 0, 200)
		for t := 0; t <= int(sm.Meta.MaxTick); t++ {
			gs := sm.States[t]
			if !gs.HasState {
				continue
			}
			plays = append(plays, PlaySnapshot{
				Tick:        gs.Elapsed,
				Quarter:     gs.Quarter,
				Down:        nullableUint8(gs.Down),
				YardsToGo:   nullableUint8(gs.YardsToGo),
				YardLine:    nullableUint8(gs.YardLine),
				HomeScore:   gs.HomeScore,
				AwayScore:   gs.AwayScore,
				Posteam:     nullableTeam(gs.Posteam),
				WinProb:     nullableWinProb(gs.WinProb),
				PlayType:    playTypeStr(gs.PlayType),
				Description: descFromArena(gs, sm.Arena),
			})
		}

		resp := GameTimelineResponse{
			GameID:   sm.Meta.GameID,
			HomeTeam: sm.Meta.HomeTeam,
			AwayTeam: sm.Meta.AwayTeam,
			MaxTick:  sm.Meta.MaxTick,
			Plays:    plays,
		}
		writeJSON(w, http.StatusOK, resp)
	}
}

// writeJSON sets Content-Type and encodes v as JSON with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}

// writeError sends a JSON ErrorResponse with the given status and message.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

// extractGameID extracts the game ID from a URL path of the form
// /games/{id}/{suffix}.
func extractGameID(path, suffix string) string {
	// Strip /games/ prefix and suffix.
	trimmed := strings.TrimPrefix(path, "/games/")
	trimmed = strings.TrimSuffix(trimmed, suffix)
	trimmed = strings.Trim(trimmed, "/")
	return trimmed
}

// gameStateToResponse converts a GameState to a GameStateResponse.
func gameStateToResponse(gs matrix.GameState, sm *matrix.StateMatrix) GameStateResponse {
	return GameStateResponse{
		Tick:        gs.Elapsed,
		Quarter:     gs.Quarter,
		Down:        nullableUint8(gs.Down),
		YardsToGo:   nullableUint8(gs.YardsToGo),
		YardLine:    nullableUint8(gs.YardLine),
		HomeScore:   gs.HomeScore,
		AwayScore:   gs.AwayScore,
		Posteam:     nullableTeam(gs.Posteam),
		Defteam:     nullableTeam(gs.Defteam),
		WinProb:     nullableWinProb(gs.WinProb),
		PlayType:    playTypeStr(gs.PlayType),
		Description: descFromArena(gs, sm.Arena),
		HasState:    gs.HasState,
	}
}

// nullableUint8 returns nil if v == 0 (no down, no yards etc.), else &v.
// Down and YardsToGo are 0 on kickoffs/special teams — those should be null.
func nullableUint8(v uint8) *uint8 {
	if v == 0 {
		return nil
	}
	return &v
}

// nullableTeam converts a [3]byte team abbreviation to *string.
// Returns nil if all bytes are zero.
func nullableTeam(b [3]byte) *string {
	end := 0
	for end < 3 && b[end] != 0 {
		end++
	}
	if end == 0 {
		return nil
	}
	s := string(b[:end])
	return &s
}

// nullableWinProb converts the uint16 WinProb sentinel to *float64.
// matrix.WinProbNull is the null sentinel — returns nil. Otherwise wp/WinProbScale.
func nullableWinProb(wp uint16) *float64 {
	if wp == matrix.WinProbNull {
		return nil
	}
	v := float64(wp) / float64(matrix.WinProbScale)
	return &v
}

// playTypeStr converts the PlayType enum to its string representation.
func playTypeStr(pt matrix.PlayType) string {
	switch pt {
	case matrix.PlayTypeRun:
		return "run"
	case matrix.PlayTypePass:
		return "pass"
	case matrix.PlayTypePunt:
		return "punt"
	case matrix.PlayTypeKickoff:
		return "kickoff"
	case matrix.PlayTypeNoPlay:
		return "no_play"
	case matrix.PlayTypeOther:
		return "other"
	default:
		return ""
	}
}

// descFromArena extracts a play description from the arena byte slice.
func descFromArena(gs matrix.GameState, arena []byte) string {
	end := int(gs.DescOffset) + int(gs.DescLen)
	if gs.DescLen == 0 || end > len(arena) {
		return ""
	}
	return string(arena[gs.DescOffset:end])
}

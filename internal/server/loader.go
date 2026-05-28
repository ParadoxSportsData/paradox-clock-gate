package server

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/ParadoxSportsData/paradox-clock-gate/internal/ingestion"
	"github.com/ParadoxSportsData/paradox-clock-gate/internal/matrix"
)

// GameCache holds compiled StateMatrices keyed by game ID.
// Load is lazy (on first request) and cached forever — game data is immutable.
// A sync.RWMutex protects the map: concurrent readers never block each other;
// write lock is held only during the map write, NOT during slow parse/compile.
type GameCache struct {
	mu      sync.RWMutex
	games   map[string]*matrix.StateMatrix
	dataDir string
}

// NewGameCache creates an empty GameCache rooted at dataDir.
func NewGameCache(dataDir string) *GameCache {
	abs, err := filepath.Abs(dataDir)
	if err != nil {
		abs = dataDir
	}
	return &GameCache{
		games:   make(map[string]*matrix.StateMatrix),
		dataDir: abs,
	}
}

// Load returns the compiled StateMatrix for gameID, parsing and compiling on
// first access. Subsequent calls return the cached pointer.
//
// Path traversal guard: gameID is cleaned and resolved; if the resulting path
// does not have dataDir as a prefix, an error is returned.
func (c *GameCache) Load(gameID string) (*matrix.StateMatrix, error) {
	// Sanitize: build the expected file path and verify it stays inside dataDir.
	candidate := filepath.Clean(filepath.Join(c.dataDir, gameID+".json"))
	// Ensure dataDir ends with separator for prefix check (prevents e.g. /foo matching /foobar).
	safeDir := c.dataDir
	if !strings.HasSuffix(safeDir, string(filepath.Separator)) {
		safeDir += string(filepath.Separator)
	}
	if !strings.HasPrefix(candidate, safeDir) {
		return nil, fmt.Errorf("game ID %q escapes data directory", gameID)
	}

	// Fast path: read lock, check map.
	c.mu.RLock()
	sm, ok := c.games[gameID]
	c.mu.RUnlock()
	if ok {
		return sm, nil
	}

	// Slow path: parse and compile outside any lock.
	header, plays, err := ingestion.ParseFile(candidate)
	if err != nil {
		// Distinguish file-not-found from parse error for cleaner error messages.
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("game %q not found", gameID)
		}
		return nil, fmt.Errorf("parse %q: %w", gameID, err)
	}
	compiled := matrix.Compile(plays, header)
	newSM := &compiled

	// Write lock: double-check in case another goroutine loaded the same game
	// while we were compiling.
	c.mu.Lock()
	if existing, ok := c.games[gameID]; ok {
		c.mu.Unlock()
		return existing, nil
	}
	c.games[gameID] = newSM
	c.mu.Unlock()

	return newSM, nil
}

// ListGames scans dataDir for *.json files, loads each via Load, and returns
// a []GameSummary sorted ascending by GameID.
func (c *GameCache) ListGames() ([]GameSummary, error) {
	pattern := filepath.Join(c.dataDir, "*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob %s: %w", pattern, err)
	}

	summaries := make([]GameSummary, 0, len(matches))
	for _, path := range matches {
		base := filepath.Base(path)
		gameID := strings.TrimSuffix(base, ".json")

		sm, err := c.Load(gameID)
		if err != nil {
			// Skip files that can't be loaded (e.g. malformed JSON) rather than
			// failing the entire list.
			continue
		}

		// Skip files whose internal game_id doesn't match the filename.
		// Canonical game files always have a matching filename and JSON game_id.
		// Sample or mis-named files (e.g. 2011_01_NO_GB_sample.json whose JSON
		// contains game_id "2011_01_NO_GB") would otherwise create phantom duplicates.
		if sm.Meta.GameID != gameID {
			continue
		}

		// Final score comes from the state at MaxTick.
		finalState := sm.States[sm.Meta.MaxTick]
		summaries = append(summaries, GameSummary{
			GameID:    gameID,
			HomeTeam:  sm.Meta.HomeTeam,
			AwayTeam:  sm.Meta.AwayTeam,
			HomeScore: finalState.HomeScore,
			AwayScore: finalState.AwayScore,
			Duration:  sm.Meta.MaxTick,
		})
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].GameID < summaries[j].GameID
	})
	return summaries, nil
}

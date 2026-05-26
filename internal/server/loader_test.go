package server

import (
	"testing"
)

func TestNewGameCache(t *testing.T) {
	cache := NewGameCache("../../testdata")
	if cache == nil {
		t.Fatal("NewGameCache returned nil")
	}
}

func TestLoadGame_HomeTeam(t *testing.T) {
	cache := NewGameCache("../../testdata")
	sm, err := cache.Load("2011_01_NO_GB")
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if sm == nil {
		t.Fatal("Load returned nil StateMatrix")
	}
	if sm.Meta.HomeTeam != "GB" {
		t.Errorf("expected HomeTeam=GB, got %q", sm.Meta.HomeTeam)
	}
}

func TestLoadGame_CacheHit(t *testing.T) {
	cache := NewGameCache("../../testdata")
	sm1, err := cache.Load("2011_01_NO_GB")
	if err != nil {
		t.Fatalf("first Load error: %v", err)
	}
	sm2, err := cache.Load("2011_01_NO_GB")
	if err != nil {
		t.Fatalf("second Load error: %v", err)
	}
	if sm1 != sm2 {
		t.Error("expected same pointer on cache hit, got different pointers")
	}
}

func TestLoadGame_PathTraversal(t *testing.T) {
	cache := NewGameCache("../../testdata")
	_, err := cache.Load("../../etc/passwd")
	if err == nil {
		t.Error("expected error for path traversal attempt, got nil")
	}
}

func TestLoadGame_NotFound(t *testing.T) {
	cache := NewGameCache("../../testdata")
	_, err := cache.Load("nonexistent_game")
	if err == nil {
		t.Error("expected error for missing game, got nil")
	}
}

func TestListGames(t *testing.T) {
	cache := NewGameCache("../../testdata")
	summaries, err := cache.ListGames()
	if err != nil {
		t.Fatalf("ListGames error: %v", err)
	}
	if len(summaries) < 1 {
		t.Errorf("expected at least 1 game summary, got %d", len(summaries))
	}
	// Verify known game is present.
	found := false
	for _, s := range summaries {
		if s.GameID == "2011_01_NO_GB" {
			found = true
			if s.HomeTeam != "GB" {
				t.Errorf("expected HomeTeam=GB, got %q", s.HomeTeam)
			}
			if s.AwayTeam != "NO" {
				t.Errorf("expected AwayTeam=NO, got %q", s.AwayTeam)
			}
			if s.Duration == 0 {
				t.Error("expected Duration > 0")
			}
		}
	}
	if !found {
		t.Error("2011_01_NO_GB not found in ListGames output")
	}
}

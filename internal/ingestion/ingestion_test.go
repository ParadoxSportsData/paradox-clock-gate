package ingestion

import "testing"

func TestGameHeaderFields(t *testing.T) {
	h := GameHeader{
		GameID:    "2011_01_ATL_CHI",
		HomeTeam:  "CHI",
		AwayTeam:  "ATL",
		HomeScore: 30,
		AwayScore: 12,
	}
	if h.GameID != "2011_01_ATL_CHI" {
		t.Fatalf("GameID = %q, want 2011_01_ATL_CHI", h.GameID)
	}
	if h.HomeTeam != "CHI" {
		t.Fatalf("HomeTeam = %q, want CHI", h.HomeTeam)
	}
	if h.AwayTeam != "ATL" {
		t.Fatalf("AwayTeam = %q, want ATL", h.AwayTeam)
	}
	if h.HomeScore != 30 {
		t.Fatalf("HomeScore = %d, want 30", h.HomeScore)
	}
	if h.AwayScore != 12 {
		t.Fatalf("AwayScore = %d, want 12", h.AwayScore)
	}
}

func TestRawPlayNullableFields(t *testing.T) {
	// Nullable JSON fields must be pointer types so they can represent null.
	var p RawPlay
	if p.Down != nil {
		t.Fatal("Down must be nil when not set")
	}
	if p.YardsToGo != nil {
		t.Fatal("YardsToGo must be nil when not set")
	}
	if p.YardLine100 != nil {
		t.Fatal("YardLine100 must be nil when not set")
	}
	if p.PlayType != nil {
		t.Fatal("PlayType must be nil when not set")
	}
	if p.Posteam != nil {
		t.Fatal("Posteam must be nil when not set")
	}
	if p.Defteam != nil {
		t.Fatal("Defteam must be nil when not set")
	}
	if p.WP != nil {
		t.Fatal("WP must be nil when not set")
	}
}

func TestRawPlayRequiredFields(t *testing.T) {
	// Non-nullable fields must have zero values, not pointers.
	p := RawPlay{
		PlayID:                42,
		Quarter:               2,
		GameClock:             "12:34",
		GameClockTotalSeconds: 754,
		Description:           "test play",
	}
	if p.PlayID != 42 {
		t.Fatalf("PlayID = %d, want 42", p.PlayID)
	}
	if p.GameClockTotalSeconds != 754 {
		t.Fatalf("GameClockTotalSeconds = %d, want 754", p.GameClockTotalSeconds)
	}
}

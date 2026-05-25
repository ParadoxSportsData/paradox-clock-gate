package ingestion

import "testing"

func TestGameHeaderFields(t *testing.T) {
	h := GameHeader{
		GameID:    "2011_01_NO_GB",
		HomeTeam:  "GB",
		AwayTeam:  "NO",
		HomeScore: 42,
		AwayScore: 34,
	}
	if h.GameID != "2011_01_NO_GB" {
		t.Fatalf("GameID = %q, want 2011_01_NO_GB", h.GameID)
	}
	if h.HomeTeam != "GB" {
		t.Fatalf("HomeTeam = %q, want GB", h.HomeTeam)
	}
	if h.AwayTeam != "NO" {
		t.Fatalf("AwayTeam = %q, want NO", h.AwayTeam)
	}
	if h.HomeScore != 42 {
		t.Fatalf("HomeScore = %d, want 42", h.HomeScore)
	}
	if h.AwayScore != 34 {
		t.Fatalf("AwayScore = %d, want 34", h.AwayScore)
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

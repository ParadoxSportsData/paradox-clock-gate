package ingestion

import (
	"testing"
)

func TestParseFileReturnsHeader(t *testing.T) {
	header, _, err := ParseFile("../../testdata/2011_01_NO_GB.json")
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if header.GameID != "2011_01_NO_GB" {
		t.Errorf("GameID = %q, want 2011_01_NO_GB", header.GameID)
	}
	if header.HomeTeam != "GB" {
		t.Errorf("HomeTeam = %q, want GB", header.HomeTeam)
	}
	if header.AwayTeam != "NO" {
		t.Errorf("AwayTeam = %q, want NO", header.AwayTeam)
	}
	if header.HomeScore != 42 {
		t.Errorf("HomeScore = %d, want 42", header.HomeScore)
	}
	if header.AwayScore != 34 {
		t.Errorf("AwayScore = %d, want 34", header.AwayScore)
	}
}

func TestParseFilePlaysNonEmpty(t *testing.T) {
	_, plays, err := ParseFile("../../testdata/2011_01_NO_GB.json")
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if len(plays) == 0 {
		t.Fatal("plays must be non-empty")
	}
}

func TestParseFilePlayFields(t *testing.T) {
	_, plays, err := ParseFile("../../testdata/2011_01_NO_GB.json")
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	// Every play must have a positive PlayID and a non-negative GameClockTotalSeconds.
	for i, p := range plays {
		if p.PlayID <= 0 {
			t.Errorf("plays[%d].PlayID = %d, want > 0", i, p.PlayID)
		}
		if p.GameClockTotalSeconds < 0 {
			t.Errorf("plays[%d].GameClockTotalSeconds = %d, want >= 0", i, p.GameClockTotalSeconds)
		}
	}
}

func TestParseFileNullableFieldsArePointers(t *testing.T) {
	// Nullable fields must remain nil for plays where the JSON value is null.
	// We verify the first play in the file — kickoff plays typically have null down/ydstogo.
	_, plays, err := ParseFile("../../testdata/2011_01_NO_GB.json")
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if len(plays) == 0 {
		t.Fatal("no plays to test")
	}
	// At least one play in the game must have a nil Down (e.g., kickoff).
	foundNilDown := false
	for _, p := range plays {
		if p.Down == nil {
			foundNilDown = true
			break
		}
	}
	if !foundNilDown {
		t.Error("expected at least one play with nil Down (kickoff/PAT), found none")
	}
}

func TestParseFileGameClockTotalSecondsRange(t *testing.T) {
	_, plays, err := ParseFile("../../testdata/2011_01_NO_GB.json")
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	for i, p := range plays {
		if p.GameClockTotalSeconds > 9000 {
			t.Errorf("plays[%d].GameClockTotalSeconds = %d, exceeds MaxTick-1 (9000)", i, p.GameClockTotalSeconds)
		}
	}
}

func TestParseFileNotFound(t *testing.T) {
	_, _, err := ParseFile("../../testdata/does_not_exist.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

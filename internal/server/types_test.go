package server

import (
	"encoding/json"
	"testing"
)

// TestGameStateResponse_NullableFieldsSerializeToNull verifies that pointer fields
// produce JSON null when nil — the contract requires null, not zero values.
func TestGameStateResponse_NullableFieldsSerializeToNull(t *testing.T) {
	gs := GameStateResponse{
		Tick:        900,
		Quarter:     2,
		Down:        nil,
		YardsToGo:   nil,
		YardLine:    nil,
		HomeScore:   7,
		AwayScore:   0,
		Posteam:     nil,
		Defteam:     nil,
		WinProb:     nil,
		PlayType:    "kickoff",
		Description: "kickoff",
		HasState:    true,
	}

	b, err := json.Marshal(gs)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	for _, field := range []string{"down", "yards_to_go", "yard_line", "posteam", "defteam", "win_prob"} {
		v, ok := m[field]
		if !ok {
			t.Errorf("field %q missing from JSON output", field)
			continue
		}
		if v != nil {
			t.Errorf("field %q: want null, got %v", field, v)
		}
	}
}

// TestGameStateResponse_PopulatedFieldsRoundTrip verifies populated nullable fields
// serialize correctly and the JSON keys match the contract (snake_case).
func TestGameStateResponse_PopulatedFieldsRoundTrip(t *testing.T) {
	down := uint8(1)
	ytg := uint8(10)
	yl := uint8(75)
	posteam := "NO"
	defteam := "GB"
	wp := 0.153

	gs := GameStateResponse{
		Tick:        1800,
		Quarter:     3,
		Down:        &down,
		YardsToGo:   &ytg,
		YardLine:    &yl,
		HomeScore:   28,
		AwayScore:   17,
		Posteam:     &posteam,
		Defteam:     &defteam,
		WinProb:     &wp,
		PlayType:    "run",
		Description: "(15:00) M.Ingram up the middle for 1 yard.",
		HasState:    true,
	}

	b, err := json.Marshal(gs)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	checks := map[string]any{
		"tick":        float64(1800),
		"quarter":     float64(3),
		"down":        float64(1),
		"yards_to_go": float64(10),
		"yard_line":   float64(75),
		"home_score":  float64(28),
		"away_score":  float64(17),
		"posteam":     "NO",
		"defteam":     "GB",
		"win_prob":    0.153,
		"play_type":   "run",
		"has_state":   true,
	}

	for key, want := range checks {
		got, ok := m[key]
		if !ok {
			t.Errorf("key %q missing from JSON", key)
			continue
		}
		if got != want {
			t.Errorf("key %q: want %v, got %v", key, want, got)
		}
	}
}

// TestPlaySnapshot_NoDefteamField verifies PlaySnapshot omits defteam
// (only GameStateResponse has defteam per the contract).
func TestPlaySnapshot_NoDefteamField(t *testing.T) {
	ps := PlaySnapshot{
		Tick:     0,
		Quarter:  1,
		PlayType: "kickoff",
	}

	b, err := json.Marshal(ps)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if _, has := m["defteam"]; has {
		t.Error("PlaySnapshot must not have defteam field — it is GameStateResponse-only")
	}
	if _, has := m["has_state"]; has {
		t.Error("PlaySnapshot must not have has_state field — it is GameStateResponse-only")
	}
}

// TestErrorResponse_MaxTickOmittedWhenNil verifies omitempty: max_tick absent
// for non-422 errors, present for 422.
func TestErrorResponse_MaxTickOmittedWhenNil(t *testing.T) {
	notFound := ErrorResponse{Error: "game not found"}
	b, _ := json.Marshal(notFound)
	var m map[string]any
	json.Unmarshal(b, &m)
	if _, has := m["max_tick"]; has {
		t.Error("max_tick must be absent when nil (non-422 response)")
	}

	mt := uint16(3600)
	unprocessable := ErrorResponse{Error: "tick exceeds game length", MaxTick: &mt}
	b, _ = json.Marshal(unprocessable)
	json.Unmarshal(b, &m)
	got, has := m["max_tick"]
	if !has {
		t.Error("max_tick must be present on 422 response")
	}
	if got != float64(3600) {
		t.Errorf("max_tick: want 3600, got %v", got)
	}
}

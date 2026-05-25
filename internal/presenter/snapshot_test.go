package presenter

import (
	"strings"
	"testing"

	"github.com/ParadoxSportsData/paradox-clock-gate/internal/matrix"
)

func makeTestState() (matrix.GameState, matrix.GameMeta, []byte) {
	arena := []byte("(12:34) 12-A.Rodgers pass right to 87-J.Nelson for 12 yards")
	gs := matrix.GameState{
		Elapsed:    150,
		Quarter:    2,
		Down:       3,
		YardsToGo:  8,
		YardLine:   45,
		HomeScore:  14,
		AwayScore:  7,
		PlayType:   matrix.PlayTypePass,
		WinProb:    7240,
		Posteam:    [3]byte{'G', 'B', 0},
		Defteam:    [3]byte{'N', 'O', 0},
		DescOffset: 0,
		DescLen:    uint16(len(arena)),
		HasState:   true,
	}
	meta := matrix.GameMeta{
		GameID:   "2011_01_NO_GB",
		HomeTeam: "GB",
		AwayTeam: "NO",
		MaxTick:  3600,
	}
	return gs, meta, arena
}

func TestRenderTextContainsTeams(t *testing.T) {
	gs, meta, arena := makeTestState()
	out := RenderText(gs, meta, arena)
	if !strings.Contains(out, "GB") {
		t.Error("RenderText output must contain home team GB")
	}
	if !strings.Contains(out, "NO") {
		t.Error("RenderText output must contain away team NO")
	}
}

func TestRenderTextContainsScore(t *testing.T) {
	gs, meta, arena := makeTestState()
	out := RenderText(gs, meta, arena)
	if !strings.Contains(out, "7") {
		t.Error("RenderText output must contain home score 7")
	}
	if !strings.Contains(out, "3") {
		t.Error("RenderText output must contain away score 3")
	}
}

func TestRenderTextContainsQuarter(t *testing.T) {
	gs, meta, arena := makeTestState()
	out := RenderText(gs, meta, arena)
	if !strings.Contains(out, "Q2") {
		t.Error("RenderText output must contain quarter Q2")
	}
}

func TestRenderTextContainsElapsed(t *testing.T) {
	gs, meta, arena := makeTestState()
	out := RenderText(gs, meta, arena)
	if !strings.Contains(out, "150") {
		t.Error("RenderText output must contain elapsed seconds 150")
	}
}

func TestRenderTextContainsWinProb(t *testing.T) {
	gs, meta, arena := makeTestState()
	out := RenderText(gs, meta, arena)
	// WinProb=7240 → 72.4%
	if !strings.Contains(out, "72.4") {
		t.Error("RenderText output must contain win probability 72.4")
	}
}

func TestRenderTextContainsDescription(t *testing.T) {
	gs, meta, arena := makeTestState()
	out := RenderText(gs, meta, arena)
	if !strings.Contains(out, "A.Rodgers") {
		t.Error("RenderText output must contain play description")
	}
}

func TestRenderTextNoStateMessage(t *testing.T) {
	gs := matrix.GameState{HasState: false}
	meta := matrix.GameMeta{HomeTeam: "GB", AwayTeam: "NO"}
	out := RenderText(gs, meta, nil)
	if !strings.Contains(out, "no state") && !strings.Contains(out, "No state") && !strings.Contains(out, "no data") && !strings.Contains(out, "No data") {
		t.Error("RenderText with HasState=false must indicate no state available")
	}
}

func TestRenderJSONContainsFields(t *testing.T) {
	gs, meta, arena := makeTestState()
	out := RenderJSON(gs, meta, arena)
	for _, field := range []string{`"elapsed"`, `"quarter"`, `"home_score"`, `"away_score"`, `"win_prob"`} {
		if !strings.Contains(out, field) {
			t.Errorf("RenderJSON output must contain field %s", field)
		}
	}
}

func TestRenderJSONNullWinProb(t *testing.T) {
	gs, meta, arena := makeTestState()
	gs.WinProb = 65535 // null sentinel
	out := RenderJSON(gs, meta, arena)
	if !strings.Contains(out, "null") {
		t.Error("RenderJSON with WinProb=65535 must render win_prob as null")
	}
}

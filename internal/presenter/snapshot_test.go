package presenter

import (
	"strings"
	"testing"

	"github.com/ParadoxSportsData/paradox-clock-gate/internal/matrix"
)

func makeTestState() (matrix.GameState, matrix.GameMeta, []byte) {
	arena := []byte("(12:34) M.Forte left end to CHI 26 for 5 yards")
	gs := matrix.GameState{
		Elapsed:    150,
		Quarter:    2,
		Down:       3,
		YardsToGo:  8,
		YardLine:   45,
		HomeScore:  7,
		AwayScore:  3,
		PlayType:   matrix.PlayTypeRun,
		WinProb:    5830,
		Posteam:    [3]byte{'C', 'H', 'I'},
		Defteam:    [3]byte{'A', 'T', 'L'},
		DescOffset: 0,
		DescLen:    uint16(len(arena)),
		HasState:   true,
	}
	meta := matrix.GameMeta{
		GameID:   "2011_01_ATL_CHI",
		HomeTeam: "CHI",
		AwayTeam: "ATL",
		MaxTick:  3600,
	}
	return gs, meta, arena
}

func TestRenderTextContainsTeams(t *testing.T) {
	gs, meta, arena := makeTestState()
	out := RenderText(gs, meta, arena)
	if !strings.Contains(out, "CHI") {
		t.Error("RenderText output must contain home team CHI")
	}
	if !strings.Contains(out, "ATL") {
		t.Error("RenderText output must contain away team ATL")
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
	// WinProb=5830 → 58.3%
	if !strings.Contains(out, "58.3") {
		t.Error("RenderText output must contain win probability 58.3")
	}
}

func TestRenderTextContainsDescription(t *testing.T) {
	gs, meta, arena := makeTestState()
	out := RenderText(gs, meta, arena)
	if !strings.Contains(out, "M.Forte") {
		t.Error("RenderText output must contain play description")
	}
}

func TestRenderTextNoStateMessage(t *testing.T) {
	gs := matrix.GameState{HasState: false}
	meta := matrix.GameMeta{HomeTeam: "CHI", AwayTeam: "ATL"}
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

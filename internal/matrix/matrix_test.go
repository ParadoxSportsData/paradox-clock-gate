package matrix

import "testing"

func TestMaxTick(t *testing.T) {
	if MaxTick != 9001 {
		t.Fatalf("MaxTick = %d, want 9001", MaxTick)
	}
}

func TestPlayTypeConstants(t *testing.T) {
	if PlayTypeNone != 0 {
		t.Fatalf("PlayTypeNone = %d, want 0", PlayTypeNone)
	}
	if PlayTypeRun != 1 {
		t.Fatalf("PlayTypeRun = %d, want 1", PlayTypeRun)
	}
	if PlayTypePass != 2 {
		t.Fatalf("PlayTypePass = %d, want 2", PlayTypePass)
	}
	if PlayTypePunt != 3 {
		t.Fatalf("PlayTypePunt = %d, want 3", PlayTypePunt)
	}
	if PlayTypeKickoff != 4 {
		t.Fatalf("PlayTypeKickoff = %d, want 4", PlayTypeKickoff)
	}
	if PlayTypeNoPlay != 5 {
		t.Fatalf("PlayTypeNoPlay = %d, want 5", PlayTypeNoPlay)
	}
	if PlayTypeOther != 6 {
		t.Fatalf("PlayTypeOther = %d, want 6", PlayTypeOther)
	}
}

func TestGameStateZeroValue(t *testing.T) {
	// Zero-value GameState must have HasState = false and WinProb = 0.
	var gs GameState
	if gs.HasState {
		t.Fatal("zero-value GameState must have HasState=false")
	}
	if gs.WinProb != 0 {
		t.Fatalf("zero-value WinProb = %d, want 0", gs.WinProb)
	}
}

func TestGameStateFields(t *testing.T) {
	gs := GameState{
		Elapsed:    900,
		Quarter:    2,
		Down:       3,
		YardsToGo:  8,
		YardLine:   45,
		HomeScore:  7,
		AwayScore:  3,
		PlayType:   PlayTypePass,
		WinProb:    7240,
		Posteam:    [3]byte{'G', 'B', 0},
		Defteam:    [3]byte{'N', 'O', 0},
		DescOffset: 0,
		DescLen:    42,
		HasState:   true,
	}
	if gs.Elapsed != 900 {
		t.Fatalf("Elapsed = %d, want 900", gs.Elapsed)
	}
	if gs.WinProb != 7240 {
		t.Fatalf("WinProb = %d, want 7240", gs.WinProb)
	}
	if gs.Posteam != [3]byte{'G', 'B', 0} {
		t.Fatalf("Posteam = %v, want GB", gs.Posteam)
	}
	if !gs.HasState {
		t.Fatal("HasState must be true when set")
	}
}

func TestStateMatrixSize(t *testing.T) {
	// StateMatrix.States must have exactly MaxTick entries.
	var m StateMatrix
	if len(m.States) != MaxTick {
		t.Fatalf("States length = %d, want %d", len(m.States), MaxTick)
	}
}

func TestGameMetaFields(t *testing.T) {
	meta := GameMeta{
		GameID:   "2011_01_NO_GB",
		HomeTeam: "GB",
		AwayTeam: "NO",
		MaxTick:  1800,
	}
	if meta.GameID != "2011_01_NO_GB" {
		t.Fatalf("GameID = %q, want 2011_01_NO_GB", meta.GameID)
	}
	if meta.MaxTick != 1800 {
		t.Fatalf("MaxTick = %d, want 1800", meta.MaxTick)
	}
}

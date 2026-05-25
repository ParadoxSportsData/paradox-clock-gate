package matrix

import (
	"testing"

	"github.com/ParadoxSportsData/paradox-clock-gate/internal/ingestion"
)

func TestCompileForwardFill(t *testing.T) {
	// Play at tick 10, next play at tick 20 — ticks 11-19 must be forward-filled from tick 10.
	header := ingestion.GameHeader{
		GameID:   "test",
		HomeTeam: "GB",
		AwayTeam: "NO",
	}
	down := 1
	ytg := 10
	yl := 50
	pt := "run"
	posteam := "GB"
	defteam := "NO"
	plays := []ingestion.RawPlay{
		{
			PlayID:                1,
			Quarter:               1,
			GameClock:             "15:00",
			GameClockTotalSeconds: 10,
			Down:                  &down,
			YardsToGo:             &ytg,
			YardLine100:           &yl,
			PlayType:              &pt,
			Description:           "run play",
			Posteam:               &posteam,
			Defteam:               &defteam,
		},
		{
			PlayID:                2,
			Quarter:               1,
			GameClock:             "14:40",
			GameClockTotalSeconds: 20,
			Down:                  &down,
			YardsToGo:             &ytg,
			YardLine100:           &yl,
			PlayType:              &pt,
			Description:           "second run",
			Posteam:               &posteam,
			Defteam:               &defteam,
		},
	}
	m := Compile(plays, header)

	if !m.States[10].HasState {
		t.Fatal("States[10].HasState must be true after compile")
	}
	// Ticks 11-19 are the gap — must be forward-filled from States[10].
	for tick := 11; tick < 20; tick++ {
		if !m.States[tick].HasState {
			t.Fatalf("States[%d] must be forward-filled from States[10]", tick)
		}
		if m.States[tick].Quarter != m.States[10].Quarter {
			t.Errorf("States[%d].Quarter = %d, want %d (forward-filled)", tick, m.States[tick].Quarter, m.States[10].Quarter)
		}
	}
}

func TestCompileTickZeroBeforeFirstPlay(t *testing.T) {
	// Before any play fires, States[0] must have HasState=false.
	header := ingestion.GameHeader{GameID: "test", HomeTeam: "GB", AwayTeam: "NO"}
	down := 1
	ytg := 10
	yl := 50
	pt := "kickoff"
	posteam := "GB"
	defteam := "NO"
	plays := []ingestion.RawPlay{
		{
			PlayID:                1,
			Quarter:               1,
			GameClock:             "15:00",
			GameClockTotalSeconds: 5,
			Down:                  &down,
			YardsToGo:             &ytg,
			YardLine100:           &yl,
			PlayType:              &pt,
			Posteam:               &posteam,
			Defteam:               &defteam,
			Description:           "kickoff",
		},
	}
	m := Compile(plays, header)
	// Ticks 0-4 have no play; they must not have HasState=true.
	for tick := 0; tick < 5; tick++ {
		if m.States[tick].HasState {
			t.Errorf("States[%d].HasState = true, want false (no play before tick 5)", tick)
		}
	}
}

func TestCompileConcurrentPlaysHigherPlayIDWins(t *testing.T) {
	// Two plays at the same tick — higher PlayID must win.
	header := ingestion.GameHeader{GameID: "test", HomeTeam: "GB", AwayTeam: "NO"}
	pt1 := "run"
	pt2 := "pass"
	posteam := "GB"
	defteam := "NO"
	plays := []ingestion.RawPlay{
		{
			PlayID:                1,
			Quarter:               1,
			GameClockTotalSeconds: 100,
			PlayType:              &pt1,
			Description:           "first play",
			Posteam:               &posteam,
			Defteam:               &defteam,
		},
		{
			PlayID:                2,
			Quarter:               1,
			GameClockTotalSeconds: 100,
			PlayType:              &pt2,
			Description:           "winning play",
			Posteam:               &posteam,
			Defteam:               &defteam,
		},
	}
	m := Compile(plays, header)
	if m.States[100].PlayType != PlayTypePass {
		t.Errorf("States[100].PlayType = %v, want PlayTypePass (higher PlayID=2 wins)", m.States[100].PlayType)
	}
}

func TestCompileOTPlays(t *testing.T) {
	// OT plays past second 3600 must compile correctly.
	header := ingestion.GameHeader{GameID: "test", HomeTeam: "GB", AwayTeam: "NO"}
	pt := "run"
	posteam := "GB"
	defteam := "NO"
	plays := []ingestion.RawPlay{
		{
			PlayID:                1,
			Quarter:               5,
			GameClockTotalSeconds: 3700,
			PlayType:              &pt,
			Description:           "OT play",
			Posteam:               &posteam,
			Defteam:               &defteam,
		},
	}
	m := Compile(plays, header)
	if !m.States[3700].HasState {
		t.Fatal("States[3700] must have HasState=true for OT play")
	}
	if m.States[3700].Quarter != 5 {
		t.Errorf("States[3700].Quarter = %d, want 5", m.States[3700].Quarter)
	}
}

func TestCompileScoreAttribution(t *testing.T) {
	// When posteam == HomeTeam: HomeScore = PosteamScore, AwayScore = DefteamScore.
	header := ingestion.GameHeader{GameID: "test", HomeTeam: "GB", AwayTeam: "NO"}
	pt := "pass"
	posteam := "GB" // posteam == HomeTeam
	defteam := "NO"
	homeScore := 7
	awayScore := 3
	plays := []ingestion.RawPlay{
		{
			PlayID:                1,
			Quarter:               2,
			GameClockTotalSeconds: 500,
			PlayType:              &pt,
			Description:           "td pass",
			Posteam:               &posteam,
			Defteam:               &defteam,
			PosteamScore:          &homeScore,
			DefteamScore:          &awayScore,
		},
	}
	m := Compile(plays, header)
	if m.States[500].HomeScore != 7 {
		t.Errorf("HomeScore = %d, want 7 (posteam==home)", m.States[500].HomeScore)
	}
	if m.States[500].AwayScore != 3 {
		t.Errorf("AwayScore = %d, want 3 (defteam==away)", m.States[500].AwayScore)
	}
}

func TestCompileScoreAttributionAwaySide(t *testing.T) {
	// When posteam == AwayTeam: AwayScore = PosteamScore, HomeScore = DefteamScore.
	header := ingestion.GameHeader{GameID: "test", HomeTeam: "GB", AwayTeam: "NO"}
	pt := "run"
	posteam := "NO" // posteam == AwayTeam
	defteam := "GB"
	awayScore := 10
	homeScore := 7
	plays := []ingestion.RawPlay{
		{
			PlayID:                1,
			Quarter:               3,
			GameClockTotalSeconds: 2000,
			PlayType:              &pt,
			Description:           "NO run",
			Posteam:               &posteam,
			Defteam:               &defteam,
			PosteamScore:          &awayScore,
			DefteamScore:          &homeScore,
		},
	}
	m := Compile(plays, header)
	if m.States[2000].AwayScore != 10 {
		t.Errorf("AwayScore = %d, want 10 (posteam==away)", m.States[2000].AwayScore)
	}
	if m.States[2000].HomeScore != 7 {
		t.Errorf("HomeScore = %d, want 7 (defteam==home)", m.States[2000].HomeScore)
	}
}

func TestCompileMaxTick(t *testing.T) {
	header := ingestion.GameHeader{GameID: "test", HomeTeam: "GB", AwayTeam: "NO"}
	pt := "run"
	posteam := "GB"
	defteam := "NO"
	plays := []ingestion.RawPlay{
		{PlayID: 1, Quarter: 1, GameClockTotalSeconds: 100, PlayType: &pt, Description: "a", Posteam: &posteam, Defteam: &defteam},
		{PlayID: 2, Quarter: 4, GameClockTotalSeconds: 3500, PlayType: &pt, Description: "b", Posteam: &posteam, Defteam: &defteam},
	}
	m := Compile(plays, header)
	if m.Meta.MaxTick != 3500 {
		t.Errorf("Meta.MaxTick = %d, want 3500", m.Meta.MaxTick)
	}
}

func BenchmarkQuery(b *testing.B) {
	header := ingestion.GameHeader{GameID: "bench", HomeTeam: "GB", AwayTeam: "NO"}
	pt := "run"
	posteam := "GB"
	defteam := "NO"
	plays := make([]ingestion.RawPlay, 150)
	for i := range plays {
		tick := i * 24
		plays[i] = ingestion.RawPlay{
			PlayID:                i + 1,
			Quarter:               (i / 40) + 1,
			GameClockTotalSeconds: tick,
			PlayType:              &pt,
			Description:           "bench play description text",
			Posteam:               &posteam,
			Defteam:               &defteam,
		}
	}
	m := Compile(plays, header)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = m.States[1800]
	}
}

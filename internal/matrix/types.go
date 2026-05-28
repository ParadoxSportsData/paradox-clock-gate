// Package matrix defines the core data types and state-matrix structure for
// O(1) temporal game-state queries.
package matrix

const MaxTick = 9001 // covers regulation (3600s) + up to ~3 OT periods

const (
	WinProbNull  uint16 = 65535 // sentinel: no win probability data available
	WinProbScale uint16 = 10000 // multiply float wp by WinProbScale before storing as uint16
)

type PlayType uint8

const (
	PlayTypeNone    PlayType = 0
	PlayTypeRun     PlayType = 1
	PlayTypePass    PlayType = 2
	PlayTypePunt    PlayType = 3
	PlayTypeKickoff PlayType = 4
	PlayTypeNoPlay  PlayType = 5
	PlayTypeOther   PlayType = 6
)

// GameState is ~28 bytes with no pointer fields — GC-free at query time.
type GameState struct {
	Elapsed    uint16
	Quarter    uint8
	Down       uint8
	YardsToGo  uint8
	YardLine   uint8   // yards to opponent endzone (0-100)
	HomeScore  uint8
	AwayScore  uint8
	PlayType   PlayType
	WinProb    uint16  // wp * WinProbScale (e.g., 0.583 → 5830); WinProbNull = no data
	Posteam    [3]byte // null-padded e.g. "CHI"
	Defteam    [3]byte
	DescOffset uint32
	DescLen    uint16
	HasState   bool
}

type GameMeta struct {
	GameID   string
	HomeTeam string
	AwayTeam string
	MaxTick  uint16
	GameDate string // ISO-8601 from JSON header e.g. "2011-09-08T00:00:00"
	Week     int
	Season   int
}

type StateMatrix struct {
	States    [MaxTick]GameState
	Arena     []byte   // all descriptions concatenated (one alloc at init)
	Meta      GameMeta
	PlayTicks []uint16 // ticks with real plays, sorted ascending; excludes forward-filled ticks
}

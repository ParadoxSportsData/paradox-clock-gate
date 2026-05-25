# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

`paradox-clock-gate` is a Go CLI tool that reads NFL play-by-play JSON files (270 games, 2011 season) and answers "what was the exact game state at elapsed second T?" with O(1) lookup and guaranteed future-state containment (no data from after tick T can leak into the result).

**Assessment context:** Take-home submission for a Software Engineer/Forward Deployed Engineer role at Prelude (Origin team). Contact: claudia@preludesecurity.com. The assessment asks for a real tool built with AI as a core workflow tool; the write-up should be honest about AI mistakes and pushback.

**Full project context:** This tool is the proof-of-concept for the temporal engine that will eventually live at `internal/engine/` in a full monolith rewrite of `paradox-platform` (sibling at `../paradox-platform`). Do not start that rewrite — this CLI is the deliverable. Full rewrite architecture decision: monolith, not microservices (team-scaling problem doesn't exist yet).

**Scope:** 4-8 hour implementation. Working binary + `WRITEUP.md`.

**GitHub:** `ParadoxSportsData/paradox-clock-gate`

---

## Initial Setup (if not yet done)

```bash
git init
go mod init github.com/ParadoxSportsData/paradox-clock-gate
```

Then create the full package directory structure and stub files with package declarations before implementing.

---

## Build & Test Commands

```bash
# Build binary
go build ./cmd/clock-gate/

# Run all tests + benchmarks
go test -bench=. -benchmem ./...

# Run a single package's tests
go test ./internal/matrix/...

# Critical benchmark — query path must show 0 allocs/op
go test -bench=BenchmarkQuery -benchmem ./internal/matrix/

# Vet
go vet ./...

# Smoke tests (run from repo root after building)
./clock-gate --tick 0 testdata/2011_01_NO_GB.json      # kickoff
./clock-gate --tick 1800 testdata/2011_01_NO_GB.json   # halftime
./clock-gate --tick 3599 testdata/2011_01_NO_GB.json   # late game
./clock-gate --tick 900 --format json testdata/2011_01_NO_GB.json
./clock-gate --tick 999999 testdata/2011_01_NO_GB.json  # expect error
./clock-gate --list testdata/
```

---

## Package Structure

```
paradox-clock-gate/
├── cmd/clock-gate/main.go        # CLI entry: flag parsing, open file, pipeline
├── internal/
│   ├── ingestion/
│   │   ├── schema.go             # RawPlay struct (maps 1:1 to JSON fields)
│   │   └── parser.go             # Stream-parse JSON with json.Decoder (token mode)
│   ├── matrix/
│   │   ├── types.go              # GameState, GameMeta, StateMatrix, MaxTick, PlayType enum
│   │   └── compiler.go           # Compile([]RawPlay, GameMeta) StateMatrix
│   ├── gate/
│   │   └── gate.go               # Validate(tick, maxTick) — bounds + containment
│   └── presenter/
│       └── snapshot.go           # RenderText(GameState, GameMeta, arena) + RenderJSON()
├── testdata/                     # 3 curated GB Packers wins; ships in repo, no external dep
├── go.mod
├── README.md
└── WRITEUP.md
```

---

## Architecture

### Core Design

Irregular play events (~150 per game) are compiled into a **flat pre-allocated array indexed directly by elapsed second**. Query time is `matrix.States[tick]` — one array dereference, zero heap pressure.

**Forward-fill algorithm:** After writing plays into the array by index, a single pass iterates 0→MaxTick copying `States[t-1]` into any empty `States[t]`. This fills gaps and guarantees O(1) lookup. Future state cannot bleed backward because the fill only copies earlier ticks forward, never later ticks backward. This is the temporal isolation guarantee.

### Key Types

```go
// internal/matrix/types.go

const MaxTick = 9001  // covers regulation (3600s) + up to ~3 OT periods

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

// ~28 bytes. No pointers = no GC pressure at runtime.
type GameState struct {
    Elapsed    uint16   // matches array index
    Quarter    uint8    // 1-6; 0 = not yet started
    Down       uint8    // 0 = no down (kickoff etc), 1-4
    YardsToGo  uint8
    YardLine   uint8    // yards to opponent endzone (0-100)
    HomeScore  uint8
    AwayScore  uint8
    PlayType   PlayType
    WinProb    uint16   // wp * 10000 (e.g., 0.583 → 5830); 65535 = null
    Posteam    [3]byte  // null-padded e.g. "GB\x00"
    Defteam    [3]byte
    DescOffset uint32   // byte offset into StateMatrix.Arena
    DescLen    uint16
    HasState   bool
}

type GameMeta struct {
    GameID   string
    HomeTeam string
    AwayTeam string
    MaxTick  uint16
}

type StateMatrix struct {
    States  [MaxTick]GameState  // ~252 KB; indexed directly by elapsed second
    Arena   []byte              // all descriptions concatenated (one alloc at init)
    Meta    GameMeta
}
```

```go
// internal/ingestion/schema.go — nullable JSON fields use pointer types

// GameHeader holds the top-level wrapper fields parsed before the plays array.
// home_team, away_team, home_score, away_score come from the JSON header —
// NOT from filename and NOT from play data (posteam/defteam change every drive).
type GameHeader struct {
    GameID    string `json:"game_id"`
    HomeTeam  string `json:"home_team"`
    AwayTeam  string `json:"away_team"`
    HomeScore int    `json:"home_score"`
    AwayScore int    `json:"away_score"`
}

type RawPlay struct {
    PlayID                int      `json:"play_id"`
    Quarter               int      `json:"quarter"`
    GameClock             string   `json:"game_clock"`
    GameClockTotalSeconds int      `json:"game_clock_total_seconds"`
    Down                  *int     `json:"down"`
    YardsToGo             *int     `json:"ydstogo"`
    YardLine100           *int     `json:"yardline_100"`
    PlayType              *string  `json:"play_type"`
    YardsGained           *int     `json:"yards_gained"`
    Description           string   `json:"description"`
    Posteam               *string  `json:"posteam"`
    Defteam               *string  `json:"defteam"`
    PosteamScore          *int     `json:"posteam_score"`
    DefteamScore          *int     `json:"defteam_score"`
    WP                    *float64 `json:"wp"`
    EPA                   *float64 `json:"epa"`
}
```

### Compiler Algorithm (`internal/matrix/compiler.go`)

`Compile(plays []RawPlay, header GameHeader) StateMatrix`:

1. Sort plays by `(GameClockTotalSeconds ASC, PlayID ASC)` — for concurrent events at the same tick, higher PlayID wins
2. Allocate `StateMatrix`; pre-allocate Arena with `make([]byte, 0, totalDescLen)` (one allocation)
3. For each play: convert to `GameState`, append description to Arena, write into `States[play.GameClockTotalSeconds]`, track `maxTick`. **Score attribution:** compare `RawPlay.Posteam` against `GameHeader.HomeTeam` to correctly map possession scores to home/away fields: if posteam == HomeTeam → HomeScore=PosteamScore, AwayScore=DefteamScore; else → AwayScore=PosteamScore, HomeScore=DefteamScore.
4. Forward-fill: `for t = 1 to maxTick: if !States[t].HasState { States[t] = States[t-1]; States[t].Elapsed = t }`
5. Set `matrix.Meta.MaxTick = maxTick`

### Ingestion (`internal/ingestion/parser.go`)

**Actual JSON structure** (verified against testdata/2011_01_NO_GB.json):
```json
{
  "game_id": "2011_01_NO_GB",
  "home_team": "GB",
  "away_team": "NO",
  "home_score": 42,
  "away_score": 34,
  "plays": [ ... ]
}
```

Each file is a wrapper object, NOT a flat array. `ParseFile` returns both `GameHeader` and `[]RawPlay`.

**Parser algorithm** — token-mode `json.Decoder` over the wrapper object:
1. Decode `{` (outer object start)
2. Loop: read next key token
   - If key matches a `GameHeader` field (`game_id`, `home_team`, `away_team`, `home_score`, `away_score`): `decoder.Decode(&headerField)`
   - If key == `"plays"`: decode `[`, loop `decoder.Decode(&play)` until `]`
   - Otherwise: call `decoder.Token()` repeatedly to skip the value (handles nested objects/arrays)
3. Decode `}` (outer object end)

**Home/away teams:** Read from `GameHeader.HomeTeam` / `GameHeader.AwayTeam` — they are top-level JSON fields. Do NOT parse from filename. Do NOT use play-level `posteam`/`defteam` (those indicate possession, which changes every drive).

**`game_clock_total_seconds`** is pre-computed in the data — use it directly. Recompute formula is a fallback only if a play is missing it.

**Validated against real data:**
- Max `game_clock_total_seconds` across all 270 games: **4500** (OT games confirmed)
- `MaxTick = 9001` is safe
- All plays have `game_clock_total_seconds` populated (confirmed in sample)

### CLI

Standard library `flag` only — no Cobra (overkill for 3 flags, zero deps). Three flags: `--tick int`, `--format text|json`, `--list`.

---

## Data Source

**Dev/demo:** `testdata/` ships 3 curated GB Packers wins directly in the repo — no external dependency required.

| File | Game | Result |
|------|------|--------|
| `testdata/2011_01_NO_GB.json` | NO @ GB, Week 1 2011 | GB 42 – NO 34 |
| `testdata/2011_09_GB_SD.json` | GB @ SD, Week 9 2011 | GB 45 – SD 38 |
| `testdata/2011_14_OAK_GB.json` | OAK @ GB, Week 14 2011 | GB 46 – OAK 16 |

**Full dataset:** `../paradox-platform/data/raw/*.json` — 270 JSON files, one per NFL game, 2011 season. Naming: `{season}_{week}_{away}_{home}.json`. ~85-95 KB each, ~120-150 plays per game. Not required for dev or demos.

**`game_clock_total_seconds`** is the primary index field — pre-computed elapsed seconds since kickoff:
- Regulation: 0 (kickoff) → 3600 (end of Q4)
- OT Q5: 3600–4500; OT Q6: 4500–5400
- Recompute formula if needed: `3600 - game_seconds_remaining` for regulation; `3600 + ((qtr-4-1)*900) + (900 - game_seconds_remaining)` for OT

**Nullable fields:** `down`, `ydstogo`, `yardline_100`, `posteam`, `defteam`, `posteam_score`, `defteam_score`, `wp`, `epa` can all be `null` (kickoffs, special teams, game markers) — handle with pointer types.

---

## Design Constraints (Do Not Change)

| Decision | Choice | Rationale |
|---|---|---|
| Language | Go | Assessment requirement, performance story |
| CLI framework | `flag` stdlib only — no Cobra | Zero deps; right tool for 3 flags |
| Lookup strategy | Pre-allocated flat array, O(1) | Not binary search — trades 252 KB for zero conditionals at query time |
| Runtime allocations after init | Zero | No strings or pointers in hot path |
| `WinProb` storage | `uint16` (wp × 10000) | Not `float64` — keeps struct GC-free |
| Team abbreviations | `[3]byte` | Not `string` — eliminates pointer in struct |
| Descriptions | Arena allocator (`[]byte`) with offset+length | Not per-play strings — one allocation at init |
| Output scope | Game state only — not cumulative stats | Focused; achievable in 4-8h |
| `StateMatrix` size | `[MaxTick]GameState` ≈ 252 KB | Acceptable tradeoff for O(1) |

---

## CLI Interface

```
Usage: clock-gate --tick <seconds> [--format text|json] <game-file>
       clock-gate --list <directory>

Examples:
  clock-gate --tick 1800 testdata/2011_01_NO_GB.json
  clock-gate --tick 900 --format json testdata/2011_01_NO_GB.json
  clock-gate --list testdata/
```

## Text Output Format

```
┌────────────────────────────────────────────────────────────┐
│  NO @ GB   │  Q3  │  Elapsed: 1800s (30:00)               │
├────────────────────────────────────────────────────────────┤
│  Score:  GB 28  –  NO 17                                   │
│  Ball:   NO possession │  1st & 10  at NO 80               │
│  Win Prob: NO 15.3%                                        │
├────────────────────────────────────────────────────────────┤
│  (15:00) 28-M.Ingram up the middle to NO 21 for 1 yard… │
└────────────────────────────────────────────────────────────┘
```

Display win prob relative to posteam: "GB Win %" when posteam == home team; "NO Win %" when posteam == away. Elapsed display: `Xs (M:SS)`.

---

## Implementation Phases

### Phase 1: Ingestion (~1.5h)
Files: `internal/ingestion/schema.go`, `internal/ingestion/parser.go`
- `ParseFile(path string) (GameHeader, []RawPlay, error)` — token-mode `json.Decoder`
- Parse home/away from JSON header (`GameHeader.HomeTeam` / `GameHeader.AwayTeam`) — NOT from filename
- Unit test: parse `testdata/2011_01_NO_GB.json`, assert correct play count, spot-check field values

**Definition of done:** `go test ./internal/ingestion/...` passes

### Phase 2: State Matrix Compiler (~2h)
Files: `internal/matrix/types.go`, `internal/matrix/compiler.go`
- Implement `Compile()` per algorithm above
- Unit tests:
  - Forward-fill: `States[11]` equals `States[10]` when no play at second 11
  - Concurrent events: two plays at same tick → higher `play_id` wins
  - OT: plays past second 3600 compile and query correctly
- Benchmark: `BenchmarkQuery` — single array lookup must show `0 allocs/op`

**Definition of done:** `go test -bench=BenchmarkQuery -benchmem ./internal/matrix/` shows `0 allocs/op`

### Phase 3: Gate + CLI + Presenter (~2h)
Files: `internal/gate/gate.go`, `internal/presenter/snapshot.go`, `cmd/clock-gate/main.go`
- Gate: `Validate(tick int, maxTick uint16) error` — clear error if tick > maxTick or tick < 0
- Presenter: `RenderText(gs GameState, meta GameMeta, arena []byte) string` and `RenderJSON(...) string`
- Main: wire flags → open file → parse → compile → gate → present → print

**Definition of done:** Binary runs; `--tick 0` returns kickoff; `--tick 999999` returns bounded error message

### Phase 4: Polish + Write-up (~1.5h)
- `--list` flag: scan directory, print game IDs and file names
- `README.md`: install instructions, usage examples, architecture overview
- `WRITEUP.md`: fill in from the narrative scaffold already in that file
- `go vet ./...` clean, `go build ./...` clean

---

## Session Workflow

**Claude's responsibilities — Aaron does not update memory files:**

Start of session:
1. Read `~/.claude-work/projects/-Users-athatcher-Documents-at-proj-paradox-clock-gate/memory/implementation_log.md` — know exactly where we left off and what's next
2. Read `ai_moments.md` — maintain continuity on the AI collaboration story

During implementation (Claude writes immediately, not at session end):
- Invoke `capture-decision` skill whenever: a design decision is made, scope changes, an AI suggestion is accepted or overridden, a guardrail is added, a known flaw is identified, or the README may need updating
- `capture-decision` 6-item checklist: (1) memory files, (2) WRITEUP.md, (3) interview prep, (4) plan, (5) plugin repo content, (6) continuous improvement log
- `go test -bench` output is auto-captured by the PostToolUse benchmark hook

End of session:
- Write session entry to `implementation_log.md` (what was done, blockers, next starting point)
- The Stop hook verifies memory files were updated — if they weren't, that is a Claude failure

Phase gates (mandatory — no exceptions):
- Invoke `superpowers:test-driven-development` at the START of each implementation phase
- Invoke `superpowers:verification-before-completion` before declaring any phase DONE
- Invoke `superpowers:finishing-a-development-branch` before final submission

---

## Related Project

`../paradox-platform` — the Python/TypeScript PoC this tool is derived from. Full dataset lives at `../paradox-platform/data/raw/` — not required for dev or demos (testdata/ ships in this repo). Do not modify that project.

## Plan File

Full master plan: `~/.claude-work/plans/now-do-you-understand-rustling-sky.md`
Implementation phases: `~/.claude-work/plans/jiggly-wishing-beacon.md`

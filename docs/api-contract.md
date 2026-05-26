# clock-gate serve mode — API Contract

> **Canonical contract for Phase 2A (Go backend) and Phase 2B (React frontend).**
> Both sides build against this document. Neither side merges changes to this file
> without creating a corresponding amendment ticket for the other side.

---

## Base URL

```
http://localhost:8080
```

Configurable via `--port` flag. CORS allowed for all origins (local dev).

---

## Endpoints

### `GET /games`

List all available games loaded from the data directory.

**Response:** `200 OK` — `application/json` — `GameSummary[]`

```json
[
  {
    "game_id": "2011_01_NO_GB",
    "home_team": "GB",
    "away_team": "NO",
    "home_score": 42,
    "away_score": 34,
    "duration": 3600
  }
]
```

---

### `GET /games/{id}/state?tick=N`

Get exact game state at elapsed second N. O(1) StateMatrix lookup.

**Query params:** `tick` (required, integer ≥ 0)

**Response: `200 OK`** — `application/json` — `GameStateResponse`

```json
{
  "tick": 1800,
  "quarter": 3,
  "down": 1,
  "yards_to_go": 10,
  "yard_line": 75,
  "home_score": 28,
  "away_score": 17,
  "posteam": "NO",
  "defteam": "GB",
  "win_prob": 0.153,
  "play_type": "run",
  "description": "(15:00) M.Ingram up the middle for 1 yard.",
  "has_state": true
}
```

**Errors:**
- `400 Bad Request` — tick < 0 or missing
- `404 Not Found` — game ID not found in data directory
- `422 Unprocessable Entity` — tick > game maxTick (with `max_tick` in error body)

---

### `GET /games/{id}/timeline`

Full game timeline — every populated tick as a sorted array. Load once, scrub locally.

**Response: `200 OK`** — `application/json` — `GameTimelineResponse`

```json
{
  "game_id": "2011_01_NO_GB",
  "home_team": "GB",
  "away_team": "NO",
  "max_tick": 3600,
  "plays": [
    {
      "tick": 0,
      "quarter": 1,
      "down": null,
      "yards_to_go": null,
      "yard_line": null,
      "home_score": 0,
      "away_score": 0,
      "posteam": "GB",
      "win_prob": 0.5,
      "play_type": "kickoff",
      "description": "(15:00) GB kicks off..."
    }
  ]
}
```

`plays` contains only ticks where `has_state = true` (irregular play events, not the forward-filled gaps). Sorted by `tick` ASC. Frontend performs local forward-fill for slider position → nearest play lookup.

**Errors:**
- `404 Not Found` — game ID not found

---

### Error Format

All error responses:

```json
{ "error": "human-readable message", "max_tick": 3600 }
```

`max_tick` is only present in 422 responses.

---

## Go Response Types

```go
// internal/server/types.go

type GameSummary struct {
    GameID    string `json:"game_id"`
    HomeTeam  string `json:"home_team"`
    AwayTeam  string `json:"away_team"`
    HomeScore uint8  `json:"home_score"`
    AwayScore uint8  `json:"away_score"`
    Duration  uint16 `json:"duration"` // maxTick in seconds
}

type GameStateResponse struct {
    Tick        uint16   `json:"tick"`
    Quarter     uint8    `json:"quarter"`
    Down        *uint8   `json:"down"`        // null for kickoff/special
    YardsToGo   *uint8   `json:"yards_to_go"` // null for kickoff/special
    YardLine    *uint8   `json:"yard_line"`    // null for kickoff/special
    HomeScore   uint8    `json:"home_score"`
    AwayScore   uint8    `json:"away_score"`
    Posteam     *string  `json:"posteam"`      // null if no possession
    Defteam     *string  `json:"defteam"`      // null if no possession
    WinProb     *float64 `json:"win_prob"`     // null if unknown; WinProb uint16/10000
    PlayType    string   `json:"play_type"`    // "run"|"pass"|"punt"|"kickoff"|"no_play"|"other"|""
    Description string   `json:"description"`
    HasState    bool     `json:"has_state"`
}

type PlaySnapshot struct {
    Tick        uint16   `json:"tick"`
    Quarter     uint8    `json:"quarter"`
    Down        *uint8   `json:"down"`
    YardsToGo   *uint8   `json:"yards_to_go"`
    YardLine    *uint8   `json:"yard_line"`
    HomeScore   uint8    `json:"home_score"`
    AwayScore   uint8    `json:"away_score"`
    Posteam     *string  `json:"posteam"`
    WinProb     *float64 `json:"win_prob"`
    PlayType    string   `json:"play_type"`
    Description string   `json:"description"`
}

type GameTimelineResponse struct {
    GameID   string         `json:"game_id"`
    HomeTeam string         `json:"home_team"`
    AwayTeam string         `json:"away_team"`
    MaxTick  uint16         `json:"max_tick"`
    Plays    []PlaySnapshot `json:"plays"`
}

type ErrorResponse struct {
    Error   string `json:"error"`
    MaxTick *uint16 `json:"max_tick,omitempty"` // 422 only
}
```

---

## TypeScript Interfaces

```typescript
// src/api/types.ts

export interface GameSummary {
  game_id: string;       // "2011_01_NO_GB"
  home_team: string;     // "GB"
  away_team: string;     // "NO"
  home_score: number;
  away_score: number;
  duration: number;      // maxTick in seconds
}

export interface GameStateResponse {
  tick: number;
  quarter: number;           // 1-6
  down: number | null;       // 1-4, null for kickoff/special
  yards_to_go: number | null;
  yard_line: number | null;  // 0-100, yards to opponent endzone
  home_score: number;
  away_score: number;
  posteam: string | null;    // "GB", "NO", etc.
  defteam: string | null;
  win_prob: number | null;   // 0.0-1.0, null if unknown
  play_type: string;         // "run"|"pass"|"punt"|"kickoff"|"no_play"|"other"|""
  description: string;
  has_state: boolean;        // false = before any play at this tick
}

export interface PlaySnapshot {
  tick: number;
  quarter: number;
  down: number | null;
  yards_to_go: number | null;
  yard_line: number | null;
  home_score: number;
  away_score: number;
  posteam: string | null;
  win_prob: number | null;
  play_type: string;
  description: string;
}

export interface GameTimelineResponse {
  game_id: string;
  home_team: string;
  away_team: string;
  max_tick: number;
  plays: PlaySnapshot[]; // only has_state=true ticks, sorted by tick ASC
}

export interface ApiError {
  error: string;
  max_tick?: number; // present on 422 only
}
```

---

## Contract Amendment Protocol

If either side (2A Go or 2B TypeScript) needs to add, remove, or modify a field:

1. Create a Jira ticket: `P2A-N: Contract amendment — [change description]`
2. Mark the requesting side's ticket **blocked by** P2A-N
3. P2A-N must implement the change in **both** Go types and TypeScript interfaces
4. Both sides unblock once P2A-N merges

---

## Notes

- **Nullable fields:** `down`, `yards_to_go`, `yard_line`, `posteam`, `defteam`, `win_prob` are JSON `null` for kickoffs, special teams, and game markers. Frontend must handle null gracefully.
- **`win_prob`:** Stored as `uint16 × 10000` in `GameState.WinProb`; converted to `float64` in API response. `65535` sentinel → `null` in JSON.
- **`play_type`:** Maps `PlayType` enum to string. `PlayTypeNone` (0) → `""`.
- **`has_state`:** Only `false` for ticks before the first play (pre-kickoff). `/timeline` omits these ticks entirely; `/state` returns them with `has_state: false`.
- **Path sanitization:** Game ID is sanitized with `filepath.Clean` and prefix-asserted before any file open. Prevents path traversal.

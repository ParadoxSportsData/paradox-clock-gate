# Deep Code Review — Go (paradox-clock-gate)
Date: 2026-05-27
Scope: FULL
Files reviewed: 21 (11 production, 10 test)
Packages reviewed: 6 (cmd/clock-gate, internal/gate, internal/ingestion, internal/matrix, internal/presenter, internal/server)

---

## Section 1 — Build / Test Status

```
go build ./...   → PASS (no output, no errors)
go vet ./...     → PASS (no output, no errors)
go test ./...    → PASS (56 tests, all packages green)

BenchmarkQuery-14    1000000000    0.2402 ns/op    0 B/op    0 allocs/op
```

The 0 allocs/op on BenchmarkQuery confirms the core architectural promise holds.

---

## Section 2 — Executive Summary

- **Magic Number proliferation (WinProb encoding):** The sentinel `65535` and scale factor `10000` appear as bare literals in 4 files across 3 packages (compiler, presenter, server/handler). Any change to the encoding strategy requires surgical edits in all four places without a compile-time guarantee that all sites were found.
- **Duplicate CORS middleware:** An identical CORS handler is defined twice — `corsMiddleware` in `cmd/clock-gate/serve.go` and `corsMiddlewareHandler` in `internal/server/handler.go`. In the `serve` mode the outer `corsMiddleware` wraps the entire mux, while each route inside `NewServeMux` is already wrapped with `corsMiddlewareHandler`, causing every response in serve mode to have CORS headers set twice.
- **Missing package-level doc comments:** All six packages lack `// Package X ...` documentation. Every exported type in `server/types.go` and several in `matrix/types.go` and `ingestion/schema.go` also lack godoc. This increases time-to-change for any new contributor.
- **Test coverage gaps in server and presenter:** The HTTP `handleState` endpoint is not tested for unknown game IDs (404 path), missing tick parameter (400 path), or `HasState=false` semantics. `TestRenderTextContainsScore` has a misleading assertion that passes accidentally.
- **HTTP server missing read/write timeouts:** `http.ListenAndServe` is called with no `http.Server` struct, leaving the server with no `ReadTimeout`, `WriteTimeout`, or `IdleTimeout`, which exposes it to Slowloris-style resource exhaustion.

---

## Section 3 — Findings

### [HIGH] Duplicate CORS middleware with double-application in serve mode
File: `cmd/clock-gate/serve.go:41` and `internal/server/handler.go:22`
Category: Design
Smell: Duplicate Code (2 sites — identical function bodies)
Finding: `corsMiddleware` in serve.go and `corsMiddlewareHandler` in handler.go have byte-for-byte identical bodies; in serve mode the mux returned by `NewServeMux` (which already wraps each route with `corsMiddlewareHandler`) is then wrapped again with `corsMiddleware`, causing `Access-Control-Allow-Origin: *` to be set twice on every serve-mode response.
Impact: Duplicate response headers are benign in browsers today but violate HTTP/1.1 specification for headers that must not be duplicated; inconsistent behavior between the two CORS implementations if either is ever changed without updating the other.
Fix: Delete `corsMiddleware` from `serve.go` entirely. Move the single canonical implementation into `internal/server/handler.go` as an exported function `server.CORSMiddleware`. Call it from both `NewServeMux` (removing the per-route wrapping) and `runServe`. This collapses 2 duplicate implementations into 1 and eliminates the double-wrap.

---

### [HIGH] Magic number `65535` (WinProb null sentinel) across 4 files
File: `internal/matrix/compiler.go:83`, `internal/presenter/snapshot.go:35`, `internal/presenter/snapshot.go:105`, `internal/server/handler.go:218`
Category: Hard-Coded Values
Smell: Magic Number (5 sites across 4 files including test)
Finding: The value `65535` encoding "WinProb is null" is embedded as a bare integer literal at 4 production sites and 1 test site with no named constant; the comment "null sentinel" appears at most sites but is not enforced at compile time.
Impact: If the sentinel value is changed (e.g., to `math.MaxUint16` for clarity or to accommodate a different null convention), all 5 sites must be found manually; a missed site silently introduces a bug where null win probability is rendered as a valid percentage.
Fix: Add `const WinProbNull uint16 = 65535` to `internal/matrix/types.go` adjacent to the `GameState` definition. Replace all 5 literal occurrences with `matrix.WinProbNull` (or `WinProbNull` within the same package). This is 1 site for the definition and 0 logic changes.

---

### [HIGH] Magic number `10000` (WinProb scale factor) across 3 files
File: `internal/matrix/compiler.go:81`, `internal/presenter/snapshot.go:106`, `internal/server/handler.go:221`
Category: Hard-Coded Values
Smell: Magic Number (3 sites across 3 files)
Finding: The fixed-point scale factor `10000` (converting float64 win probability to uint16 and back) appears as a bare literal at 3 production sites with no named constant; the types.go comment documents the encoding but does not enforce it.
Impact: If precision is increased (e.g., to `100000` for 5 decimal places), the type must also change; a site missed in the conversion causes silent precision loss or overflow.
Fix: Add `const WinProbScale uint32 = 10000` to `internal/matrix/types.go`. Replace all 3 literal occurrences with `matrix.WinProbScale`. The `uint32` type prevents overflow when used in intermediate arithmetic (`uint16(wp * WinProbScale)`).

---

### [MEDIUM] HTTP server missing read/write/idle timeouts
File: `cmd/clock-gate/serve.go:36`
Category: Design (Security / Resource Management)
Smell: N/A
Finding: `http.ListenAndServe(addr, handler)` uses Go's default server with no `ReadTimeout`, `WriteTimeout`, or `IdleTimeout`; an attacker or misbehaving client can hold connections open indefinitely.
Impact: In a demo or local-network context the risk is low, but the server will exhaust file descriptors under a Slowloris attack or when many slow clients connect. Adding timeouts is a one-line change with no behavior change for well-behaved clients.
Fix: Replace `http.ListenAndServe` with:
```go
srv := &http.Server{
    Addr:         addr,
    Handler:      handler,
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  120 * time.Second,
}
log.Fatal(srv.ListenAndServe())
```

---

### [MEDIUM] Missing package-level doc comments on all packages
File: `internal/ingestion/schema.go:1`, `internal/matrix/types.go:1`, `internal/gate/gate.go:1`, `internal/presenter/snapshot.go:1`, `internal/server/handler.go:1`, `cmd/clock-gate/main.go:1`
Category: Maintainability
Smell: N/A
Finding: No package in the codebase has a `// Package X ...` comment above the `package` declaration; `go doc` produces no package-level summary for any package.
Impact: Per K10, each missing package doc adds 15 minutes to time-to-change for a new contributor; 6 missing docs = 90 minutes of avoidable ramp-up time. For a portfolio project this is the first thing a technical reviewer reads.
Fix: Add one line per package immediately before `package X`:
- `// Package ingestion parses NFL game JSON files into GameHeader and RawPlay slices.`
- `// Package matrix compiles RawPlay slices into an O(1)-queryable StateMatrix.`
- `// Package gate validates temporal query bounds against a compiled game's MaxTick.`
- `// Package presenter renders GameState as human-readable text or JSON.`
- `// Package server exposes the StateMatrix as a REST API via HTTP.`
- No doc needed for `package main`.

---

### [MEDIUM] Missing godoc on 8 exported symbols across 3 files
File: `internal/ingestion/schema.go:17` (RawPlay), `internal/matrix/types.go:5` (PlayType), `internal/matrix/types.go:35` (GameMeta), `internal/matrix/types.go:45` (StateMatrix), `internal/server/types.go:3` (GameSummary), `internal/server/types.go:15` (GameStateResponse), `internal/server/types.go:31` (PlaySnapshot), `internal/server/types.go:45` (GameTimelineResponse), `internal/server/types.go:53` (ErrorResponse)
Category: Maintainability
Smell: N/A
Finding: Nine exported types are missing godoc comments; `RawPlay` has none; all five exported types in `server/types.go` have none; `PlayType`, `GameMeta`, and `StateMatrix` in `matrix/types.go` have none. `GameState` has a comment but it is an inline struct note, not a godoc comment.
Impact: +10 minutes per missing symbol for a new contributor. At 9 symbols that is 90 minutes of additional ramp-up time before someone can confidently modify the API contract types.
Fix: Add `// TypeName ...` godoc immediately before each type declaration. For `server/types.go` types, the JSON field names are the API contract — document that: `// GameStateResponse is the JSON body for GET /games/{id}/state?tick=N.`

---

### [MEDIUM] `TestRenderTextContainsScore` has misleading assertions that pass accidentally
File: `internal/presenter/snapshot_test.go:48`
Category: Testing
Smell: N/A
Finding: The test sets `HomeScore=14, AwayScore=7` but its comment says "want home score 7" and "want away score 3". The check `strings.Contains(out, "3")` passes because elapsed `150s (2:30)` contains the character `"3"` in `"30"` — not because any score is 3. The check for `"7"` happens to match AwayScore=7, so it is accidentally correct.
Impact: A regression that changes AwayScore from 7 to another value (e.g., 14) would still pass because the test does not check that "7" is in a score context. The test provides false confidence.
Fix: Replace with explicit score string checks. Use `strings.Contains(out, "GB 14")` for home score and `strings.Contains(out, "NO 7")` for away score, matching the exact format produced by `RenderText`.

---

### [MEDIUM] `gate_test.go`: 5 separate test functions that should be a table-driven test
File: `internal/gate/gate_test.go:5`
Category: Testing
Smell: Duplicate Code (5 nearly identical test functions)
Finding: `TestValidateNegativeTick`, `TestValidateZeroTick`, `TestValidateTickAtMaxTick`, `TestValidateTickExceedsMaxTick`, and `TestValidateTickFarExceedsMaxTick` each call `Validate` with different inputs and check error presence; this is the canonical table-driven pattern.
Impact: Per K5, "No TestXxx function with more than 3 sequential branches that could be table rows." Adding a new boundary case (e.g., testing maxTick=0) requires a new function rather than a new row. The 5 functions add clutter without adding clarity.
Fix: Consolidate into one `TestValidate` with `tests []struct{tick int; maxTick uint16; wantErr bool}`, covering all 5 cases plus any new ones.

---

### [MEDIUM] `forward-fill` leaves pre-game ticks in an unhandled state for the HTTP API
File: `internal/server/handler.go:88`
Category: Design
Smell: N/A
Finding: A request to `/games/{id}/state?tick=0` when the first play is at tick 5 returns HTTP 200 with `has_state: false` and all zero values (quarter=0, score=0-0). This is undocumented API behavior — the response is not an error but is not meaningful game state either. No test covers this path.
Impact: API consumers that don't check `has_state` will silently render corrupted scoreboards for any tick before the first recorded play.
Fix (two-part): (1) Add a `TestServeStateBeforeFirstPlay` test that requests tick=0 and asserts `has_state=false`. (2) Document the behavior in the godoc for `handleState`: "Ticks before the first recorded play return HTTP 200 with `has_state: false`."

---

### [LOW] `padRight` uses byte length instead of rune count — potential display misalignment
File: `internal/presenter/snapshot.go:124`
Category: Design
Smell: N/A
Finding: `padRight(s, n)` uses `len(s)` (byte count) to pad to width `n`. If play descriptions contain multi-byte UTF-8 characters (em-dashes, ellipses), the visual box width will be shorter than intended because `len` counts bytes, not display columns.
Impact: The box-drawing frame will appear misaligned in terminals for any play description containing non-ASCII characters. NFL play descriptions are ASCII-dominant so this rarely manifests, but the ellipsis appended at truncation (`"…"` at line 69) is itself a 3-byte UTF-8 character, meaning the truncated line is visually 1 character wider than expected.
Fix: Replace `len(s)` with `len([]rune(s))` in `padRight`, and replace the `"…"` literal at line 69 with the ASCII `"~"` or change `maxDesc` accounting to subtract 3 bytes for the ellipsis instead of 1.

---

### [LOW] `ingestion_test.go`: `TestGameHeaderFields` is a pure struct construction test with no value
File: `internal/ingestion/ingestion_test.go:5`
Category: Testing
Smell: N/A
Finding: `TestGameHeaderFields` builds a `GameHeader` struct literal and asserts the fields equal what was just assigned. This test cannot fail unless the struct is removed; it tests Go's assignment semantics, not any project logic.
Impact: Adds 29 lines to the test file that maintain themselves indefinitely without protecting any real invariant.
Fix: Delete the test. The invariant it is trying to protect (field names match JSON tags) is better covered by `TestParseFileReturnsHeader`, which exercises the same fields through the real JSON parser.

---

### [LOW] `ingestion_test.go`: `TestRawPlayNullableFields` tests Go zero values, not project logic
File: `internal/ingestion/ingestion_test.go:30`
Category: Testing
Smell: N/A
Finding: `TestRawPlayNullableFields` asserts that pointer fields of a zero-value `RawPlay{}` are nil. This is guaranteed by Go's zero-value semantics; no project code is exercised.
Impact: Same as above — false confidence without real coverage.
Fix: Delete the test. The `TestParseFileNullableFieldsArePointers` parser test already covers the meaningful case (JSON `null` round-trips to nil pointer) with real data.

---

### [LOW] `matrix_test.go`: `TestMaxTick`, `TestPlayTypeConstants`, `TestGameStateZeroValue` test language semantics
File: `internal/matrix/matrix_test.go:5`
Category: Testing
Smell: N/A
Finding: `TestMaxTick` asserts `MaxTick == 9001` (a constant cannot differ from itself), `TestPlayTypeConstants` asserts `PlayTypeRun == 1` etc. (iota constants cannot change without a compile-time break), and `TestGameStateZeroValue` asserts zero-value fields are zero.
Impact: These tests add noise and maintenance burden without protecting any project behavior. A real regression in this area (e.g., reordering iota values) would be caught by the compiler, not by these tests.
Fix: Delete these three test functions. Replace with a single compile-time check if iota ordering matters: `var _ = [1]struct{}{}[PlayTypeNone]` style assertions or a comment referencing the wire format.

---

### [LOW] `playTypeStr` returns `""` for `PlayTypeNone` but this is undocumented
File: `internal/server/handler.go:226`
Category: Maintainability
Smell: N/A
Finding: The default case in `playTypeStr` returns `""` for `PlayTypeNone` (value 0), which is the zero value of the enum and appears whenever `p.PlayType == nil` during compilation. API consumers receive `play_type: ""` with no documentation of what this means.
Impact: Low — the behavior is deterministic — but consumers may not know that empty string means "no play type recorded" rather than "unknown play type."
Fix: Add a case `case matrix.PlayTypeNone: return ""` with a comment: `// PlayTypeNone: play_type field was null in source data.` This documents the contract without changing behavior.

---

### [LOW] `compiler.go` silently drops plays with `tick >= MaxTick`
File: `internal/matrix/compiler.go:42`
Category: Error Handling
Smell: N/A
Finding: Plays with `tick >= MaxTick` are silently skipped with a `continue` and no log message; if real game data contains plays beyond second 9000 (e.g., a 7th OT period), they vanish from the compiled matrix without any observable signal.
Impact: The tool would return incorrect state for extreme OT games rather than returning an error. Given that `MaxTick=9001` is designed to cover 6 OT periods and no 2011 game had 7, this is low probability but high severity if it occurs silently.
Fix: Add `fmt.Fprintf(os.Stderr, "warning: play %d at tick %d exceeds MaxTick %d — skipped\n", p.PlayID, tick, MaxTick)` before `continue`. This makes the silent drop audible without changing behavior.

---

## Section 4 — Package-Level Verdicts

| Package | Files | Critical | High | Medium | Low | Time-to-Change (min) |
|---------|-------|----------|------|--------|-----|----------------------|
| cmd/clock-gate | 4 | 0 | 1 | 1 | 0 | 40 (missing pkg doc +15, CORS dup +20, HTTP timeout medium) |
| internal/gate | 2 | 0 | 0 | 1 | 0 | 25 (missing pkg doc +15, table-driven +10) |
| internal/ingestion | 4 | 0 | 0 | 1 | 2 | 40 (missing pkg doc +15, RawPlay godoc +10, 2 trivial tests) |
| internal/matrix | 4 | 0 | 2 | 1 | 1 | 65 (missing pkg doc +15, 3 missing type docs +30, magic numbers +15, trivial tests) |
| internal/presenter | 2 | 0 | 0 | 1 | 1 | 40 (missing pkg doc +15, test assertion bug +15, padRight +10) |
| internal/server | 5 | 0 | 1 | 3 | 1 | **75** (missing pkg doc +15, 5 type godocs +50, CORS dup +5, undocumented has_state +5) |

**server package exceeds 60 minutes: HIGH — maintainability barrier.** The server package has the densest concentration of missing documentation across all five exported types and no package-level doc.

---

## Section 5 — File Inventory

| File | Status | Finding count |
|------|--------|---------------|
| cmd/clock-gate/main.go | REVIEWED | 0 issues |
| cmd/clock-gate/serve.go | REVIEWED | 2 issues (CORS duplicate, HTTP timeouts) |
| internal/gate/gate.go | REVIEWED | 0 issues |
| internal/ingestion/schema.go | REVIEWED | 2 issues (RawPlay godoc, package doc) |
| internal/ingestion/parser.go | REVIEWED | 0 issues |
| internal/matrix/types.go | REVIEWED | 3 issues (PlayType/GameMeta/StateMatrix godoc, package doc, WinProb const candidates) |
| internal/matrix/compiler.go | REVIEWED | 2 issues (magic numbers 65535/10000, silent tick drop) |
| internal/presenter/snapshot.go | REVIEWED | 2 issues (magic number 65535/10000, padRight bytes vs runes) |
| internal/server/handler.go | REVIEWED | 3 issues (CORS duplicate, magic numbers, playTypeStr undocumented) |
| internal/server/loader.go | REVIEWED | 0 issues |
| internal/server/types.go | REVIEWED | 1 issue (all 5 exported types missing godoc) |
| cmd/clock-gate/integration_test.go | REVIEWED | 0 issues |
| cmd/clock-gate/serve_test.go | REVIEWED | 0 issues |
| internal/gate/gate_test.go | REVIEWED | 1 issue (table-driven refactor) |
| internal/ingestion/ingestion_test.go | REVIEWED | 2 issues (2 trivial tests to delete) |
| internal/ingestion/parser_test.go | REVIEWED | 0 issues |
| internal/matrix/compiler_test.go | REVIEWED | 0 issues |
| internal/matrix/matrix_test.go | REVIEWED | 1 issue (3 trivial tests to delete) |
| internal/presenter/snapshot_test.go | REVIEWED | 1 issue (misleading score assertion) |
| internal/server/serve_test.go | REVIEWED | 1 issue (missing coverage for unknown game ID + HasState=false) |
| internal/server/types_test.go | REVIEWED | 0 issues |

---

## What's Well Done

**Architecture is sound and internally consistent.** The flat array with forward-fill is the right tradeoff for O(1) queries — the 252 KB fixed allocation is negligible, and `0 allocs/op` at query time proves the design promise. The forward-fill algorithm is correct: it propagates state forward in a single left-to-right pass because each filled tick gets `HasState=true`, allowing the next iteration to continue propagation across arbitrarily long gaps.

**Concurrency in `GameCache` is correct and well-documented.** The double-checked locking pattern in `loader.go` is properly implemented: read lock → fast path → release → parse outside lock → write lock → double-check → insert. The comment at lines 18–19 explains the invariant clearly. This is production-quality concurrent code.

**Path traversal guard is present and tested.** `loader.go` uses `filepath.Clean` + prefix check with a trailing separator (to prevent `/datadir` matching `/datadir-other`), and `TestServePathTraversal` verifies it with a percent-encoded `..` sequence. This is correct.

**Error wrapping is consistent throughout.** Every error at package boundaries is wrapped with `fmt.Errorf("context: %w", err)`, enabling callers to use `errors.Is`. No error is silently discarded except the two documented `//nolint:errcheck` sites.

**Score attribution logic is correct and well-tested.** The `posteam == HomeTeam` comparison with both home-possession and away-possession test cases (`TestCompileScoreAttribution` and `TestCompileScoreAttributionAwaySide`) correctly pins final score to `GameHeader` home/away, not to possession team labels that change every drive.

**Integration test coverage is strong for the happy path.** All five smoke tests cover the CLI binary end-to-end, and the seven server tests exercise the real mux + cache + parser pipeline against actual testdata fixtures. The path traversal test and CORS header test are specifically good security/contract regression tests.

**`GameCache.ListGames` gracefully handles malformed files.** Silently skipping unparseable files during listing (rather than failing the entire list) is the right choice for a data directory that may contain non-game JSON files.

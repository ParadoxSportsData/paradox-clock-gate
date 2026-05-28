# clock-gate: Assessment Write-up

*Submission for Prelude / Origin team — Software Engineer / Forward Deployed Engineer role*

---

## Problem & Motivation

The use case: a user is watching a recorded NFL game on a streaming service — Super Bowl XLV from 2011, say, on YouTube. They want a second screen that knows exactly where they are in the broadcast and shows live-style game state for that moment: score, down-and-distance, win probability, last play. Not the final box score. Not season totals. The exact state at the elapsed second they're watching. Spoiler-free.

That experience depends on one primitive: given elapsed second T, return the exact game state that was true at T — nothing from after it, nothing missing from before it.

The `paradox-platform` Python PoC proved the concept: index NFL play-by-play data by `game_clock_total_seconds` and answer "what was the game state at moment T?" But it had real technical debt — linear scans, N+1 queries, hardcoded localhost URLs, and no temporal isolation guarantee (a bug could let future-play data bleed into historical queries). The core innovation deserved a clean implementation.

`clock-gate` extracts and perfects that one mechanism: the forward-fill compiler makes the temporal isolation guarantee mathematical rather than conventional. There is no code path by which a play at tick T+1 can affect the result returned for tick T.

`ParadoxSportsData` — the GitHub org this lives under — is a deliberate platform foundation, not a one-product account. Clock-gate is the first artifact: the temporal query engine that every future product in the org depends on. The rewind league, the momentum engine, the injury-aware simulation — all of those require accurate point-in-time game state as a primitive. You build that once, correctly, and everything else composes on top of it.

---

## How AI Was Used

Claude Code (this session) was the primary execution vehicle — architecture design, code generation, debugging, iteration. Gemini (earlier session) provided the initial `clock-gate` concept framing and naming.

This was a deliberate choice: use AI for everything I'd normally do by muscle memory (scaffolding, boilerplate, type definitions) so my cognitive budget stays on the decisions that actually matter (data representation, allocation strategy, temporal correctness proof).

Beyond code generation, AI was used to build reusable workflow infrastructure: a set of project-level skills (READ-DO checklists for the agent) covering artifact creation, requirements capture, and decision logging. These skills are submitted alongside the code — they are artifacts of the process, not scaffolding discarded after use. The skill system demonstrates that AI collaboration at a team level requires workflow design, not just prompting.

The centerpiece of that infrastructure is a skill called `capture-decision`. Whenever a significant decision was made during this build — a trade-off chosen, a guardrail added, an AI suggestion overridden, a known flaw identified — the agent invoked this skill before moving on. It walks a six-item checklist: update memory files, update the WRITEUP, update interview prep, update the plan, evaluate for the reusable workflow tool, log to the improvement backlog. The agent maintained its own audit trail in real time. By the time implementation was complete, the interview prep was already current, the known flaws were already documented with improvement paths, and the improvement roadmap was already written — not reconstructed from memory after the fact.

---

## How I Set Up the Agent to Work Safely

Using AI as a core workflow tool doesn't mean turning it loose. Before writing a line of implementation code, I built an access and constraint layer around the agent — the same way I'd configure a service account before granting it access to a production system.

**Principle of least privilege.** The project `.claude/settings.json` contains an explicit `toolApprovals` allowlist: the exact Bash commands the agent is permitted to run without prompting. Nothing outside that list runs silently. This is not a default — it's a deliberate access control decision, documented and auditable.

**Scope isolation.** A PreToolUse hook warns if the agent attempts to write or modify files outside the `paradox-clock-gate/` project root. The sibling repository (`../paradox-platform`) is live data I care about. The agent cannot touch it by accident.

**Architectural constraint enforcement.** Two of the core technical claims — zero external dependencies and zero query-time allocations — are enforced at the infrastructure level, not just by convention. A PostToolUse hook on `go.mod`-touching commands warns if an external dependency appears. The benchmark capture hook flags immediately if `allocs/op` is non-zero. The constraints are machine-verified, not trust-based.

**Test enforcement.** A second PreToolUse hook on `git commit` runs `go test ./...` and blocks the commit if any test is failing. Tests-passing is a gate enforced at the infrastructure level, not a convention. The TDD discipline — writing tests first — is handled by a mandatory `superpowers:test-driven-development` skill invoked at the start of each implementation phase; the hook enforces the outcome, the skill enforces the process.

**Secrets hygiene.** A PreToolUse hook on `git commit` scans staged files for credential patterns before anything is committed. clock-gate has no real secrets, but the habit costs nothing and the serve mode (Phase 2A) introduces configuration surface. The check runs regardless.

**Path traversal protection (serve mode).** Any endpoint that serves files from a directory is a path traversal surface. The gate layer sanitizes the game ID parameter — `filepath.Clean`, absolute path resolution, prefix assertion against the data directory — before any file is opened. Documented as a design constraint, not an afterthought.

The through-line: every constraint that matters to the correctness or safety of this tool is enforced somewhere other than my memory or convention. If I forgot the zero-dep rule tomorrow, the hook would catch it. That's the standard I hold for agent-assisted development.

---

## Where AI Got It Wrong / Where I Pushed Back

**1. Cobra framework suggestion**
The initial Gemini session suggested Cobra for CLI handling. Overruled: this tool has three flags. Cobra adds a dependency, import overhead, and abstraction that helps nothing here. `flag` stdlib is the right tool, easier for a reviewer to audit, and makes the binary's dependency tree trivially inspectable.

**2. String fields in the hot struct**
First struct design had `string` fields for `Posteam` and `Description` directly in `GameState`. Rejected. Strings in Go are (`pointer`, `length`) — that's a GC-visible pointer in the hot struct. Replaced with `[3]byte` for team abbreviations and an arena allocator (`[]byte` with `DescOffset`+`DescLen`) for descriptions. Result: `GameState` contains no pointers, no GC pressure, no per-query allocations.

**3. Binary search at query time**
Gemini's prompt suggested binary search over sorted play timestamps. Overruled in favor of forward-fill + O(1) array lookup. Trade-off: 252 KB of pre-allocated memory (`[9001]GameState`) versus true O(1) with zero runtime conditionals. At game-file scale (~90 KB input), 252 KB is trivially cheap. The O(1) claim is cleaner and provably correct.

**4. `float64` for `WinProb`**
AI defaulted to `float64` in the struct. Rejected: replaced with `uint16` storing `wp × 10000` (e.g., 0.583 → 5830). Gives 0.01% precision — more than sufficient — while keeping the struct free of floating-point and GC-visible memory. Reconvert to float only at display time.

**5. Home/away teams from filename**
AI initially described parsing home and away team abbreviations from the filename (e.g., `2011_01_NO_GB.json` → away=NO, home=GB). Overruled after inspecting the actual data: `home_team` and `away_team` are explicit top-level JSON fields. Parsing from the filename introduces fragile assumptions about naming conventions; parsing from JSON is authoritative and eliminates them entirely.

**6. Flat array parser assumption**
AI described the parser as opening `[`, looping `decoder.Decode(&play)`, closing `]` — assuming the file was a flat JSON array of plays. Overruled after reading an actual game file: the format is a wrapper object with `game_id`, `home_team`, `away_team`, `home_score`, `away_score`, and a nested `"plays"` array. The token-mode parser must handle the outer object structure first, key by key, before entering the plays array.

**7. Forward-fill test wrote a structurally incorrect test**
During TDD Red phase, the agent wrote a forward-fill test with only one play at tick 10 and asserted that `States[11]` would be forward-filled. But with a single play, `maxTick = 10` and the fill loop stops at 10 — it never reaches tick 11. The test would have passed incorrectly if `States[11]` happened to have HasState=false (which it did, by zero-value). The error was caught immediately by the mandatory RED-verify step: run the test, confirm it fails for the right reason. It failed for the wrong reason — the test was bad, not the code. Fixed by adding a second play at tick 20 to create a real gap to fill. This is exactly why TDD's "watch it fail correctly" step is not optional.

**8. AI hedged on NFL timing, user corrected**
During Phase 2B UI debugging, the timeline scrubber appeared frozen — the kickoff play was still showing at what the UI labelled "Q1 1:43" elapsed. The AI suggested this *might* be valid NFL timing — perhaps the game clock runs differently than assumed. The user corrected immediately: "A kickoff does not take 1 minute and 17 seconds. That's patently ridiculous." The AI accepted the correction without argument. Root cause, confirmed: `handleTimeline` was emitting all ~3,501 forward-filled ticks rather than the ~154 real plays. Every tick from 0 to maxTick had `HasState=true` after forward-fill, so the endpoint serialized thousands of entries sharing the same description bytes. The frontend binary search worked correctly — it found the exact tick — but descriptions didn't change for hundreds of consecutive ticks. This is the failure mode of leaking an internal implementation detail (the forward-fill optimization) through an API contract. Fixed: `PlayTicks []uint16` added to `StateMatrix`, populated at compile time, `handleTimeline` now ranges over only real play ticks. Lesson: on domain-specific facts, the domain expert's correction is authoritative. The AI's uncertainty was wrong to surface as a legitimate alternative.

**9. AI recommended a per-tick LRU cache for stats; user corrected all games are recorded**
During paradox-stats design, the agent proposed an LRU cache keyed on `(game_id, tick)` — cache the computed stats result for a given tick, evict least-recently-used entries as new ticks are requested. The implicit mental model was a live game where future state is uncertain: cache what you've computed so far, because tomorrow's plays don't exist yet. The user corrected immediately: our use case is always recorded games. The complete play-by-play is in the data file at access time — the future isn't uncertain, it's just not needed yet. Caching a per-tick aggregate is the wrong abstraction: it caches an intermediate computation keyed on a query parameter, not the underlying data. The right design is to pre-compute stats at each play boundary at first access (~150 boundaries per game), store those snapshots, and binary-search at query time. One build per game, not one cache entry per tick per user. The agent accepted the correction immediately and proposed the StatsMatrix architecture: a sorted slice of ~150 play-boundary snapshots, binary search per query, one LRU at the game level shared across all users. The resulting design is simpler and faster — a per-tick LRU would grow unboundedly with query breadth; StatsMatrix has a fixed 150 entries regardless of how many ticks are queried.

**10. AI wrote compiler.go without normalizing WinProb to home-team perspective**
The initial `compiler.go` stored `wp` from nflfastR directly as `WinProb` — without recognizing that the raw field is possession-team win probability, not home-team. The bug was invisible in the CLI output (the label said `{posteam} Win Prob` which was technically accurate) but surfaced immediately when the WinProbChart in the React UI showed values that jumped with possession changes: the chart was labeled "GB Win %" but the plotted line spiked whenever the away team took possession. The visual was obviously wrong. Fix: `compiler.go` normalizes at compile time — `1 - wp` when `posteam != homeTeam` — so `WinProb` in `GameState` always means home-team win probability for all consumers. The UI was updated to display both teams' current percentages in the chart header; the away percentage is `1 - homeWP`, consistent regardless of possession.

**11. Filter implementation: dots correct, slider lookup wrong**
During PDX-56 (timeline filter visual feedback), the agent correctly added play tick dots to the timeline bar — `filteredTicks` derived from `filteredPlays` memo, rendered only at filtered-play positions, updating instantly when filter changes. But `handleSliderChange` still called `findNearestPlay(query.data.plays, newTick)` — all plays — so dragging the slider between two run-play dots could land on a pass play's tick and display it even with "Run" filter active. The visual layer (dots) was correct; the interaction layer was wrong. User caught it: "Run selection still shows a pass play. Scoring shows a non-scoring play." Root cause required reading the actual play type distribution from the live API, tracing from `handleSliderChange` through `findNearestPlay`, confirming the API returns only five play types — and identifying that dragging between filtered dots silently bypassed the filter. Fix: when `activeFilter !== 'all'`, both `handleSliderChange` and `commitTimeInput` resolve from `filteredPlays` instead of `allPlays`. The lesson: "half the problem solved" bugs are invisible until someone actually uses the feature — building a dot layer that looks right doesn't mean the interaction model matches the user's expectation.

---

## Honest Rough Edges

- **`--tick` takes raw elapsed seconds, not wall-clock time.** A user thinking "Q2 15:00 remaining" has to calculate 900 themselves (15 minutes × 60). A `--clock "Q2 15:00"` flag would be the right UX fix.
- **OT game detection relies on play data having `quarter >= 5`** — there's no explicit header in the JSON indicating OT. Works correctly but is implicit.
- **`--list` is minimal** — just filenames and no play counts or final scores.
- **No `--range` flag** to diff state between two ticks (the obvious next feature for drive summaries).

---

## Benchmark

The central technical claim: zero allocations on the query path. Verified with `go test -bench=BenchmarkQuery -benchmem ./internal/matrix/` on Apple M4 Pro:

```
goos: darwin
goarch: arm64
pkg: github.com/ParadoxSportsData/paradox-clock-gate/internal/matrix
cpu: Apple M4 Pro
BenchmarkQuery-14    	1000000000	         0.2337 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/ParadoxSportsData/paradox-clock-gate/internal/matrix	0.689s
```

0 allocs/op confirms the design: `GameState` contains no pointer fields (`[3]byte` for team abbreviations, `uint16`/`uint8` for all numeric fields). The GC has nothing to collect at query time. 0.23 ns/op is one array dereference — the expected cost of a cache-warm L1 access.

---

## What I'd Add With More Time

**Near-term CLI polish:**
- **`--clock "Q2 15:00"` input format:** Parse quarter + clock string instead of requiring raw elapsed seconds. The current UX asks users to do arithmetic that the tool should do for them.
- **`--range T1 T2`:** Return what changed between two ticks — possession, score, drive progression. The natural bridge from point query to interval query, and the foundation for drive summaries.
- **`--filter scoring`:** Seek to scoring plays only. Useful for highlight generation and the "catch me up" use case.

**Production infrastructure:**
- **`clock-gate serve` — HTTP mode (Phase 2A):** Expose the compiled StateMatrix over `net/http` stdlib (same zero-dep rationale as `flag` over Cobra). One StateMatrix per game in server memory, shared across all concurrent users — O(1) query path with no GC pressure becomes a horizontally scalable read layer.
- **React + TypeScript + Vite UI (Phase 2B):** Timeline scrubber, score display, field position, win probability chart over the serve mode. Different constraints from the backend: React's component model and TypeScript's compile-time API safety are the right tools for a UI; zero-dep thinking doesn't apply to a layer where maintainability and extensibility matter more than allocation counts.
- **Plugin repo — AI workflow system:** The skills, hooks, memory structure, and settings patterns built during clock-gate setup packaged as a Claude Code marketplace plugin. Any project installs it once and gets the full collaboration infrastructure: memory system, skill-builder, Stop hook, toolApprovals allowlist, scope isolation guard. The plugin repo makes the system reusable without reinventing what the marketplace already provides.
- **StatsMatrix incremental accumulator:** The current paradox-stats `StatsMatrix.build()` calls pandas aggregation functions once per play boundary (~150 calls per game), each scan covering all plays from tick 0 to that boundary — O(n × plays) total. An incremental accumulator would walk plays in tick order, maintain running per-player and per-team totals, and snapshot at each boundary — O(plays) total build time. Deferred from v1 because 1–3 seconds per game is acceptable when the cost amortizes across all users of that game. High priority as paradox-stats expands to cover more seasons.

**The harder engineering problem — live in-progress games:**
The replay case (completed game) is actually the easier temporal isolation problem: the future doesn't exist yet in the data. The more interesting case is a game that has started but not yet ended — the user starts watching from the beginning while the game is still being played. Now the future *does* exist in the data (the ongoing game), and the gate must actively hide it. The StateMatrix is no longer statically compiled at load time; it's growing as plays arrive. `MaxTick` is no longer known at startup — it's the user's current video position, not the game's live edge. This requires streaming compilation (append plays as they arrive without recompiling the full matrix), a dynamic gate (compare tick against video position, not game completion), and a spoiler firewall (the serve endpoint must know where each user's video is, not just where the game is). That's a fundamentally different concurrency model from the static compiled case.

**The longer arc — rewind league, momentum, injury:**
The temporal engine is the foundation for a simulation product that's the real vision: a "Rewind Sports League" where historical performance data through year X drives a randomized but realistic league simulation. To build that, three additional primitives need to exist on top of what clock-gate provides:

- *Momentum quantification:* The `WinProb` trajectory across a drive sequence is already a momentum signal embedded in the data — WP moving 0.45 → 0.38 → 0.31 across three plays is measurable pressure. Quantifying momentum requires extracting series-level patterns from the play sequence, not just point-in-time state. The forward-fill structure clock-gate produces is the right shape for this: every elapsed second is populated, so rolling windows over the sequence are trivial.

- *Injury-aware simulation:* NFL performance follows ebbs and flows — a player who was at peak performance for six weeks and then took a significant hit performs differently in the weeks that follow. A realistic simulation has to model not just historical averages but the state a player was in at a given point in their career. `nflfastR` carries injury report data that clock-gate's ingestion layer doesn't currently surface. Extending the ingestion to capture player-level health signals is the first step.

- *Coaching persona simulation:* The domain expert personas in the `paradox-platform` knowledge base — offensive coordinator, defensive coordinator, QB coach — encode football decision logic. In a rewind league, these become AI agent roles: an OC agent with "script the first 15 plays, exploit the matchup" as its prime directive making play-calling decisions constrained by its team's personnel health and momentum state. The temporal engine provides the historical truth; the persona agents provide the decision layer that generates realistic simulation from it.

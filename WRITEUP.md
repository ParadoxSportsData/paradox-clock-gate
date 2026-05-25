# clock-gate: Assessment Write-up

*Submission for Prelude / Origin team — Software Engineer / Forward Deployed Engineer role*

---

## Problem & Motivation

The `paradox-platform` Python PoC proved the concept: index NFL play-by-play data by a linear timestamp and answer "what was the game state at moment T?" in sub-second time. But it had real technical debt — linear scans, N+1 queries, hardcoded localhost URLs, and no temporal isolation guarantee (a bug could let future-play data bleed into historical queries). The core innovation (linear timestamp for video sync) deserved a clean implementation.

`clock-gate` extracts and perfects that one mechanism: given a game file and an elapsed-second offset, return the exact game state that was true at that tick — nothing from after it, nothing missing from before it. The forward-fill compiler makes this a mathematical guarantee, not just a convention.

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
AI defaulted to `float64` in the struct. Replaced with `uint16` storing `wp × 10000` (e.g., 0.583 → 5830). Gives 0.01% precision — more than sufficient — while keeping the struct free of floating-point and GC-visible memory. Reconvert to float only at display time.

**5. Forward-fill test wrote a structurally incorrect test**
During TDD Red phase, the agent wrote a forward-fill test with only one play at tick 10 and asserted that `States[11]` would be forward-filled. But with a single play, `maxTick = 10` and the fill loop stops at 10 — it never reaches tick 11. The test would have passed incorrectly if `States[11]` happened to have HasState=false (which it did, by zero-value). The error was caught immediately by the mandatory RED-verify step: run the test, confirm it fails for the right reason. It failed for the wrong reason — the test was bad, not the code. Fixed by adding a second play at tick 20 to create a real gap to fill. This is exactly why TDD's "watch it fail correctly" step is not optional.

---

## Honest Rough Edges

- **`--tick` takes raw elapsed seconds, not wall-clock time.** A user thinking "Q2 15:00 remaining" has to calculate 900 themselves (15 minutes × 60). A `--clock "Q2 15:00"` flag would be the right UX fix.
- **OT game detection relies on play data having `quarter >= 5`** — there's no explicit header in the JSON indicating OT. Works correctly but is implicit.
- **`--list` is minimal** — just filenames and no play counts or final scores.
- **No `--range` flag** to diff state between two ticks (the obvious next feature for drive summaries).

---

## What I'd Add With More Time

- **`--range T1 T2`:** Show what changed between two ticks — possession changes, score changes, drive progression. This is the bridge from "point query" to "interval query" and the most natural next primitive.
- **`--filter scoring`:** Jump to only scoring plays. Useful for highlight generation.
- **`--clock "Q2 15:00"` input format:** Parse quarter + clock string instead of requiring raw elapsed seconds.
- **Plugin repo — AI workflow system (post-MVP):** The skills, hooks, memory structure, and settings patterns built during clock-gate setup packaged as a Claude Code marketplace plugin. Any project installs it once and gets the full collaboration infrastructure: memory system, skill-builder, Stop hook, toolApprovals allowlist, scope isolation guard. Clock-gate is the proof-of-concept; the plugin repo makes the system reusable without reinventing what the marketplace already provides.
- **`clock-gate serve` — HTTP mode (Phase 2A):** Expose the compiled StateMatrix over HTTP via `net/http` stdlib (same zero-dep rationale). One StateMatrix per game loaded into memory, shared across all concurrent users — the O(1) query path with no GC pressure becomes a horizontally scalable read layer. The direct bridge to the full `paradox-platform` rewrite.
- **React + TypeScript UI (Phase 2B):** Timeline scrubber, score display, field position, win probability chart over the serve mode. Different constraints from the backend — React's component model and TypeScript's compile-time API safety are the right tools for a UI, even if zero-dep thinking doesn't apply here.

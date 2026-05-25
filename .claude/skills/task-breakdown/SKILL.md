---
name: task-breakdown
description: Creates the Confluence task tracker page with phase milestones, dependency map, DoD commands, and parallel opportunity callouts. TRIGGER when: user says "create the task tracker", "create the milestone doc", "build the task breakdown", "set up the project tracker", or the pre-code setup sequence reaches the tracker step.
---

## Overview

Task trackers written as flat to-do lists omit the two things evaluators look for: explicit parallel work opportunities (evidence of planning) and exact definitions of done per phase (evidence of rigor). This skill enforces both before creating the Confluence page.

**Iron Law:** Every phase row must have an exact shell command or observable state as its definition of done. "Implementation complete" is not a definition of done.

---

## Announce at Start

"Creating task tracker for paradox-clock-gate. Drafting all phase milestones with exact DoD commands, mapping bidirectional dependencies and parallel opportunities, then creating the Confluence page. PAGE URL will be written for downstream reference."

---

## Inputs Required

Before proceeding, confirm: have all Jira tickets been created (via `create-jira-ticket` or `task-breakdown` itself)? If no tickets exist yet, create-jira-ticket must be run first — this skill links to existing ticket IDs, it does not create them.

If tickets exist: proceed. Write the ticket IDs you have.

---

## Steps

### Step 1 — Invoke `artifact-context` (Killer)

Invoke the `artifact-context` skill. Pass: phase = "Pre-code setup", component = "setup".

Write the returned metadata block in the conversation.

**Output:** Artifact-context metadata block written in conversation.

---

### Step 2 — Draft All Phase Milestone Rows (Killer)

Write a table with one row per phase. Columns: Phase, Status, Definition of Done, Ticket IDs, Notes.

Write all phases in order. For each row, the Definition of Done must be an exact command or observable state — not a description.

**Required phases and their authoritative DoD commands (from CLAUDE.md):**

| Phase | Status | Definition of Done | Ticket IDs | Notes |
|---|---|---|---|---|
| Setup | Not started | All 4 verification commands pass: `ls ~/.claude-work/.../memory/` shows 6 files; `cat .claude/settings.json` shows hooks key; `grep "Session Workflow" CLAUDE.md` matches; `grep "LIVE DOCUMENT" WRITEUP.md` matches | [list SETUP ticket IDs] | Gate before MVP |
| MVP Phase 1 — Ingestion | Not started | `go test ./internal/ingestion/...` → exit 0 | [list MVP1 ticket IDs] | |
| MVP Phase 2 — Matrix/Compiler | Not started | `go test -bench=BenchmarkQuery -benchmem ./internal/matrix/` → `0 allocs/op` | [list MVP2 ticket IDs] | Blocked by Phase 1 DoD |
| MVP Phase 3 — Gate + CLI + Presenter | Not started | `./clock-gate --tick 0 [game file]` returns state; `./clock-gate --tick 999999 [game file]` returns bounded error | [list MVP3 ticket IDs] | Blocked by Phase 2 DoD |
| MVP Phase 4 — Polish + Write-up | Not started | `go vet ./...` clean; `go build ./...` clean; README has Mermaid diagrams; WRITEUP.md has all sections; LIVE DOCUMENT annotation removed | [list MVP4 ticket IDs] | Blocked by Phase 3 DoD |
| Phase 2A — Serve Mode | Not started | `./clock-gate serve` starts HTTP server; `curl localhost:8080/games` returns JSON game list | [list P2A ticket IDs] | Blocked by MVP DoD |
| Phase 2B — React UI | Not started | `npm run build` exits 0; UI loads in browser; game selector functional | [list P2B ticket IDs] | Blocked by Phase 2A DoD |

Replace `[list X ticket IDs]` with actual ticket IDs from the conversation or write "TBD" if tickets haven't been created yet.

**Output:** Phase milestone table written in conversation with all rows populated.

---

### Step 3 — Auditor Scan (Killer)

Read each row's Definition of Done column.

For each row write one line:
`Phase [name]: EXACT ([the command]) or VAGUE ([what makes it vague])`

A DoD is VAGUE if it:
- Says "implementation complete", "all tickets done", "tests pass" without naming which tests or which command
- Contains no runnable shell command or unambiguous observable state
- Could apply to any project without modification

Count VAGUE rows.

Write: "Auditor scan: [N] phases, [K] vague DoD. [List each vague phase if K > 0]"

If K > 0: replace each vague DoD with an exact command before continuing.

**Output:** Scan result written; K = 0 before proceeding.

---

### Step 4 — Draft Dependency Map and Parallel Opportunities (Killer)

Write two subsections:

**Dependency Map**
For each phase-to-phase blocking relationship, write one line in both directions:
```
[Phase A] blocks [Phase B] — [Phase B] cannot start until [Phase A] DoD passes
[Phase B] is blocked by [Phase A]
```

Required blocking relationships (from CLAUDE.md implementation phases):
- MVP Phase 1 DoD gates MVP Phase 2
- MVP Phase 2 DoD gates MVP Phase 3
- MVP Phase 3 DoD gates MVP Phase 4
- MVP Phase 4 DoD gates Phase 2A
- Phase 2A DoD gates Phase 2B
- Setup gate must pass before any MVP phase begins

**Parallel Opportunities**
For each phase containing ≥2 tickets, write an explicit parallel callout:
```
[Phase name] — Parallel: [ticket A] and [ticket B] can start simultaneously; [ticket C] blocked by [ticket A]
```

Required parallel callouts (from master plan):
- MVP Phase 1: schema.go (MVP1-1) and matrix types (MVP2-1) can be drafted in parallel — types don't depend on parser
- MVP Phase 3: gate (MVP3-1) and presenter (MVP3-2) can be built in parallel
- MVP Phase 4: --list flag (MVP4-1), README (MVP4-2), WRITEUP (MVP4-3) can be built in parallel
- Phase 2A: serve.go (P2A-1), handler.go (P2A-2), loader.go (P2A-3) can be built in parallel
- Phase 2B: components P2B-3 through P2B-7 can be built in parallel once P2B-2 is done

**Output:** Dependency map and parallel opportunity callouts written.

---

### Step 5 — Engineering Manager Scan (Killer)

Read the parallel opportunity callouts from Step 4.

List every phase that contains ≥2 tickets. For each, confirm it has an explicit "Parallel:" callout.

Write: "Engineering manager scan: [N] phases with ≥2 tickets, [K] missing parallel callout. [List missing if K > 0]"

If K > 0: add the missing parallel callout before continuing.

**Output:** Scan result written; K = 0 before proceeding.

---

### Step 6 — Invoke `confluence-page` (Killer)

Invoke the `confluence-page` skill. Pass:
- `title`: "paradox-clock-gate — Task Tracker & Milestone Doc"
- `content`: artifact-context metadata block (Step 1) + phase milestone table (Step 2, with Step 3 fixes) + dependency map and parallel opportunities (Step 4), assembled in order

The `confluence-page` skill handles space discovery, ADF formatting, ADF compliance, page creation, and URL extraction.

**Output:** `confluence-page` invocation complete; PAGE URL returned.

---

### Step 7 — Capture PAGE URL (Killer)

Read the output from Step 6. Find the line beginning "PAGE URL:".

Write that line verbatim:
`PAGE URL: [full url]`

If no "PAGE URL:" line is present: read the MCP response, extract the URL, write it.

Write: "Task tracker created. Reference this URL in the setup summary and in CLAUDE.md if needed."

**Output:** Line beginning exactly "PAGE URL:" written in conversation.

---

### Step 8 — Output Summary

Write:

```
Task tracker created: "paradox-clock-gate — Task Tracker & Milestone Doc"
PAGE URL: [url]
Phases: Setup, MVP 1–4, Phase 2A, Phase 2B (7 total)
Blocking relationships: [N] bidirectional edges mapped
Parallel opportunities: [N] callouts across phases
```

---

## Proof Requirements

- **Step 1:** `artifact-context` invocation visible.
- **Step 2:** All 7 phase rows written with Status, DoD, Ticket IDs, Notes columns.
- **Step 3:** "Auditor scan: K = 0" written before Step 4.
- **Step 4:** Dependency map with bidirectional edges; parallel callout for each phase with ≥2 tickets.
- **Step 5:** "Engineering manager scan: K = 0" written before Step 6.
- **Step 6:** `confluence-page` invocation visible.
- **Step 7:** Line beginning exactly "PAGE URL:" present in conversation.

---

## State Persistence

N/A — single session. Output is the PAGE URL written in conversation.

---

## Red Flags

| Thought | Reality |
|---------|---------|
| "The DoD is 'all tests pass' — that's clear enough" | Which tests? What command? "All tests pass" could mean `go test ./...` or `go test -bench=BenchmarkQuery -benchmem` — these produce different pass/fail signals. Name the exact command. |
| "The phases are sequential — I'll list them in order and that implies the blocking" | Implied blocking is invisible to a follow-on agent reading the tracker. Write bidirectional edges explicitly. |
| "Everyone knows schema.go and types can be drafted in parallel" | "Everyone knows" is not in the tracker. Parallel opportunities must be stated as explicit callouts — the engineering manager scan will flag their absence. |
| "The tickets don't exist yet — I'll create a tracker as a placeholder" | Confirm ticket IDs before running. Write "TBD" only where genuinely unknown; do not fabricate IDs. |
| "confluence-page returned a response — the URL is somewhere" | Step 7 requires a line beginning exactly "PAGE URL:" in the conversation. Extract and write it explicitly. |

---

## Integration

```
Called by: pre-code setup sequence (after all Jira tickets created); standalone
Calls: artifact-context, confluence-page
Output: PAGE URL written in conversation; task tracker page with 7 phases, dependency map, parallel callouts
MCP tools: invoked by confluence-page subskill
```

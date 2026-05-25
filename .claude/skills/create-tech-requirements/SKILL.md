---
name: create-tech-requirements
description: Creates the Technical Requirements Confluence page for paradox-clock-gate with embedded Go types, all 4 Mermaid diagrams, and performance requirements. Returns PAGE URL for downstream skills. TRIGGER when: user says "create the tech requirements", "write the technical requirements", "generate the tech req doc", or the pre-code setup sequence reaches this step.
---

## Overview

Technical requirements documents written from memory abbreviate type definitions, omit diagrams, and state performance requirements as prose. A follow-on agent implementing from tickets needs exact struct layouts and verifiable performance thresholds — not "see CLAUDE.md."

**Iron Law:** Every Go type must be embedded verbatim with zero abbreviated fields. No `// ...`. No "see CLAUDE.md." The document is the authoritative reference — it must be executable without opening any other file.

---

## Announce at Start

"Creating Technical Requirements doc for paradox-clock-gate. Invoking artifact-context, embedding all Go types verbatim, generating all 4 Mermaid diagrams, then creating the Confluence page. PAGE URL will be written for downstream skills."

---

## Steps

### Step 1 — Invoke `artifact-context` (Killer)

Invoke the `artifact-context` skill. Pass: phase = "Pre-code setup", component = "cli", PRD URL = [URL from create-prd output, or "Not yet created"].

Write the returned metadata block in the conversation.

**Output:** Artifact-context metadata block written in conversation.

---

### Step 2 — Embed All Go Types Verbatim (Killer)

Read `/Users/athatcher/Documents/at-proj/paradox-clock-gate/CLAUDE.md` — specifically the "Key Types" section.

Write every type definition below verbatim, with every field, no abbreviation:

**Types to embed:**
- `const MaxTick = 9001`
- `PlayType` enum with all 7 constants (`PlayTypeNone` through `PlayTypeOther`)
- `GameState` struct — all 15 fields, exact types (`uint16`, `uint8`, `[3]byte`, `uint32`, `bool`, etc.)
- `GameMeta` struct — all 4 fields
- `StateMatrix` struct — all 3 fields including `[MaxTick]GameState`
- `GameHeader` struct — all 5 fields with JSON tags
- `RawPlay` struct — all 17 fields with pointer types and JSON tags

Write each type in a Go code block. No field may be abbreviated with `// ...` or replaced with a comment.

**Output:** All 7 type definitions written verbatim in Go code blocks.

---

### Step 3 — Type Completeness Scan (Killer)

Read the type definitions written in Step 2.

For each type, count fields in the written definition. Compare against the authoritative counts from CLAUDE.md:
- `GameState`: 15 fields
- `GameMeta`: 4 fields
- `StateMatrix`: 3 fields
- `GameHeader`: 5 fields
- `RawPlay`: 17 fields
- `PlayType` enum: 7 constants

Scan for: `// ...`, `// fields omitted`, `see CLAUDE.md`, or any field count mismatch.

Write: "Type completeness scan: [N] types checked, [K] with abbreviations or missing fields. [List each issue if K > 0]"

If K > 0: correct each abbreviated type before continuing.

**Output:** Scan result written; K = 0 before proceeding.

---

### Step 4 — Generate All 4 Mermaid Diagrams (Killer)

Invoke `create-mermaid-diagram` four times — once per diagram type:

1. **Data Flow diagram** (Type 1 — `flowchart LR`): JSON file → ingestion → matrix → gate → presenter → CLI output
2. **StateMatrix Internals diagram** (Type 2 — `flowchart TB`): irregular plays → forward-fill → dense array → O(1) query
3. **Package Structure diagram** (Type 3 — `graph TD`): import dependencies between all packages
4. **Serve Mode diagram** (Type 4 — `flowchart LR`): future-state concurrent users → HTTP → StateMatrix cache → response

Write each returned Mermaid code block in the conversation as it is returned.

**Output:** 4 Mermaid code blocks written in conversation.

---

### Step 5 — Diagram Count Scan (Killer)

Count fenced Mermaid code blocks written in Step 4.

Write: "Diagram scan: [N] of 4 diagrams present."

If N < 4: identify which diagram type is missing and invoke `create-mermaid-diagram` for that type before continuing.

**Output:** Scan confirms 4 diagrams present.

---

### Step 6 — Write All Specification Sections (Killer)

Write the following sections in order. Each section must contain clock-gate-specific content — not generic template text.

**Package Specifications**
For each package (`cmd/clock-gate`, `internal/ingestion`, `internal/matrix`, `internal/gate`, `internal/presenter`):
- Responsibility (one sentence)
- Files it creates or owns
- Key function signatures it exports

**Compiler Algorithm**
Numbered steps from CLAUDE.md `Compile()` algorithm section — verbatim, naming exact constructs:
sort key `(GameClockTotalSeconds ASC, PlayID ASC)`, forward-fill loop bounds, Arena pre-allocation, `maxTick` tracking.

**Performance Requirements**
Each requirement has an exact verifying command:
```
- 0 allocs/op on query path — verified by: `go test -bench=BenchmarkQuery -benchmem ./internal/matrix/` → "0 allocs/op"
- StateMatrix size ≈ 252 KB — verified by: `go test -v -run=TestMatrixSize ./internal/matrix/` (or sizeof assertion)
- go vet clean — verified by: `go vet ./...` → exit 0
- Zero external dependencies — verified by: `go list -m all | wc -l` → "1"
```

**Data Format Specification**
JSON wrapper structure (game_id, home_team, away_team, home_score, away_score, plays array). Nullable fields list. `game_clock_total_seconds` range (0–4500 validated against dataset).

**Design Constraints Table**
All 9 locked decisions from CLAUDE.md with their rationale — verbatim, not paraphrased.

**CLI Interface**
Flag definitions, usage examples, text output format (box-drawing), JSON output keys.

**Output:** All 6 specification sections written in conversation.

---

### Step 7 — Invoke `confluence-page` (Killer)

Invoke the `confluence-page` skill. Pass:
- `title`: "paradox-clock-gate — Technical Requirements"
- `content`: artifact-context metadata block (Step 1) + type definitions (Step 2) + Mermaid diagrams (Step 4) + specification sections (Step 6), assembled in order

The `confluence-page` skill handles space discovery, ADF formatting, ADF compliance, page creation, and URL extraction.

**Output:** `confluence-page` invocation complete; PAGE URL returned.

---

### Step 8 — Capture PAGE URL (Killer)

Read the output from Step 7. Find the line beginning "PAGE URL:".

Write that line verbatim:
`PAGE URL: [full url]`

If no "PAGE URL:" line is present: read the MCP response from the confluence-page call, extract the URL, write it.

Write: "Tech Req page created. Downstream skills: record this URL as the Tech Req cross-reference in artifact-context."

**Output:** Line beginning exactly "PAGE URL:" written in conversation.

---

### Step 9 — Output Summary

Write:

```
Tech Requirements created: "paradox-clock-gate — Technical Requirements"
PAGE URL: [url]
Types embedded: GameState (15 fields), GameMeta, StateMatrix, GameHeader, RawPlay (17 fields), PlayType
Diagrams: Data Flow, StateMatrix Internals, Package Structure, Serve Mode (future)
Next: pass this URL to create-design-doc as the Tech Req reference.
```

---

## Proof Requirements

- **Step 1:** `artifact-context` invocation visible in conversation.
- **Step 2:** All 7 type definitions written in Go code blocks — no `// ...`, no "see CLAUDE.md."
- **Step 3:** "Type completeness scan: K = 0" written before Step 4.
- **Step 4:** 4 `create-mermaid-diagram` invocations visible; 4 Mermaid code blocks written.
- **Step 5:** "Diagram scan: 4 of 4" written before Step 6.
- **Step 6:** All 6 specification sections written with clock-gate-specific content.
- **Step 7:** `confluence-page` invocation visible.
- **Step 8:** Line beginning exactly "PAGE URL:" present in conversation after Step 7.

---

## State Persistence

N/A — single session. Output is the PAGE URL written in conversation for downstream skills.

---

## Red Flags

| Thought | Reality |
|---------|---------|
| "The types are in CLAUDE.md — I'll reference it instead of embedding" | A reference is not a type. A follow-on agent that opens only this document must have every field. Embed verbatim. |
| "GameState has the main fields — I'll use // ... for the rest" | `// ...` is a different struct. The type completeness scan counts fields. K must be 0. |
| "Three diagrams cover the architecture — serve mode is future state" | Serve mode exists as future state, explicitly marked with dashed styling. All 4 diagrams are required. The diagram count scan will flag 3 of 4. |
| "The performance requirement is '0 allocs/op' — that's specific enough" | A requirement without a command is a wish. Write the exact `go test -bench=BenchmarkQuery -benchmem` command with its expected output. |
| "confluence-page returned a response — the URL is in there somewhere" | "Somewhere" is not a PAGE URL written in the conversation. Step 8 requires a line beginning exactly "PAGE URL:". |
| "I'll generate the diagrams after creating the page" | Diagrams must be in the content passed to confluence-page. Generate them in Step 4, before Step 7. |

---

## Integration

```
Called by: pre-code setup sequence; standalone
Calls: artifact-context, create-mermaid-diagram (×4), confluence-page
Output: PAGE URL written in conversation for create-design-doc and create-jira-ticket
        (via artifact-context Tech Req link field)
MCP tools: invoked by confluence-page subskill
```

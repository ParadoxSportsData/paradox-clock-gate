---
name: create-prd
description: Creates the PRD Confluence page for paradox-clock-gate and returns the PAGE URL for downstream skills. TRIGGER when: user says "create the PRD", "write the product requirements", "generate the PRD", "create the requirements doc", or the pre-code setup sequence reaches the PRD step.
---

## Overview

PRDs generated without structure produce generic prose — missing non-goals, vague success criteria, no assessment framing — that downstream skills cannot link to and evaluators cannot assess. This skill enforces section completeness and specificity before any Confluence call.

**Iron Law:** All 7 sections must pass the assessment evaluator scan and lawyer scan before `confluence-page` is invoked. A page with generic content is worse than no page.

---

## Announce at Start

"Creating PRD for paradox-clock-gate. Invoking artifact-context for metadata, drafting all 7 sections with specificity scans, then creating the Confluence page. PAGE URL will be written in conversation for downstream skills."

---

## Steps

### Step 1 — Invoke `artifact-context` (Killer)

Invoke the `artifact-context` skill. Pass: phase = "Pre-code setup", component = "cli", GitHub org = "ParadoxSportsData/paradox-clock-gate".

Write the returned metadata block in the conversation.

**Output:** Artifact-context metadata block written in conversation.

---

### Step 2 — Draft All 7 PRD Sections (Killer)

Write all 7 sections below in order. Each section must name specific Go types, functions, or constraints from CLAUDE.md — not generic product language.

**Section 1 — Problem Statement**
One paragraph. State: what specific problem clock-gate solves, what `paradox-platform` does wrong (linear scans, no temporal isolation guarantee), and what the forward-fill compiler provides that a scan cannot. Name the actual failure mode being fixed (future-state bleed).

**Section 2 — Goals**
Bulleted list. Each goal is a measurable outcome:
- `go test -bench=BenchmarkQuery -benchmem` shows `0 allocs/op`
- Query at tick T returns only state from plays at or before T — never after (temporal isolation guarantee)
- Binary builds with `go build ./cmd/clock-gate/` from a clean clone with zero external dependencies
- `--tick`, `--format`, `--list` flags work per CLI spec in CLAUDE.md

**Section 3 — Non-Goals** *(minimum 3 named exclusions, each with a reason clause)*
Bulleted list. Each item names the specific thing excluded and states why:
- Format: "[specific thing] — [reason it is excluded]"
- Example: "Cumulative statistics across drives (total yards, time of possession) — requires a different data model outside this scope"
- Must have ≥3 items. Each must have a reason clause ("— because..." or "— requires..." or "— is Phase...").

**Section 4 — User Stories**
Format: `As [who], I want [what], so that [why].`
Write at minimum:
- The video sync use case (matching game state to broadcast elapsed seconds)
- The serve-mode use case (concurrent users querying the same compiled game state)
- The developer verification use case (confirming no future-state bleed at a specific tick)

**Section 5 — Business Rules**
Numbered list. Rules that constrain implementation — each traceable to a design decision:
1. Tick T query returns `HasState = false` if no play exists at or before T
2. Concurrent plays at the same tick resolve by higher `PlayID` (tiebreak rule)
3. Home/away teams come from JSON header fields (`home_team`, `away_team`) — not filename
4. Queries above `MaxTick = 9001` return a bounded error — not a zero struct

**Section 6 — Success Criteria**
Each criterion has an exact verifying command:
```
- [ ] [criterion] — verified by: `[exact command]` → [expected output]
```
At minimum: `BenchmarkQuery 0 allocs/op`, `go vet ./...` clean, smoke tests pass (tick 0, tick 1800, tick 999999 returns error), `go build` clean with no external modules.

**Section 7 — Out of Scope** *(minimum 2 named items)*
Items explicitly not part of this deliverable — named specifically, not as "future features":
At minimum: `clock-gate serve` HTTP endpoints (Phase 2A), React/TypeScript web UI (Phase 2B).

**Output:** All 7 sections drafted in conversation with clock-gate-specific content.

---

### Step 3 — Assessment Evaluator Scan (Killer)

Read each of the 7 sections. For each section write one line:
`Section [N] ([name]): SPECIFIC or GENERIC`

A section is GENERIC if it contains only these patterns not tied to clock-gate:
- "The goal of this project is to..."
- "We will implement..."
- "This system aims to..."
- "The tool should..."
- Any sentence that could apply word-for-word to a different project

Count GENERIC sections.

Write: "Assessment evaluator scan: [N] sections, [K] generic. [List each generic section if K > 0]"

If K > 0: rewrite each flagged section to name a clock-gate-specific constraint, type, or decision before continuing.

**Output:** Scan result written; K = 0 before proceeding.

---

### Step 4 — Lawyer Scan (Killer)

**Non-goals check:** Count named exclusions in Section 3. Each must include a reason clause.

Write: "Non-goals: [N] exclusions. [List any missing reason clause]"

If N < 3: add exclusions until N ≥ 3, each with a reason clause.

**Out-of-scope check:** Count named items in Section 7.

Write: "Out-of-scope: [N] named items."

If N < 2: add items until N ≥ 2.

**Output:** Lawyer scan written; non-goals ≥3 with reason clauses, out-of-scope ≥2 named items.

---

### Step 5 — Invoke `confluence-page` (Killer)

Invoke the `confluence-page` skill. Pass:
- `title`: "paradox-clock-gate — Product Requirements Document"
- `content`: artifact-context metadata block (Step 1) + all 7 sections (Step 2, with fixes from Steps 3–4 applied)

The `confluence-page` skill handles space discovery, ADF formatting, ADF compliance, page creation, and URL extraction.

**Output:** `confluence-page` invocation complete; PAGE URL returned from confluence-page.

---

### Step 6 — Capture and Write PAGE URL (Killer)

Read the output from Step 5. Find the line beginning "PAGE URL:".

Write that line verbatim on its own line:
`PAGE URL: [full url]`

If no "PAGE URL:" line is present in confluence-page output: read the MCP response from the confluence-page invocation, extract the page URL, and write it.

Write: "PRD created. Downstream skills: record this URL as the PRD cross-reference in artifact-context."

**Output:** Line beginning exactly "PAGE URL:" written in conversation.

---

### Step 7 — Output Summary

Write:

```
PRD created: "paradox-clock-gate — Product Requirements Document"
PAGE URL: [url]
Sections complete: Problem Statement, Goals, Non-Goals ([N] exclusions),
  User Stories, Business Rules, Success Criteria, Out of Scope ([N] items)
Next: pass this URL to create-tech-requirements as the PRD reference.
```

---

## Proof Requirements

- **Step 1:** `artifact-context` invocation visible in conversation before section drafting.
- **Step 2:** All 7 sections written with content naming clock-gate-specific types, functions, or constraints.
- **Step 3:** "Assessment evaluator scan: K = 0" written before Step 5.
- **Step 4:** "Non-goals: N ≥ 3" and "Out-of-scope: N ≥ 2" written before Step 5.
- **Step 5:** `confluence-page` invocation visible.
- **Step 6:** Line beginning exactly "PAGE URL:" present in conversation after Step 5.

---

## State Persistence

N/A — single session. Output is the PAGE URL written in conversation for downstream skills.

---

## Red Flags

| Thought | Reality |
|---------|---------|
| "I know what a PRD looks like — I'll skip artifact-context and write the sections directly" | Assessment context (GitHub org, contact, phase links) comes from artifact-context. Downstream skills' cross-reference fields stay "Not yet created" without it. |
| "The non-goals section has two items — that's enough" | Minimum is 3, each with a reason clause. Two items leaves scope undefined in the gap. |
| "The section reads fine to me" | "Reads fine" is not the assessment evaluator scan. Run it: does every sentence name a clock-gate-specific constraint or decision? |
| "The success criteria are checkboxes — that's sufficient" | A checkbox without an exact command is not a criterion. Each row needs the exact command: `go test -bench=BenchmarkQuery -benchmem ./internal/matrix/` — not "benchmarks pass." |
| "confluence-page returned a response — the URL is in there somewhere" | "Somewhere in the response" is not a PAGE URL in the conversation. Step 6 requires a line beginning exactly "PAGE URL:" before the skill ends. |
| "Out-of-scope is obvious — anyone would know serve mode isn't in Phase 1" | Unstated exclusions are in scope by default. Name them. |

---

## Integration

```
Called by: pre-code setup sequence; standalone
Calls: artifact-context, confluence-page
Output: PAGE URL written in conversation ("PAGE URL: [url]") for create-tech-requirements,
        create-design-doc, and create-jira-ticket (via artifact-context PRD link field)
```

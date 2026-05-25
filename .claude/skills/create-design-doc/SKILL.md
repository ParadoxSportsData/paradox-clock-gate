---
name: create-design-doc
description: Creates the architecture design doc Confluence page for paradox-clock-gate with decision narratives, rejected alternatives, and all 4 Mermaid diagrams. Returns PAGE URL for downstream skills. TRIGGER when: user says "create the design doc", "write the architecture doc", "generate the design document", or the pre-code setup sequence reaches this step.
---

## Overview

Design documents that list decisions without naming rejected alternatives read as specs, not engineering judgment. This skill enforces ADR-style decision narratives — every choice documented alongside what was rejected and why — before any Confluence call.

**Iron Law:** Every locked design decision must name the specific alternative that was rejected and the reason the chosen option wins in terms of clock-gate's constraints. "It's simpler" fails. "Binary search requires a runtime conditional on every query; the flat array eliminates conditionals and makes the correctness proof trivial" passes.

---

## Announce at Start

"Creating architecture design doc for paradox-clock-gate. Invoking artifact-context, writing ADR-style decision narratives for all 9 locked decisions, embedding all 4 Mermaid diagrams with connecting sentences, then creating the Confluence page. PAGE URL will be written for downstream skills."

---

## Steps

### Step 1 — Invoke `artifact-context` (Killer)

Invoke the `artifact-context` skill. Pass: phase = "Pre-code setup", component = "cli", PRD URL and Tech Req URL from prior skill outputs (or "Not yet created").

Write the returned metadata block in the conversation.

**Output:** Artifact-context metadata block written in conversation.

---

### Step 2 — Write Decision Narratives for All 9 Locked Decisions (Killer)

For each of the 9 locked design decisions, write a block in this format:

```
**[Decision topic]**
Chosen: [what was chosen]
Rejected: [the specific alternative that was considered]
Why chosen wins: [reason stated in terms of clock-gate's constraints — not generic principles]
Trade-off owned: [what was given up and why it's acceptable at this scale]
```

Write all 9 blocks. The 9 decisions (from CLAUDE.md Design Constraints table):
1. Language — Go vs. alternatives
2. CLI framework — `flag` stdlib vs. Cobra
3. Lookup strategy — flat array O(1) vs. binary search
4. Runtime allocations — zero after init vs. per-query allocation
5. WinProb storage — `uint16` vs. `float64`
6. Team abbreviations — `[3]byte` vs. `string`
7. Description storage — arena allocator vs. per-play strings
8. Output scope — game state only vs. cumulative stats
9. Phase 2A backend — `net/http` vs. Gin/Echo

**Output:** All 9 decision blocks written in conversation.

---

### Step 3 — ADR Scan (Killer)

Read each of the 9 decision blocks from Step 2.

For each block write one line:
`Decision [N] ([topic]): COMPLETE or INCOMPLETE`

A block is INCOMPLETE if:
- "Rejected" names a category ("other approaches", "alternatives") instead of a specific option
- "Why chosen wins" states a generic principle ("simpler", "better") without naming a clock-gate-specific constraint
- "Trade-off owned" is absent or says "none"

Count INCOMPLETE blocks.

Write: "ADR scan: [N] decisions, [K] incomplete. [List each incomplete block if K > 0]"

If K > 0: rewrite each incomplete block before continuing.

**Output:** ADR scan result written; K = 0 before proceeding.

---

### Step 4 — Generate and Embed All 4 Diagrams with Connecting Sentences (Killer)

For each diagram type, write a connecting sentence first, then invoke `create-mermaid-diagram`, then embed the returned code block.

**Diagram 1 — Data Flow** (`flowchart LR`)
Connecting sentence: "This diagram shows the full pipeline from raw JSON to CLI output — each stage is a pure transformation with no shared state, which is what makes the temporal isolation guarantee possible."
Invoke `create-mermaid-diagram` for Type 1. Embed returned code block.

**Diagram 2 — StateMatrix Internals** (`flowchart TB`)
Connecting sentence: "This diagram shows why the O(1) lookup claim is a mathematical guarantee rather than an optimization — every second between 0 and MaxTick is pre-assigned during compilation, so query time is a single array dereference with no conditional logic."
Invoke `create-mermaid-diagram` for Type 2. Embed returned code block.

**Diagram 3 — Package Structure** (`graph TD`)
Connecting sentence: "This diagram shows the dependency graph — `internal/` packages form a directed acyclic graph, which is why any package can be tested in isolation without loading the full pipeline."
Invoke `create-mermaid-diagram` for Type 3. Embed returned code block.

**Diagram 4 — Serve Mode** (`flowchart LR`, future-state nodes dashed)
Connecting sentence: "This diagram shows the serve mode architecture (Phase 2A, future state) — one compiled StateMatrix per game, shared across all concurrent users, with the GC-free query path that makes horizontal read scaling trivial."
Invoke `create-mermaid-diagram` for Type 4. Embed returned code block.

**Output:** 4 connecting sentences written; 4 `create-mermaid-diagram` invocations visible; 4 Mermaid code blocks embedded.

---

### Step 5 — Diagram-Decision Connection Scan (Killer)

Read the 4 diagram sections from Step 4.

For each diagram, confirm: does the connecting sentence name a specific decision or property from the decision narrative (not just describe what the diagram shows visually)?

Write one line per diagram:
`Diagram [N] ([type]): CONNECTED — references [decision/property named] or DISCONNECTED`

Count DISCONNECTED diagrams.

Write: "Diagram connection scan: [N] of 4 connected."

If any diagram is DISCONNECTED: rewrite its connecting sentence to reference the specific decision it illustrates before continuing.

**Output:** Scan result written; all 4 diagrams connected.

---

### Step 6 — Invoke `confluence-page` (Killer)

Invoke the `confluence-page` skill. Pass:
- `title`: "paradox-clock-gate — Architecture Design Document"
- `content`: artifact-context metadata block (Step 1) + problem framing paragraph + all 9 decision blocks (Step 2) + all 4 diagram sections with connecting sentences (Step 4), assembled in order

Add a brief problem framing paragraph before the decision blocks: why clock-gate exists, what temporal isolation means, what the forward-fill compiler guarantees (2–3 sentences — not a copy of the PRD problem statement, but the architectural framing).

The `confluence-page` skill handles space discovery, ADF formatting, ADF compliance, page creation, and URL extraction.

**Output:** `confluence-page` invocation complete; PAGE URL returned.

---

### Step 7 — Capture PAGE URL (Killer)

Read the output from Step 6. Find the line beginning "PAGE URL:".

Write that line verbatim:
`PAGE URL: [full url]`

If no "PAGE URL:" line is present: read the MCP response from the confluence-page call, extract the URL, write it.

Write: "Design doc created. Downstream skills: record this URL as the Architecture Diagram cross-reference in artifact-context."

**Output:** Line beginning exactly "PAGE URL:" written in conversation.

---

### Step 8 — Output Summary

Write:

```
Design doc created: "paradox-clock-gate — Architecture Design Document"
PAGE URL: [url]
Decisions documented: 9 (each with rejected alternative and win reason)
Diagrams: Data Flow, StateMatrix Internals, Package Structure, Serve Mode (future)
Next: this URL feeds into artifact-context diagram link field for all subsequent skills.
```

---

## Proof Requirements

- **Step 1:** `artifact-context` invocation visible.
- **Step 2:** All 9 decision blocks written with Chosen / Rejected / Why chosen wins / Trade-off owned fields.
- **Step 3:** "ADR scan: K = 0" written before Step 4.
- **Step 4:** 4 connecting sentences written; 4 `create-mermaid-diagram` invocations visible; 4 code blocks embedded.
- **Step 5:** "Diagram connection scan: 4 of 4 connected" written before Step 6.
- **Step 6:** `confluence-page` invocation visible.
- **Step 7:** Line beginning exactly "PAGE URL:" present in conversation.

---

## State Persistence

N/A — single session. Output is the PAGE URL written in conversation.

---

## Red Flags

| Thought | Reality |
|---------|---------|
| "The rationale is 'zero dependencies' — that's specific enough" | Zero dependencies is a principle, not a clock-gate-specific constraint. State what the dependency would have cost: "Cobra adds 3 transitive imports; a reviewer can audit `flag` in one file." |
| "Rejected: other CLI frameworks" | Name the specific framework: "Rejected: Cobra." A category is not a rejected alternative. |
| "The diagram speaks for itself — no connecting sentence needed" | A diagram without a connecting sentence is decoration. The scan counts disconnected diagrams; 4 of 4 must be connected. |
| "Trade-off owned: none" | Every design choice trades something. Name it. "Trade-off: 252 KB per loaded game vs. ~12 KB for a binary-searched sorted list — acceptable because 300 MB for 1000 simultaneously popular games fits in a single server." |
| "confluence-page returned a response — the URL is in there" | "In there" is not a line beginning "PAGE URL:" in the conversation. Step 7 requires the exact format. |
| "The serve mode diagram is future state — I'll skip it" | Future-state nodes are marked with dashed styling, not omitted. The diagram is required. The diagram count scan in create-tech-requirements already enforced this; create-design-doc enforces it again. |

---

## Integration

```
Called by: pre-code setup sequence; standalone
Calls: artifact-context, create-mermaid-diagram (×4), confluence-page
Output: PAGE URL written in conversation for artifact-context diagram link field
MCP tools: invoked by confluence-page subskill
```

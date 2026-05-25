---
name: jira-ticket-template
description: Produces a complete agent-executable Jira ticket body from structured inputs. TRIGGER when: create-jira-ticket needs to format a ticket body, or user says "format this ticket", "write the ticket body", "create a ticket for [phase/feature]".
---

## Overview

Tickets that lack an embedded algorithm, inline types, or verifiable acceptance criteria require an agent to ask clarifying questions — defeating the purpose of the ticket. This skill produces ticket bodies an agent can execute with zero follow-up.

**Iron Law:** An agent reading only this ticket must be able to implement it exactly, producing the specified output, without reading any other file.

---

## Announce at Start

"Formatting ticket body for: [title]. Building all sections — types and algorithm embedded inline so the ticket is executable without external lookups."

---

## Sections

Every ticket body must contain all 10 sections below, in order. Omit none.

---

### Section 1 — TITLE (Killer)

Format: `[Verb] [specific thing] — [measurable outcome]`

- **Verb:** Create / Implement / Add / Fix (not "Work on" or "Handle")
- **Specific thing:** the exact file or function (e.g., `internal/ingestion/parser.go — ParseFile()`)
- **Measurable outcome:** the definition of done in one phrase (e.g., `go test ./internal/ingestion/... passes`)

Write the title as the first line of the ticket body.

**Output:** One title line.

---

### Section 2 — CONTEXT (Killer)

Write exactly 2–3 sentences:
1. What this component is and what it does
2. Why it exists in the system — its role in the larger pipeline
3. Where it sits in the package structure and what calls it

Do not copy the ticket title. Do not write "this ticket will..."

**Output:** 2–3 sentence context paragraph.

---

### Section 3 — FILE (Killer)

List every file the agent must create or modify:

```
Create: [exact path from repo root]
Modify: [exact path from repo root]
```

**Output:** One line per file with Create/Modify label and exact path.

---

### Section 4 — DEPENDS ON (Killer)

List every prerequisite by ticket ID and the specific artifact the agent needs from it:

```
[TICKET-ID]: [specific type definition, function, or interface that must exist]
```

If there are no dependencies: write "None — implement standalone."

Do not write "see ticket X" — name the specific artifact.

**Output:** One line per dependency, or "None."

---

### Section 5 — RELEVANT TYPES + FUNCTION SIGNATURES (Killer)

Embed every Go type definition and function signature the agent needs — inline, verbatim. The agent must not need to open any other file to find these.

Format:
```go
// From internal/matrix/types.go
type GameState struct {
    // all fields, no abbreviation
}

// Function to implement:
func ParseFile(path string) ([]RawPlay, GameHeader, error)
```

Include all struct fields — do not abbreviate with `// ...`.

**Output:** All relevant types and target signatures in a Go code block.

---

### Section 6 — ALGORITHM (Killer)

Write numbered implementation steps. Each step names the exact data structure, method, or constant to use — not "process the plays" but "iterate plays sorted by `(GameClockTotalSeconds ASC, PlayID ASC)`."

```
1. [Step — name the exact construct]
2. [Step — include sort order, loop bounds, tiebreak logic as applicable]
...
```

Include explicitly:
- Sort order if processing a slice
- Forward-fill loop bounds and condition if applicable
- How null pointer fields (`*int`, `*string`) must be handled
- How concurrent events at the same tick are resolved (higher PlayID wins)
- What NOT to do if the omission would surprise an implementer

**Output:** Numbered algorithm steps with no ambiguous "process" or "handle" verbs.

---

### Section 7 — CONSTRAINTS

List explicit prohibitions with reasons:

```
- Do NOT [thing]: [reason]
```

Always include these when applicable to the file being changed:
- No external dependencies (any Go source file): "clock-gate is zero-dependency; stdlib only"
- No pointers in `GameState` fields (if modifying hot struct): "pointers cause GC pressure on the query path"
- No per-query heap allocations (if modifying query path): "0 allocs/op claim must hold"

**Output:** Bulleted constraint list, or "No additional constraints beyond CLAUDE.md design constraints."

---

### Section 8 — ACCEPTANCE CRITERIA (Killer)

List specific, testable conditions. Each criterion has an exact verifying command producing a binary result.

Format:
```
- [ ] [Condition] — verified by: `[exact command]` → [expected output or pattern]
```

Acceptable verifying outputs: `ok github.com/...`, `0 allocs/op`, exit code 0, specific stderr/stdout text.
Not acceptable: "works correctly," "tests pass" (which tests?), "no errors."

**Output:** Checklist of criteria, each with an exact shell command.

---

### Section 9 — DEFINITION OF DONE (Killer)

One command + one expected output. This is the phase gate the calling skill uses to verify completion.

```
Command: [exact command to run from repo root]
Expected output: [exact string, regex pattern, or "exit 0"]
```

**Output:** Single Command + Expected output pair.

---

### Section 10 — METADATA

```
BLOCKS: [comma-separated ticket IDs this ticket blocks, or "none"]
LABELS: [phase] [type: implementation|test|doc|setup] [component: ingestion|matrix|gate|presenter|cli]
```

**Output:** BLOCKS and LABELS lines.

---

## Proof Requirements

- **Section 1:** Title has all three parts — Verb, specific thing, measurable outcome.
- **Section 3:** Every file has a Create or Modify label and an exact path from repo root.
- **Section 4:** Every dependency names the specific artifact (type, function, interface) — not just the ticket ID.
- **Section 5:** All type definitions and signatures embedded inline in a Go code block — no "see CLAUDE.md."
- **Section 6:** Every algorithm step names the exact construct to use — zero "process" or "handle" verbs without specifics.
- **Section 8:** Every criterion has an exact shell command producing a binary result.
- **Section 9:** Exactly one Command + Expected output pair.

---

## State Persistence

N/A — single session. Output is the formatted ticket body returned to create-jira-ticket.

---

## Red Flags

| Thought | Reality |
|---------|---------|
| "The type definitions are in CLAUDE.md — I'll reference them" | Military order anchor: the ticket must be executable without opening any other file. Embed the types inline. |
| "The algorithm is straightforward — an engineer would know it" | Each step must name the exact construct. "Sort plays" fails; "Sort plays by (GameClockTotalSeconds ASC, PlayID ASC)" passes. |
| "The acceptance criterion is 'tests pass'" | Which tests? What command? What output? Every criterion needs an exact shell command. |
| "I'll use '// ...' to abbreviate the type definition" | An abbreviated type is not the type. An agent will invent the missing fields. Embed in full. |
| "The constraint is implied by the project's zero-dep rule" | The ticket is self-contained. State the constraint explicitly with its reason. |
| "I'll add a note to refer to CLAUDE.md for the algorithm" | CLAUDE.md is not available to an agent that has only the ticket. Write the algorithm steps here. |

---

## Integration

```
Called by: create-jira-ticket (once per ticket)
Calls: none
Output: Formatted 10-section ticket body string, returned to create-jira-ticket
        for submission via Jira MCP tool
```

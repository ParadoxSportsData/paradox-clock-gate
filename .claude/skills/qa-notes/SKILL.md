---
name: qa-notes
description: Generates QA testing notes from code evidence and posts them to a Jira ticket. TRIGGER when: user says "add QA notes", "write QA notes", "QA notes for", "add testing notes", "document QA coverage", "write test notes", "add notes to the ticket", "generate QA notes", or when implementation of a phase is complete and QA documentation is needed.
---

# QA Notes

## Overview

QA notes written from descriptions miss callers, omit pass/fail criteria, and leave regression as vague names that cannot be acted on. This skill starts from code evidence — the diff — not descriptions.

**Executor:** AI agent, post-implementation of a phase or feature. Invoked standalone.

**Iron Law:** The diff is ground truth. Every scenario must trace to a modified function or its caller.

## Announce at Start

"Generating QA notes for [TICKET-ID]. I'll map all impacted code paths from the diff, identify what the change branches on, and build scenario tables with explicit pass/fail outcomes. Output is visible in the conversation before posting to Jira."

## Steps

**Step 1 — Read ticket and diff**

Read the Jira ticket for context. Identify the files changed. Write a list of every modified function, method, and interface. This list drives every subsequent checkpoint.

---

**☑ CHECKPOINT 1 (Gawande) — Build caller-to-entry-point map**

For every function in the Step 1 list:
- Grep the repo for all callers, excluding test files (`*_test.go`)
- For each caller: determine if it is an entry point — code called directly by the runtime rather than by application code. For clock-gate, entry points are: `cmd/clock-gate/main.go` (CLI), any `_test.go` `TestXxx` or `BenchmarkXxx` function.
- If the caller is NOT an entry point, grep for that caller's callers and repeat until every chain reaches an entry point or terminates in test-only code.

Write the map in the conversation, showing the full chain:
```
modified_function → caller → ... → entry point (type)
```

> **PAUSE — Map must be written in the conversation before continuing.**

---

**☑ CHECKPOINT 2 — Identify variable dimension(s)**

Read the bodies of the modified functions. Determine what the code actually branches on — the axis of variation that produces distinct test scenarios. Examples: tick value (before first play, mid-game, OT), format flag (text vs. json), play type (kickoff, run, pass), null field handling.

Write one line per dimension:
`Variable dimension: [name] — [the specific condition the code checks]`

If the code does not branch (single linear path): write `Variable dimension: none — single path; scenarios vary by caller only`.

> **PAUSE — Dimension line(s) must be posted before building tables.**

---

**☑ CHECKPOINT 3 — Build Changed Components tables**

For each modified component (function, package, or CLI flag):

1. Write one line: what it does differently now (state the behavioral change, not a copy of the ticket description)
2. Build a table. Use the dimension name from Checkpoint 2 as the column header:

| Scenario | [Variable Dimension] | Expected Result |
|---|---|---|
| [description] | [dimension value] | [exact outcome — output text, return value, error type, or benchmark result] |

One row per distinct scenario. Core component under test comes first; secondary impacted components follow in their own tables.

---

**☑ CHECKPOINT 4 — Build Edge Cases table**

Build a separate table for conditions the happy-path rows do not cover:

| Scenario | Condition | Expected Result |
|---|---|---|
| [edge case] | [boundary or error condition] | [exact outcome] |

Add a row for each of the following that applies to this change:
- Null or missing JSON fields in play data
- Tick value of 0 (before any play)
- Tick value greater than game's MaxTick
- OT plays (quarter >= 5, tick > 3600)
- Concurrent plays at the same tick (play_id tiebreak)
- Zero-allocation claim: `BenchmarkQuery` must show `0 allocs/op`
- Any pre-existing behavior exposed or changed by the diff

---

**☑ CHECKPOINT 5 (Gawande) — Audit caller coverage**

List every caller from Checkpoint 1. Count how many have ≥1 row in either the Changed Components or Edge Cases tables.

For any caller with zero rows: add a row before continuing.

Post the result: `Caller coverage: [N] of [N] callers have ≥1 row.`

> **PAUSE — Post the coverage count before continuing.**

---

**☑ CHECKPOINT 6 (QA Test Case Author) — All rows must be self-contained**

For every row in every table: "Could someone who has never read this code determine pass/fail from this row alone?"

For any row where the answer is no: rewrite the Expected Result until the answer is yes.

Acceptable outcomes: exact output text, specific return values, specific error types, specific benchmark numbers (`0 allocs/op`), specific exit codes.
Not acceptable: "works correctly", "returns successfully", "behaves as expected".

---

**☑ CHECKPOINT 7 (Kent Beck) — Build Regression section**

For each code path touched but whose behavior must not change, write one line:
`ComponentName — [specific behavior that must still work]`

Apply Kent Beck test to each item: "If this change were fully reverted, would QA catch the rollback from this item alone?" Any item that is a bare component name fails — rewrite it as a behavioral assertion.

---

**☑ CHECKPOINT 8 — Verify artifact accuracy**

For every testable artifact referenced in any table (CLI flags, output format fields, type names, benchmark names), verify against the authoritative source before including in the posted comment.

| Artifact | Authoritative Source | How to Verify |
|---|---|---|
| CLI flag names and types | `cmd/clock-gate/main.go` | Read the `flag.String/Int/Bool` declarations |
| Text output format | `internal/presenter/snapshot.go` | Read the `RenderText` function |
| JSON output keys | `internal/presenter/snapshot.go` | Read the `RenderJSON` function |
| Type definitions (`GameState`, `StateMatrix`, etc.) | `internal/matrix/types.go` | Read the type declarations |
| Benchmark names | `internal/matrix/compiler_test.go` or equivalent | Read the `Benchmark*` function signatures |
| `go test` command and flags | `CLAUDE.md` Build & Test Commands section | Read the exact commands listed |

For each artifact in the tables:
1. Identify which type it is
2. Read the authoritative source file
3. Confirm the artifact matches the source exactly
4. If a mismatch is found: correct it before continuing

> **PAUSE — List every artifact checked and the source file used. Post this list in the conversation before composing the Jira comment.**

---

**Step 2 — Compose and post**

Compose the Jira comment in ADF format. Sections in order:

1. **In Scope** — 1–2 sentences only
2. **Changed Components** — tables from Checkpoint 3, with dimension rationale line above each table
3. **Edge Cases** — table from Checkpoint 4
4. **Regression** — items from Checkpoint 7
5. **Ticket** — link to the Jira ticket only

Post using `mcp__plugin_greenfi-engineering_atlassian-custom__jira_add_comment`.

## Proof Requirements

The posted Jira comment contains:
- In Scope: 1–2 sentences, no more
- Changed Components: ≥1 table; every Expected Result cell contains a specific outcome (not blank, not "N/A", not "works correctly")
- Edge Cases: ≥1 table with ≥1 row
- Regression: ≥1 item; zero items are bare component names without a behavioral assertion
- Ticket link present

Visible in the conversation (not required in the Jira comment):
- Caller-to-entry-point map and dimension rationale lines
- Artifact verification list (Checkpoint 8): every testable artifact with its named source file

## State Persistence

N/A — single session, no persistence required.

## Red Flags

| Thought | Reality |
|---------|---------|
| "The ticket description covers the scenarios" | Ticket descriptions are summaries. Callers not mentioned are still touched. Start from the diff. |
| "I know what the variable dimension is" | The dimension comes from reading the if/switch branches in the modified function bodies, not from domain knowledge. |
| "The expected result is obvious — I don't need to write it" | Obvious to whom? Every row needs a stated outcome. |
| "I'll list the component names in regression" | A name is not an assertion. "Compile" tells QA nothing. "Compile — forward-fill still produces HasState=true for tick 10 when only tick 5 has a play" does. |
| "I've covered the main path — edge cases are minor" | Edge cases are where regressions hide. Checkpoint 4 is not optional. |
| "There's only one caller to worry about" | Run the grep. CLI tools have multiple internal callers. |
| "I know the flag name" | Flag names come from `main.go`, not from memory. Read the source before writing it in a row. |
| "The tables look complete" | Run Checkpoint 5. Looks complete ≠ coverage complete. |
| "The zero-allocation claim is a benchmark concern, not a QA concern" | `0 allocs/op` is a testable acceptance criterion verifiable with `go test -bench -benchmem`. It belongs in the Edge Cases table if any change touches the query path. |

## Integration

```
Called by: standalone after phase implementation; or after create-jira-ticket when documenting QA coverage
Calls: mcp__plugin_greenfi-engineering_atlassian-custom__jira_add_comment (to post the comment)
Output: Jira comment on [TICKET-ID] — Changed Components tables, Edge Cases table, Regression items, ticket link
```

---
*Imported from GreenFi marketplace `qa-notes` skill. Modified for clock-gate: removed k8s/kustomize/gRPC/SQS infrastructure references; adapted Checkpoint 8 artifact table for Go CLI (main.go flags, presenter output, matrix types); adapted entry point definition for CLI (main.go); adapted edge cases for clock-gate domain (OT, null fields, tick bounds, zero-alloc claim).*

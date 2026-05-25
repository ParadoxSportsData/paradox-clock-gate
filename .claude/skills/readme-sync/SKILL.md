---
name: readme-sync
description: Generates specific paste-ready README update text by comparing current source code against the existing README. Does NOT write to README — presents updates for review only. TRIGGER when: capture-decision Step 9 fires (interface change confirmed), pre-commit hook warns that .go files changed without README update, user says "sync the README", "README is stale", "update the README", "what changed in the README".
---

## Overview

A changed CLI flag or output format that isn't reflected in the README is a credibility gap at submission time. This skill finds every divergence between code reality and README content and produces update text Aaron can paste directly.

**Iron Law:** Generate the diff. Present it. Do not write to README.md. Aaron applies updates; the skill does not.

---

## Announce at Start

"Invoking readme-sync. I'll read the actual source files, extract current interface facts, compare against the README, and produce paste-ready update text for every gap."

---

## Steps

### Step 1 — Establish What Changed

State what interface change triggered this invocation:
- Which flags changed (added, removed, renamed, default changed)?
- Which subcommands changed?
- Did output format change?
- Did install instructions change?

If invoked from pre-commit hook warning with no specific change stated: write "Invoked from pre-commit warning — scanning all interface sections."

**Output:** One-sentence trigger statement before proceeding.

---

### Step 2 — Check README Existence (Killer)

Run: `ls /Users/athatcher/Documents/at-proj/paradox-clock-gate/README.md`

Two paths:
- **File exists:** Proceed to Step 3.
- **File does not exist:** Write "README does not exist — generating initial README content from scratch." Skip Step 5 comparison; Steps 6–7 generate full initial content instead of diffs.

**Output:** README existence stated explicitly before proceeding.

---

### Step 3 — Read Actual Source Files (Killer)

Read these files using the Read tool — not from memory, not from prior session context:

1. `/Users/athatcher/Documents/at-proj/paradox-clock-gate/cmd/clock-gate/main.go` — CLI flags, subcommands, usage string, error messages
2. `/Users/athatcher/Documents/at-proj/paradox-clock-gate/internal/presenter/snapshot.go` — output format fields, box-drawing characters, display labels
3. `/Users/athatcher/Documents/at-proj/paradox-clock-gate/CLAUDE.md` — CLI interface section, build commands, package structure

If any file does not exist yet: write "FILE NOT YET WRITTEN: [path] — using CLAUDE.md spec as source of truth for this file's interface."

**Output:** All three files read (or absence noted) before proceeding.

---

### Step 4 — Extract Interface Facts (Killer)

From the files read in Step 3, write this explicit list in the conversation:

```
FLAGS:
- --tick <int>: [description, default if any]
- --format <text|json>: [description, default]
- --list: [description]
[add any others present in source]

SUBCOMMANDS / MODES:
- [list any, or "none"]

OUTPUT FORMAT (text mode):
- Header line: [exact format]
- Score line: [exact format]
- Ball/possession line: [exact format]
- Win prob line: [exact format]
- Last play line: [exact format]

OUTPUT FORMAT (json mode):
- Top-level keys: [list]

USAGE EXAMPLES (from code):
- [exact invocations shown in usage string or CLAUDE.md]

BUILD / INSTALL:
- Build command: [exact]
- Binary location after build: [exact]
```

Do not infer or fill from memory. If a field is not present in the source files, write "NOT IN SOURCE."

**Output:** The completed facts list written in the conversation, with explicit source attribution.

---

### Step 5 — Compare README Against Facts (Killer)

Read README.md in full.

For each README section that covers interface content — flags, usage examples, output format, install instructions, build commands, architecture overview — compare against the facts list from Step 4.

For each section write:
```
Section: [README section name]
README says: [exact current content, quoted]
Code says: [exact fact from Step 4]
Status: MATCH | STALE | MISSING IN README | NOT IN SCOPE
```

Sections outside interface content (project motivation, contributing, license) → Status: NOT IN SCOPE. Do not compare these.

**Output:** Full comparison table written in the conversation.

---

### Step 6 — First-Time User Scope Gate (Killer)

For every gap marked STALE or MISSING IN README in Step 5:

Ask: would a developer who has never seen this codebase, reading only this README section, take the wrong action or be unable to use the tool?

- If **yes**: keep in the update list.
- If **no**: write "Excluded — [specific reason first-time user is not affected]" and remove from list.

Do not exclude based on how minor the change seems. Exclude only when a first-time user is genuinely unaffected.

**Output:** Final update list with any exclusions and their reasons stated.

---

### Step 7 — Generate Paste-Ready Update Text (Killer)

For each gap in the final list from Step 6, produce:

```
## UPDATE: [README section name]

REPLACE THIS:
"""
[exact current README content for this section]
"""

WITH THIS:
"""
[complete replacement content, ready to paste — no placeholders, no [...] fragments]
"""

REASON: [one sentence — what changed in code that requires this update]
```

If the README section is MISSING (no current content): use "ADD THIS:" instead of "REPLACE THIS / WITH THIS."

The replacement content must be complete. Aaron pastes it and is done.

**Output:** One UPDATE block per gap, all written in the conversation.

---

### Step 8 — Declare Scope

Write: "Sections not checked in this invocation: [list of README sections outside interface scope]. To check those sections, invoke readme-sync with explicit scope."

**Output:** Explicit scope declaration written.

---

### Step 9 — Present; Do Not Write (Killer)

Write: "readme-sync complete. [N] update(s) above. Apply any you approve — I will not write to README.md."

Do not call Write or Edit on README.md. Do not apply changes on Aaron's behalf unless he explicitly says to do so in a follow-up message.

**Output:** Completion statement written. No file writes to README.md.

---

## Proof Requirements

- **Step 2:** README existence stated explicitly — "exists" or "does not exist."
- **Step 3:** All three source files read with explicit Read tool calls visible in conversation; no file reads skipped or substituted with memory.
- **Step 4:** Interface facts list written in conversation with "NOT IN SOURCE" for absent fields; no inferred content.
- **Step 5:** Full comparison table written — every README interface section has an explicit Status row.
- **Step 6:** Every STALE or MISSING gap reviewed; exclusions have explicit first-time-user reasoning.
- **Step 7:** Each UPDATE block contains complete replacement content — no placeholder text, no "[...]" fragments.
- **Step 9:** No Write or Edit tool call targeting README.md appears in this invocation.

---

## State Persistence

N/A — single session, all outputs are conversation text and presented update blocks.

---

## Red Flags

| Thought | Reality |
|---------|---------|
| "I know what the flags are — no need to read main.go" | The journalist anchor requires a source citation. Memory is not a source. Read the file. |
| "The README looks pretty complete" | "Looks complete" is confirmation bias — the README confirming itself. The Step 4 facts list is the source of truth. |
| "This change is too minor to mention in the README" | The first-time user gate in Step 6 makes this call, not the agent. Surface the gap; let the gate decide. |
| "I'll generate a quick summary of what needs updating" | Quick summaries are not paste-ready. Aaron cannot use them without additional work. Step 7 requires complete replacement content. |
| "I'll just apply the changes since they're obvious" | The Iron Law: generate and present. Writing to README.md without Aaron's review is a scope violation. |
| "The README section is mostly right — I'll note the one wrong line" | Generate the complete replacement for the section. A partial snippet with surrounding context removed is harder to apply than the stale original. |
| "README doesn't exist yet, so there's nothing to sync" | Nothing to diff means everything is a gap. Step 2 routes to full initial content generation, not a skip. |

---

## Integration

```
Called by: capture-decision (Step 9, when interface change confirmed); pre-commit hook warning
Calls: none
Output: Paste-ready UPDATE blocks presented in conversation for each README gap; no files written
```

---
name: coding-standards-go
description: Go implementation standards for paradox-clock-gate. TRIGGER when: about
  to write or modify any .go file, starting GREEN phase, user says "implement",
  "add function", "write the code", "create [file].go", or any code-writing task
  in the paradox-clock-gate repo.
---

## Overview

Without this skill, Go code that compiles and passes tests silently accumulates
debt: bare error returns that lose call-site context, magic strings in three files,
functions that mix concerns, duplicated logic, and log lines that can't be queried
in production.

**Executor:** AI agent, before writing any .go function body in paradox-clock-gate.

**Iron Law:** No function body is written until the Pike gate (Step 2) and Constants
gate (Step 3) results are written in conversation. No commit until the 7-item
compliance checklist (Step 9) shows all PASS.

---

## Announce at Start

"Applying coding-standards-go to [file/function]. Seven gates before commit:
Pike → Constants → Fowler → implement → Cheney → Uncle Bob → Bourgon.
Each gate result written in conversation."

---

## Steps

### Step 1 — State Scope

Write: `"Scope: [file path] — [function name(s) to implement]"`

**Output:** Scope written before any gate.

---

### Step 2 — Pike Gate: Name Before Body (KILLER)

For each function to be written, write its intended name as a single `VerbNoun`
identifier before writing any body code.

Examples: `ParseFile`, `CompileMatrix`, `ValidateTick`, `RenderText`.

If the intended behavior requires "and" anywhere in a natural description of
what the function does — split it into two functions with separate names.
Do not write either body until both names pass this gate.

Write in conversation:
```
Pike gate:
  [FunctionName] — single-purpose: PASS
  [FunctionName] — FAIL: splits into [NameA] + [NameB]
```

**Output:** Pike gate written; zero FAIL entries remain before Step 3.

---

### Step 3 — Constants Gate (KILLER)

List every string literal the planned implementation will contain.

For each literal:
- Used more than once in this changeset → `const`
- Referenced in any `_test.go` file → `const`
- Used exactly once, never in tests → inline permitted

Define all `const` literals at package level in the relevant `.go` file
before writing any function body.

Write in conversation:
```
Constants gate:
  "game_clock_total_seconds" — 3 uses — const GameClockField
  "GB" — test-only — const TeamGB
  "text" — 1 use, no test reference — inline permitted
```

**Output:** Constants gate table written; all promoted literals declared
as `const` before function body is written.

---

### Step 4 — Fowler Duplication Scan (KILLER)

Before writing any new logic block longer than 3 lines, run:

```bash
grep -rn "[first distinctive identifier or type from planned logic]" ./internal/ ./cmd/
```

If a matching logic block of >3 lines already exists anywhere in the repo:
1. Write: `"Fowler: duplication at [file:line] — extracting to [FunctionName]"`
2. Write the extracted function first
3. Replace the existing inline instance with a call to it
4. The new implementation calls the same extracted function

Write in conversation:
```
Fowler scan: [no duplication found]
  — or —
Fowler scan: duplication at internal/matrix/compiler.go:47 — extracted sortPlays()
```

**Output:** Fowler scan result written; any extracted function exists before
new logic is written.

---

### Step 5 — Write Implementation

Write the function body or bodies from Step 1.

If the body reaches 25 lines, split at the first operation that can be given
an independent name before continuing.

**Output:** Implementation written.

---

### Step 6 — Cheney Error Scan (KILLER)

Read every line of the written code. For each `return err` and `_ = err`:

**`return err`** — rewrite as:
```go
return fmt.Errorf("packagename.FunctionName: %w", err)
// when an input value is diagnostic:
return fmt.Errorf("packagename.FunctionName(%v): %w", inputValue, err)
```

**`_ = err`** — not permitted. Either return the wrapped error, or log with
`slog.Error` and add a one-line comment naming the reason swallowing is correct.

Write in conversation:
```
Cheney scan:
  2 bare `return err` → wrapped with fmt.Errorf context
  0 `_ = err`
  Result: clean
```

**Output:** Cheney scan written; zero bare `return err`; zero silent `_ = err`
in the changeset.

---

### Step 7 — Uncle Bob Abstraction Scan (KILLER)

For each function in the changeset, read the body and answer both questions:

1. Does the body contain a decision: `if`, `for`, `switch`, or `select`?
2. Does the body contain a direct I/O call: `os.*`, `json.*`, `http.*`,
   `fmt.Fprint*`, or `bufio.*` — not inside a called helper, in this body?

If both answers are YES: extract the I/O call into a named helper. The
decision function receives the result as a parameter or return value.

Write in conversation:
```
Uncle Bob scan:
  ParseFile:     decision=yes  I/O=yes  → extracted decodeJSON()
  CompileMatrix: decision=yes  I/O=no   → clean
```

**Output:** Uncle Bob scan written per function; no function in the changeset
contains both decision and direct I/O in the same scope.

---

### Step 8 — Bourgon Logging Scan (KILLER)

For every `slog.*` call in the changeset:

**Message argument** — must be a static string literal:
```go
// PASS
slog.Info("matrix.Compile: parsed plays",
    slog.String("game_id", meta.GameID),
    slog.Int("play_count", len(plays)),
)

// FAIL — message is not a static literal
slog.Info(fmt.Sprintf("compiled %d plays for %s", len(plays), meta.GameID))
```

**Variable values** — every variable in a keyed field: `slog.String`,
`slog.Int`, `slog.Bool`, `slog.Any`.

Any `fmt.Println` or `log.Printf` used for observability: replace with `slog.*`.

Write in conversation:
```
Bourgon scan:
  3 slog calls — all clean
  1 fmt.Println → replaced with slog.Info + keyed fields
```

**Output:** Bourgon scan written; all slog calls use a static message literal
and keyed fields; no fmt.Println/log.Printf for observability.

---

### Step 9 — Compliance Checklist (KILLER)

Write this checklist in conversation with PASS or FAIL for each item:

```
Compliance — coding-standards-go:
  [PASS/FAIL] 1. Pike:      every function named VerbNoun; no mixed-concern body
  [PASS/FAIL] 2. Constants: all repeated/test-referenced literals are package const
  [PASS/FAIL] 3. Fowler:    no >3-line logic block duplicated; extractions first
  [PASS/FAIL] 4. Cheney:    no bare `return err`; every error wrapped with context
  [PASS/FAIL] 5. Uncle Bob: no function mixes decision + direct I/O in same scope
  [PASS/FAIL] 6. Bourgon:   all slog calls use static message + keyed fields
  [PASS/FAIL] 7. Length:    no function body exceeds 30 lines
```

Any FAIL: fix the item, re-run its scan step, update the entry to PASS.
Do not run `git add` until all 7 show PASS.

**Output:** 7/7 PASS checklist written in conversation before any `git add`.

---

## Proof Requirements

- **Step 2:** Pike gate table written; zero FAIL entries before Step 3.
- **Step 3:** Constants gate table written; all promoted literals declared
  as `const` before any function body is written.
- **Step 4:** Fowler scan result written; any extracted function exists
  in the codebase before new logic is written.
- **Step 6:** Cheney scan written; zero bare `return err`; zero `_ = err`.
- **Step 7:** Uncle Bob scan written per function; zero functions with
  both decision and direct I/O in the same scope.
- **Step 8:** Bourgon scan written; all slog calls conform.
- **Step 9:** 7/7 PASS checklist written before `git add`.

---

## State Persistence

N/A — single session. Output is the compliance checklist in conversation
and the committed code.

---

## Red Flags

| Thought | Reality |
|---|---|
| "The name is a bit long but the intent is clear" | A name requiring "and" or "or" is a split diagnosis. Write both names, then write both bodies. |
| "I'll promote the string to a const after the function works" | The const is written before the function body. A literal that survives to commit will survive to the next three files. |
| "I know there's no duplication — I wrote the whole codebase" | The Fowler scan is a grep command, not a memory assertion. Run it. |
| "The error is obvious from context — wrapping adds noise" | The context is obvious now. At 3am, the log line has no call site. Every error gets wrapped. |
| "These two concerns are closely related — splitting would be artificial" | "Closely related" is the justification for every scope-creep function. Pike gate: name it first. If "and" appears, split. |
| "I used fmt.Sprintf in the slog message but the variable is right there" | The message is a static string. The variable goes in a keyed field. This is what makes logs queryable. |
| "The checklist is obviously all passing — I watched the scans run" | The checklist is a written gate, not a mental state. Write it. |

---

## Integration

```
Called by: superpowers:test-driven-development (GREEN phase, before writing
           implementation); standalone before any .go file creation or
           significant modification in paradox-clock-gate
Calls: none
Output: 7-item compliance checklist written in conversation with 7/7 PASS;
        committed .go code that passes all 7 gates
```

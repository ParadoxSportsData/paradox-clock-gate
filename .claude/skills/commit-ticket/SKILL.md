---
name: commit-ticket
description: Per-ticket git commit discipline for clock-gate. TRIGGER when: a TDD phase completes (RED tests written and confirmed failing, or GREEN go test passes), user says "commit this ticket", "commit PDX-N", "make a commit", "time to commit".
---

## Overview

Without this skill, commits accumulate until end-of-session, producing a monolithic diff that loses per-ticket reviewability, TDD write-order evidence, and GitHub-visible progress.

**Iron Law:** One commit per TDD phase, per ticket. Not after "just one more ticket." Not at end of session.

## Announce at Start

"Committing PDX-[N] [RED/GREEN]. Staging files, auditing scope, drafting message, running phase gate, then committing. Visible output at each step."

## Steps

### Step 1 — State Intent (Killer — Incident Commander)

Write exactly:
`"Committing PDX-N [RED|GREEN] now, before any next step."`

Substitute the real ticket number and tag. Do not proceed to Step 2 without writing this line in the conversation.

**Output:** Intent statement written in conversation.

---

### Step 2 — Determine Phase (Killer)

Write one of:
- `"Phase: RED — staging test files only. Will confirm all staged files are *_test.go."`
- `"Phase: GREEN — staging implementation and test files. Will run go test ./... before committing."`

RED = tests just written and confirmed failing.
GREEN = go test ./... just passed.

**Output:** Phase declaration written.

---

### Step 3 — Stage Files

Run `git add` for the files belonging to this ticket. Do not run `git add .` or `git add -A` — name each file explicitly.

**Output:** `git add` command(s) executed.

---

### Step 4 — Audit Staged Files (Killer — Incident Commander)

Run:
```bash
git diff --name-only --cached
```

Write a table for every file:
```
File: [path]   Ticket: PDX-N   Belongs: YES/NO
```

For every file where Belongs = NO: run `git restore --staged [path]`. Re-run the audit until zero NO entries remain before proceeding to Step 5.

**Output:** Audit table written; zero NO entries.

---

### Step 5a — GREEN Phase Gate (Killer)

*Skip to Step 5b if phase is RED.*

Run:
```bash
go test ./...
```

If exit code is non-zero: stop. Fix failing tests, re-run. Do not proceed until exit code is 0. Write `"go test ./... exit code 0"` in the conversation.

**Output:** `go test ./...` exit code 0 written in conversation.

---

### Step 5b — RED Phase Gate (Killer)

*Skip to Step 6 if phase is GREEN.*

Run:
```bash
git diff --name-only --cached | grep -v '_test\.go$'
```

If this produces any output: non-test files are staged. Run `git restore --staged [each non-test file]`. Re-run until the command produces no output.

**Output:** Command produces no output — zero non-test files staged.

---

### Step 6 — Draft Commit Message

Write the full message in the conversation using this template:

```
PDX-N type(scope)[RED|GREEN]: short description (72 chars max)

Why this approach over the alternative — the decision this diff encodes.
Specific named choice: [type / algorithm / constraint / sentinel value chosen].

TDD: [RED: "tests written: TestFoo, TestBar — confirmed failing before this commit"]
     [GREEN: "confirmed RED before implementing — tests: TestFoo, TestBar"]
Implements: PDX-N
[Benchmark: BenchmarkX: 0 allocs/op  — only when this commit proves a benchmark]

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>
```

Valid types: `feat`, `test`, `fix`, `refactor`, `docs`, `chore`
Valid scopes: `ingestion`, `matrix`, `gate`, `presenter`, `cli`, `setup`

**Output:** Full message written in conversation before any git command.

---

### Step 7 — Structural Field Count (Killer — Compliance Auditor)

Read the drafted message. Write a 1 or 0 for each field:

```
Ticket prefix (PDX-\d+ at start of subject):   [1/0]
type(scope) pattern:                            [1/0]
[RED] or [GREEN] tag:                           [1/0]
TDD: line:                                      [1/0]
Implements: PDX- line:                          [1/0]
Co-Authored-By: line:                           [1/0]
─────────────────────────────────────────────────────
Total:                                          [N/6]
```

Any field scoring 0: rewrite that field in the message draft before proceeding. Total must be 6/6 before Step 8.

**Output:** Score of 6/6 written in conversation.

---

### Step 8 — Content Gate (Killer — Linus Torvalds)

Read the body of the drafted message (lines after subject, before `TDD:`).

Write: `"Body names decision: [extract the specific decision stated]"`

The body fails this gate if it restates the diff ("add GameState struct", "implement ParseFile", "write compiler tests") without naming a choice made over an alternative. Rewrite until the decision is explicit.

**Output:** `"Body names decision: [decision]"` written before Step 9.

---

### Step 9 — Commit

Execute:
```bash
git commit -m "$(cat <<'EOF'
[full message from Step 6 with Step 7/8 fixes applied]
EOF
)"
```

**Output:** Commit command executed.

---

### Step 10 — Read the Log (Killer)

Run:
```bash
git log --oneline -1
```

Write the output in the conversation. It must begin with the abbreviated SHA followed by `PDX-N`.

If the line does not start with `PDX-`: the commit was blocked by a hook. Read the hook output, fix the issue, return to Step 9. Do not proceed to the next ticket until this line starts with the correct format.

**Output:** `git log --oneline -1` output written; starts with SHA + `PDX-N`.

---

### Step 11 — Invoke `qa-notes` (Killer)

Write: "Invoking qa-notes for PDX-[N] using git show HEAD as the diff source."

Invoke the `qa-notes` skill. Pass:
- Diff source: `git show HEAD`
- Ticket ID: PDX-[N] from Step 1

The qa-notes skill posts an ADF comment to the Jira ticket. This is the QA documentation record for the commit.

Do not skip this step even if the diff is small. A simple struct ticket produces a minimal qa-notes output — that is correct behavior, not a reason to skip.

**Output:** `qa-notes` invocation complete; ADF comment posted to PDX-[N].

---

## Proof Requirements

- **Step 1:** Intent statement written before any git command.
- **Step 4:** Audit table with zero NO entries before Step 5.
- **Step 5a:** `"go test ./... exit code 0"` written (GREEN only).
- **Step 5b:** grep command produces no output (RED only).
- **Step 7:** 6/6 score written before Step 8.
- **Step 8:** `"Body names decision: [decision]"` written before Step 9.
- **Step 10:** `git log --oneline -1` written; starts with SHA + `PDX-N`.
- **Step 11:** `qa-notes` invocation visible; ADF comment posted to Jira ticket.

## State Persistence

N/A — single session. Output is one git commit per invocation.

## Red Flags

| Thought | Reality |
|---------|---------|
| "I'll commit all the tickets at the end — it's faster" | End-of-session commits are the monolithic history this skill exists to prevent. Commit now. |
| "These files are close enough to this ticket" | "Close enough" is scope bleed. Step 4 requires naming the ticket every file belongs to. If you can't, unstage it. |
| "The body explains what the code does — that's sufficient" | What the code does is the diff. The body names the decision and why this approach over the alternative. Step 8 rejects bodies that restate the diff. |
| "I'll skip go test — I just ran it" | Step 5a runs immediately before committing. State changes between phases are real. |
| "The commit message looks right — I'll skip the field count" | "Looks right" is not 6/6. A missing Co-Author line looks fine at a glance and fails the hook. |
| "The commit ran — I can move on" | Step 10 is not optional. A hook-blocked commit exits non-zero silently. Read the log. |
| "The diff is tiny — qa-notes is overkill" | Step 11 is not optional. A simple diff produces a minimal qa-notes output — that is the correct and proportional result, not a reason to skip. |

## Integration

```
Called by: TDD workflow at each phase boundary; standalone
Calls: qa-notes (Step 11)
Output: One git commit per invocation visible in git log with PDX-N-prefixed conventional commit message;
        QA notes ADF comment posted to the Jira ticket
```

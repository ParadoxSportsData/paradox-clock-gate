---
name: auto-qa
description: Executes QA automatically after a ticket's implementation is committed. TRIGGER when: user or agent says "run auto-qa for PDX-N", "QA PDX-N", "auto-qa PDX-N", commit-ticket GREEN phase completes and QA is the next step, or implementation agent signals handoff after fixing a failure.
---

## Overview

Without this skill, QA results exist only in conversation and are lost at context reset — the ticket has no machine-readable evidence of what was tested or whether acceptance criteria passed. This skill fetches the ticket's DoD command and QA Notes, executes every required command with verbatim evidence capture, and posts a durable ADF comment before transitioning status.

**Executor:** AI agent. Invoked after `commit-ticket` GREEN. Repeated, templated task.

**Iron Law:** No pass comment is posted without chain-of-custody evidence for every command. No failure comment is posted without all three consumer fields (exact command, verbatim output, specific AC violated) present for every failing command.

---

## Announce at Start

"Running auto-qa for [ticket ID]. Fetching ticket and QA Notes, enumerating the full command list, printing each command before it runs, capturing verbatim output, then posting evidence comment and transitioning status. Every command appears in conversation before it executes."

---

## Steps

### Step 1 — Fetch Ticket and Extract DoD Command (KILLER)

Call `mcp__paradox-confluence__jira_get_issue` with the ticket ID, passing `fields=*all,comments` to include comments in the response.

Extract:
- **DoD command**: exact shell command from the "DEFINITION OF DONE" → "Command:" line
- **Expected output**: pattern or string from "Expected output:" line
- **Acceptance Criteria**: all bullet points from the "ACCEPTANCE CRITERIA" section

Write in conversation:
```
Ticket: [ID] — [title]
DoD command: [exact command]
Expected output: [pattern]
AC count: [N] criteria
```

If the DoD section is absent or Command field is blank: write `"BLOCKED — no DoD command in ticket [ID]. Cannot proceed."` Stop.

**Output:** DoD command, expected output, and AC list written before Step 2.

---

### Step 2 — Fetch QA Notes Comment and Parse Test Strategy (KILLER)

From the `jira_get_issue` response comments array, find the comment whose body contains the string `[AI CODING AGENT COMMENT]` and a "Test Strategy" section. If multiple exist, use the most recent.

Extract:
- **Test Strategy commands**: every bullet item that is a runnable shell command (starts with `go`, `./`, `curl`, or is backtick-wrapped)
- **Edge Cases**: all bullet items under "Edge Cases" (descriptive — review during output validation, not direct execution)
- **Regression commands**: every bullet item under "Regression" that is a runnable command

Write in conversation:
```
QA Notes found: yes/no
Test Strategy commands: [N]
  1. [command]
  2. [command]
Edge Cases noted: [N] (descriptive — reviewed during output validation)
Regression commands: [N]
  1. [command]
```

If no QA Notes comment exists: write `"No QA Notes comment found — proceeding with DoD command + baseline health checks only."` Continue to Step 3.

**Output:** Test Strategy commands, Edge Cases, Regression commands extracted and written.

---

### Step 3 — Enumerate Complete Command List (KILLER — Gawande Gate)

Assemble the ordered list:
1. DoD command (from Step 1)
2. All Test Strategy commands (from Step 2), in order
3. All Regression commands (from Step 2), in order
4. `go test ./...` — baseline health check
5. `go vet ./...` — baseline vet

Deduplicate: if any Test Strategy or Regression command is identical to `go test ./...` or `go vet ./...`, keep one instance only.

Write the complete numbered list:
```
Full command list — [N] commands total:
1. [DoD command]           ← Definition of Done
2. [command]               ← Test Strategy
3. go test ./...           ← baseline health
4. go vet ./...            ← baseline vet
```

Write: `"Pre-execution count: [N]. Pass comment requires [N] executed commands with captured output."`

Do not execute any command until this list is written.

**Output:** Complete numbered command list and count written before any Bash call.

---

### Step 4 — Print Each Command Before Executing (KILLER — User Verification)

For each command in the Step 3 list, before running it, write in conversation:
```
Running [N]/[total]: [exact command]
```

Then execute via Bash. Then write the output per Step 5.

Do not batch commands into a single Bash call. Do not execute a command without its "Running N/total" line appearing first.

**Output:** `"Running [N]/[total]: [command]"` written in conversation before each Bash call.

---

### Step 5 — Execute and Capture Verbatim Output (KILLER — Chain-of-Custody)

For each command:
1. Run via Bash
2. Capture complete stdout+stderr — no filtering, no summarizing
3. Record the numeric exit code (not "pass" or "fail" — the integer)

After each execution, write:
```
Command [N] complete.
Exit code: [integer]
Output:
─────────────────────────────
[verbatim stdout+stderr]
─────────────────────────────
```

If output exceeds 100 lines: write first 50 lines + last 20 lines with `[... N lines omitted ...]` marker in conversation; include the full untruncated output in the ADF comment assembled in Steps 7 or 8.

**Output:** Verbatim output and numeric exit code written in conversation for every command.

---

### Step 6 — Determine Pass/Fail

A command passes if exit code = 0. A command fails if exit code ≠ 0.
Overall: PASS if all commands pass. FAIL if any command fails.

Write: `"Result: [PASS/FAIL] — [N] commands run, [K] failed."`

Count commands executed against the pre-execution count from Step 3. If executed count < pre-execution count: write `"COUNT MISMATCH: expected [N], ran [M]."` Run the missing commands before writing the result.

**Output:** Result and count written; executed count = pre-execution count.

---

### Step 7 — PASS Path: Chain-of-Custody Gate → Post Comment → Transition In Testing → Done (KILLER)

*Execute only when overall result = PASS.*

**Chain-of-custody gate:** For each of the N commands, write one line:
```
Command [N] [label]: command ✓ | output ✓ | exit code ✓
```

If any field is absent: collect it from Step 5 output now. Do not call `jira_add_comment` until all N lines show three checkmarks.

**Assemble ADF pass comment** and call `mcp__paradox-confluence__jira_add_comment`:

```json
{"type":"doc","version":1,"content":[
  {"type":"heading","attrs":{"level":2},
   "content":[{"type":"text","text":"auto-qa: PASS — [ticket ID]"}]},
  {"type":"paragraph",
   "content":[{"type":"text","text":"[N] commands run. All passed. [ISO timestamp]"}]},
  // Repeat block below for each command:
  {"type":"heading","attrs":{"level":3},
   "content":[{"type":"text","text":"Command [N]: [label]"}]},
  {"type":"codeBlock","attrs":{"language":"bash"},
   "content":[{"type":"text","text":"$ [exact command]"}]},
  {"type":"codeBlock","attrs":{},
   "content":[{"type":"text","text":"[verbatim output — full, not truncated]"}]},
  {"type":"paragraph",
   "content":[{"type":"text","text":"Exit code: 0"}]}
]}
```

Call `mcp__paradox-confluence__jira_transition_issue` to transition ticket from `"In Testing"` to `"Done"`.

If the transition returns an error: read the error, determine the correct transition name or ID from the response, and retry before concluding.

Write: `"auto-qa PASS — comment posted, ticket [ID] In Testing → Done."`

**Output:** N/N chain-of-custody lines written; ADF pass comment posted; ticket transitioned In Testing → Done; confirmation written.

---

### Step 8 — FAIL Path: Consumer Test Gate → Post Failure Comment → Transition In Testing → In Progress → Handoff (KILLER)

*Execute only when overall result = FAIL.*

**Consumer test gate:** For each failing command, write:
```
Failing command [N]:
  (1) Exact command: ✓ [copy-paste runnable command]
  (2) Verbatim output: ✓ [first line of output]
  (3) AC violated: [exact AC text from Step 1] — ✓ / MISSING
```

If field (3) is MISSING: map the failure to an AC now. If no AC matches, write `"AC violated: no matching AC — possible regression or infrastructure failure"` and include that text in the comment.

**Assemble ADF failure comment** and call `mcp__paradox-confluence__jira_add_comment`:

```json
{"type":"doc","version":1,"content":[
  {"type":"heading","attrs":{"level":2},
   "content":[{"type":"text","text":"auto-qa: FAIL — [ticket ID]"}]},
  {"type":"paragraph",
   "content":[{"type":"text","text":"[N] commands run. [K] failed. [ISO timestamp]"}]},
  {"type":"heading","attrs":{"level":3},
   "content":[{"type":"text","text":"Failed Commands"}]},
  // Repeat for each failing command:
  {"type":"heading","attrs":{"level":4},
   "content":[{"type":"text","text":"Command [N] — FAILED"}]},
  {"type":"codeBlock","attrs":{"language":"bash"},
   "content":[{"type":"text","text":"$ [exact command]"}]},
  {"type":"codeBlock","attrs":{},
   "content":[{"type":"text","text":"[verbatim output — full]"}]},
  {"type":"paragraph",
   "content":[{"type":"text","text":"Exit code: [N]"}]},
  {"type":"paragraph",
   "content":[{"type":"text","text":"AC violated: [exact AC text]"}]},
  // Passing commands summary:
  {"type":"heading","attrs":{"level":3},
   "content":[{"type":"text","text":"Passing Commands ([M] of [N])"}]},
  {"type":"bulletList","content":[
    // one listItem per passing command with label and exit code 0
  ]}
]}
```

Call `mcp__paradox-confluence__jira_transition_issue` to transition ticket from `"In Testing"` to `"In Progress"`.

If the transition returns an error: read the error, determine the correct transition name or ID, retry.

Write in conversation:
```
Implementation flaw detected — handing back. See [ticket ID] for failure details.
Failing commands: [K]
Next: implementation agent reads the failure comment from [ticket ID], fixes the issue,
commits via commit-ticket (which transitions back to In Testing), then invokes auto-qa from Step 1.
```

**Output:** Consumer test gate table written (3 fields confirmed for each failure); ADF failure comment posted; ticket transitioned In Testing → In Progress; handoff message written in conversation.

---

## Proof Requirements

- **Step 1:** DoD command, expected output, and AC list written before any Step 2 action.
- **Step 2:** QA Notes parse result written with Test Strategy count; "No QA Notes found" stated explicitly if absent.
- **Step 3:** Complete numbered command list with count written; `"Pre-execution count: [N]"` present before any Bash call.
- **Step 4:** `"Running [N]/[total]: [command]"` present in conversation log before each corresponding Bash call.
- **Step 5:** Verbatim output and integer exit code written after each command.
- **Step 6:** Executed count = pre-execution count written before result is stated.
- **Step 7 (PASS):** N/N chain-of-custody lines (three checkmarks each) written before `jira_add_comment`; ticket transitioned In Testing → Done; confirmation written.
- **Step 8 (FAIL):** Consumer test table with 3 fields per failing command written before `jira_add_comment`; ticket transitioned In Testing → In Progress; handoff message in conversation.

---

## State Persistence

N/A — single session. Durable output is the ADF comment posted to the Jira ticket. The ticket is the handoff medium; conversation text is not load-bearing across context resets.

---

## Red Flags

| Thought | Reality |
|---|---|
| "The tests looked fine when I ran them earlier — I'll skip the DoD command" | Phantom pass. Every command in the pre-execution list executes and produces captured output. Manual pre-clearance does not exist. |
| "The QA Notes don't seem to have runnable commands" | Parse QA Notes before concluding this. Write "No runnable Test Strategy commands found" explicitly if true. The count must reflect what was found, not assumed. |
| "go test passed, so go vet is probably fine" | Both commands are in the list. Exit codes are captured, not inferred from adjacent results. |
| "The failure is obvious — I'll summarize the output in the comment" | The implementation agent reads the ticket comment, not this conversation. Summarized output loses test names, line numbers, and exact error text — precisely what is needed to locate the fix. |
| "I can see what AC this maps to — I don't need to write field (3) explicitly" | The consumer test gate requires the AC text written in conversation before posting. "I can see it" is not a written gate pass. |
| "I combined two commands into one Bash call to save time" | Each command in the pre-execution list is a separate Bash call with its own output and exit code. Combined calls collapse the per-command evidence. |
| "The transition failed — I'll note it and move on" | A ticket not in Done is still In Progress. Read the transition error, find the correct transition name or ID, retry. |

---

## Integration

```
Called by: commit-ticket (after GREEN phase completes + In Testing transition); standalone; coding-agent (after fixing failure)
Calls: Bash (test execution), mcp__paradox-confluence__jira_get_issue,
       mcp__paradox-confluence__jira_add_comment,
       mcp__paradox-confluence__jira_transition_issue
Expects: ticket in "In Testing" status (set by commit-ticket Step 12)
Output: ADF comment posted to ticket with verbatim evidence per command;
        ticket In Testing → Done (pass) or ticket In Testing → In Progress + handoff message (fail)
Loop: On fail — coding-agent reads failure comment from ticket, fixes implementation,
      commits via commit-ticket (which transitions back to In Testing), invokes auto-qa from Step 1. Repeats until all pass.
```

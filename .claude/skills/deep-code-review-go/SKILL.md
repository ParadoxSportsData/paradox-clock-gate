---
name: deep-code-review-go
description: Exhaustive code review for the paradox-clock-gate Go codebase. TRIGGER when: user says "review this codebase", "deep review", "full review", "code audit", "review all code", "review this PR", "review this ticket", "review this Jira", or provides a GitHub PR URL or Jira ticket ID. Asks for exhaustive-mode confirmation if no artifact is given.
---

# Deep Code Review — Go (paradox-clock-gate)

## Overview

Unfocused reviews sample rather than audit — peripheral files, workers, and utilities go unchecked and systemic issues compound. This skill drives a subagent fleet through every Go file in scope, applies Go-specific quality checks grounded in five expert lenses, and surfaces systemic patterns, not per-file symptoms.

**Iron Law:** Every file in scope is either cited in the report with a verdict, or explicitly listed as "no issues found." No file leaves the inventory silently.

**Executor:** AI agent, invoked by a developer for PR review or full codebase audit. Spawns parallel subagents for large file sets. Completeness over speed.

## Announce at Start

Say this before any other action:

> "Starting deep code review (Go / paradox-clock-gate). [If no PR or Jira artifact given: 'No PR or Jira ticket was provided — do you want an exhaustive review of the entire codebase? If not, paste a PR URL or Jira ticket ID to scope the review.'] I'll build a complete file inventory, dispatch parallel subagents, apply Go-specific checks and five expert-anchor checkpoints, and write a structured markdown report. No file leaves the inventory unexamined."

Wait for scope confirmation before reading any file.

## Steps

### K1 — Determine Review Mode

**1a. PR link given:** Extract PR number; fetch diff via GitHub MCP tools (`mcp__plugin_greenfi-engineering_github-operations__get_pull_request`). Set SCOPE=PR. File list = changed files in the diff only.

**1b. Jira ticket given:** Fetch ticket via Jira MCP tools. Extract any linked PR URLs from description and comments. If a PR is found, proceed as 1a. If no PR is linked, set SCOPE=TICKET; list files mentioned in the ticket description plus files matching the ticket's component/area. State the interpreted file list and ask: "Is this the right scope? Reply yes, or remove any entries." Wait for confirmation.

**1c. Nothing given:** Ask: "No PR or Jira ticket was provided. Should I perform an exhaustive review of the entire codebase? (yes / no — provide a link instead)" Wait. If yes: SCOPE=FULL. If no: wait for a link and return to 1a or 1b.

**Gate:** State `SCOPE={PR|TICKET|FULL}` before reading any file.

---

### K2 — Fetch Scope Artifacts

**SCOPE=PR:** List every filename from the diff. State the count.

**SCOPE=TICKET:** List every derived filename. State the count.

**SCOPE=FULL:**
```bash
find . -name "*.go" -not -path "*/vendor/*" -not -name "*_test.go"
find . -name "*_test.go" -not -path "*/vendor/*"
```
Record production files and test files as separate counts.

**Gate:** State: `"N production files, M test files in scope."`

---

### K3 — Build Complete File Inventory (Gawande FM1+FM2)

Create this table with every file from K2:

| File path | Status |
|-----------|--------|
| ...       | PENDING |

Every file starts as PENDING. No file is added or removed after this step without a logged reason.

**Gate:** Inventory row count == file count from K2. State: `"Inventory: N files queued. 0 may be skipped without a logged reason."`

---

### K4 — Dispatch Parallel Subagents

Group files by package directory. Assign each group to a subagent with:
- Its file list
- The K5 checklist
- The five anchor checkpoint questions (K6–K10)
- Output format: findings as `file:line`, severity (CRITICAL/HIGH/MEDIUM/LOW/INFO)

Collect all subagent results before proceeding to K6.

---

### K5 — Go-Specific Review Checklist

Apply every item below to every file in scope. Record findings with `file:line` and severity. Cross-reference `coding-standards-go` for project-specific conventions.

**Error Handling**
- Every `error` return is handled; no `_ = err` and no silently discarded error return
- Errors at package boundaries are wrapped with context: `fmt.Errorf("doing X: %w", err)`
- No bare `panic()` in library or worker code without an inline comment explaining why recovery is intentional
- Sentinel errors defined as package-level `var`; error strings do not include the calling function's name (that is added by wrapping)

**Concurrency**
- Every goroutine has a defined shutdown path: context cancellation, channel close, or WaitGroup Done
- No goroutine launched without the caller knowing how it terminates
- Every `sync.Mutex` unlock is via `defer`; manual unlock without `defer` requires an explanatory comment
- Buffered vs unbuffered channel choice is justified with a comment when non-obvious

**Interfaces and Design**
- Interfaces are defined at the consumer (the package that depends on the behavior), not the producer
- No interface with more than 5 methods without documented justification; prefer composition
- No package named "util", "helper", "common", or "misc" — all code has a named, purposeful home

**Context Propagation**
- `context.Context` is the first parameter of every function that performs I/O or long computation
- No `context.Background()` called inside a function that already received a context parameter

**Hard-Coded Values — flag every instance**
- String literals appearing more than once (API paths, status strings, config keys, error messages) → named `const`
- Numeric literals other than 0 and 1 appearing in logic → named `const`
- Timeout, retry count, and limit values inline → config struct field or named constant

**Package and Naming**
- All exported types, functions, and methods have godoc comments
- Package names are singular nouns, not verbs or adjectives
- Test files use `_test` package suffix for black-box tests where appropriate

**Testing**
- Multiple-case tests use table-driven format with a `tests []struct{...}` slice
- No `TestXxx` function with more than 3 sequential `if` branches that could be table rows
- Test helpers call `t.Helper()` as their first line
- No test modifies global state without restoring it via `t.Cleanup`

**Resource Management**
- Every `os.File`, `http.Response.Body`, `sql.Rows`, and `net.Conn` is closed via `defer` or in every error return path
- No resource opened inside a loop without being closed within the same loop iteration

---

### K6 — Gawande Gate: Inventory Completion

**6a.** Count inventory rows still PENDING. That number must be 0.

**6b.** For every package directory in scope, at least one of these is true: at least one finding exists with a file in that package, or the statement `"No issues found in {package}"` appears explicitly.

If 6a > 0: name every PENDING file; assign to a subagent immediately; do not write the report until count reaches 0.

If 6b has gaps: list the unchecked packages and review them now.

**Gate:** State: `"Inventory complete: N files reviewed, 0 pending. All M packages have verdicts."`

---

### K7 — Schneier Gate: Boundary and Failure Mode Review

For each item below, the code either has explicit, documented handling (state what it does) or receives a CRITICAL or HIGH finding:

- Every exported HTTP handler or gRPC handler: what does it do on a malformed or missing request body field? What does it do when a downstream service times out or returns an error?
- Every worker or goroutine entrypoint: what does it do when its context is cancelled mid-operation? Does it drain, abort, or leave state partially written?
- Every database call or external service call: what does it do on a transient network error? On a permanent error (connection refused, auth failure)?
- Every value parsed from an external source (JSON body, env var, CLI flag, config file): what does it do on an unexpected type, a missing required field, or an out-of-range value?

---

### K8 — Fowler Gate: Pattern Classification

For every finding involving repeated code, duplicated logic, or hard-coded values:

1. Name the smell: Magic Number, Duplicate Code, Shotgun Surgery, Feature Envy, Data Clumps, Long Method, Primitive Obsession, or equivalent.
2. State the single abstraction that eliminates the pattern: a named const, a helper function, a shared type, a config struct.
3. State how many sites in the codebase contain this pattern — not just the instance that triggered the finding.

A finding that identifies a problem without naming the smell, the fix, and the site count is incomplete — rewrite it before including it in the report.

---

### K9 — Beck Gate: Abstraction Survivability

For every interface, abstraction layer, or generalized helper found during review:

- State the requirement shift (a 20% change in behavior or input shape) that would break this abstraction.
- If the abstraction breaks on any plausible extension: flag HIGH — "brittle abstraction."
- If the abstraction has exactly one concrete implementation with no documented expansion plan: flag MEDIUM — "premature abstraction."
- If the abstraction survives the 20%-shift test: state that explicitly as a passing verdict.

---

### K10 — Kim Gate: Time-to-Safe-Change

For every package reviewed, estimate in minutes how long a Go engineer unfamiliar with this codebase would need to safely make a targeted, single-function change.

Add time for each signal present:
- Missing godoc on an exported symbol: +10 min per missing block
- Untested function: +15 min per function
- Missing package-level doc comment: +15 min
- Global state or package-level side effects: +20 min
- Undocumented concurrency invariants (shared state, lock ordering): +20 min

**Estimate > 60 minutes: flag the package HIGH — "maintainability barrier." List every signal that contributed to the estimate.**

---

### K11 — Produce Structured Report

Write `review-report.md` in the current working directory. The file must contain all five sections below before this gate passes.

**Section 1 — Header:**
```
# Deep Code Review — Go (paradox-clock-gate)
Date: {YYYY-MM-DD}
Scope: {FULL | PR #{number} | Jira {ticket-id}}
Files reviewed: {N}
Packages reviewed: {M}
```

**Section 2 — Executive Summary:**
3–5 bullet points naming the most critical systemic issues found. Each bullet names the issue category, the affected area, and the risk if unaddressed.

**Section 3 — Findings (one block per finding):**
```
### [{SEVERITY}] {Short title}
File: `path/to/file.go:line`
Category: {Error Handling | Concurrency | Design | Hard-Coded Values | Testing | Security | Maintainability}
Smell: {Fowler smell name, or "N/A"}
Finding: {One sentence describing what is wrong}
Impact: {One sentence describing what degrades or fails as a result}
Fix: {Concrete action or short code snippet — never "improve this" or "consider X"}
```

**Section 4 — Package-Level Verdicts:**
```
| Package | Files | Critical | High | Medium | Low | Time-to-Change (min) |
|---------|-------|----------|------|--------|-----|----------------------|
```

**Section 5 — File Inventory:**
```
| File | Status | Finding count |
|------|--------|---------------|
```

**Gate:** State: `"Report written to review-report.md. N findings: X critical, Y high, Z medium, W low."`

---

## Proof Requirements

| Deliverable | Done when |
|-------------|-----------|
| SCOPE set | One of {PR, TICKET, FULL} stated before any file is read |
| Artifacts fetched | File list stated with production + test counts |
| Inventory built | Row count == artifact count; all rows start PENDING |
| K5 applied | Every file has findings or explicit "no issues found" |
| Gawande gate (K6) | 0 PENDING rows; every package has a verdict — stated explicitly |
| Schneier gate (K7) | Every boundary item documented as handled or flagged CRITICAL/HIGH |
| Fowler gate (K8) | Every pattern finding has smell name + abstraction fix + site count |
| Beck gate (K9) | Every abstraction has a survivability verdict |
| Kim gate (K10) | Every package has a time estimate; >60 min packages flagged HIGH |
| Report (K11) | `review-report.md` exists with all five sections; finding count stated |

---

## State Persistence

For large codebases requiring multiple context turns, write to `review-state.json` after K3:
```json
{
  "scope": "FULL",
  "inventory": [{"file": "path/to/file.go", "status": "PENDING"}],
  "findings": [],
  "packages_complete": []
}
```
Update `status` and `findings` after each subagent returns. If context resets, resume from this file rather than restarting.

---

## Red Flags

| Thought | Reality |
|---------|---------|
| "I've reviewed the main files — that covers the important stuff" | K6 requires 0 PENDING files. There is no "important stuff" exemption. |
| "This looks like the same pattern as an earlier finding" | K8 requires counting all sites. One-instance findings are incomplete. |
| "This abstraction is fine — it's a standard Go pattern" | K9 requires a stated 20%-shift failure mode, not an assertion of normalcy. |
| "I'll note this could be a maintainability concern" | K10 requires a minute estimate. A note is not a gate. |
| "The error handling looks reasonable throughout" | K7 requires per-boundary documentation or a finding. "Reasonable" is not documentation. |
| "No PR or ticket was given — I'll just review what's visible" | K1c requires asking the user before choosing scope. |
| "I'll add the inventory to the report after I finish writing findings" | K11 requires all five sections before the gate passes. |
| "The subagents covered most packages — I'll fill in the rest from memory" | Subagent results must be collected before K6. Memory is not a source. |

---

## Integration

Called by: User or any agent starting a review in paradox-clock-gate
Calls: Parallel subagents (one per package group), GitHub MCP tools, Jira MCP tools
References: `coding-standards-go` for project-specific conventions
Output: `review-report.md` in current directory; `review-state.json` if multi-turn

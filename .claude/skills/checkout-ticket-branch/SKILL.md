---
name: checkout-ticket-branch
description: Creates and checks out a PDX-N-short-description branch at the start of ticket work. TRIGGER when: starting work on a Jira ticket, "start ticket PDX-N", "create branch for PDX-N", "begin PDX-N", before invoking test-driven-development for any ticket.
---

## Overview

Without this skill, agents commit to main or an unrelated branch, breaking the PR workflow and losing the PDX-N branch traceability that links GitHub history to Jira tickets.

**Executed by:** AI agent, once per ticket, before TDD or any file change.

**Iron Law:** No implementation work begins until `git branch --show-current` output is written in the conversation and matches `PDX-N-*`.

---

## Announce at Start

"Creating branch for PDX-[N]. Switching to main, pulling latest, deriving branch name, validating format, creating branch. Branch name written in conversation before any work begins."

---

## Steps

### Step 1 — Switch to Main (Killer — Incident Commander)

Run:
```bash
git checkout main
```

If exit code is non-zero: the working tree has uncommitted changes or the branch name `main` doesn't exist. Write the error. Do not proceed until exit code is 0.

**Output:** On `main` branch, clean checkout.

---

### Step 2 — Pull Latest Main (Killer — Incident Commander)

Run:
```bash
git pull
```

If exit code is non-zero: remote is unreachable or there is a merge conflict. Write the error. Do not proceed until exit code is 0 and output contains either `Already up to date` or a fast-forward summary.

**Output:** `git pull` exit code 0 written in conversation.

---

### Step 3 — Check if Branch Already Exists (Killer)

Run:
```bash
git branch --list "PDX-[N]-*"
```

Replace `[N]` with the actual ticket number.

- If output is non-empty: the branch already exists. Run `git checkout [existing-branch-name]`. Write: `"Branch already exists: [name] — switching to it."` Skip to Step 6.
- If output is empty: proceed to Step 4.

**Output:** Branch existence stated explicitly before Step 4.

---

### Step 4 — Derive Branch Name (Killer — Air Traffic Controller)

Construct the branch name:
- Prefix: `PDX-[N]-` where N is the exact Jira ticket number
- Descriptor: 1–4 words from the ticket title, lowercased, spaces replaced with hyphens, non-alphanumeric characters removed
- Full format: `PDX-[N]-[descriptor]`

Examples:
- Ticket "PDX-13: Add README.md with usage examples" → `PDX-13-readme`
- Ticket "PDX-14: Finalize WRITEUP.md assessment narrative" → `PDX-14-writeup`
- Ticket "P2A-1: Implement clock-gate serve subcommand" → `P2A-1-serve-cmd`

Write the proposed name in the conversation before creating anything.

**Output:** Proposed branch name written in conversation.

---

### Step 5 — Validate Name Format (Killer — Air Traffic Controller)

Read the proposed name from Step 4.

The name passes if it satisfies all three:
1. Starts with a ticket prefix (`PDX-`, `P2A-`, `P2B-`) followed by digits and a hyphen
2. Descriptor portion contains only lowercase letters, digits, and hyphens
3. Total length ≤ 50 characters

Write one line:
`"Branch name: [name] — Format: PASS"` or `"Branch name: [name] — Format: FAIL: [reason]"`

If FAIL: rewrite the name to fix the issue before proceeding.

**Output:** `"Format: PASS"` written before `git checkout -b` runs.

---

### Step 6 — Create and Switch to Branch (Killer)

Run:
```bash
git checkout -b [branch-name]
```

If exit code is non-zero: write the error. Do not proceed.

**Output:** Branch created and active.

---

### Step 7 — Confirm Active Branch (Killer — Surgeon)

Run:
```bash
git branch --show-current
```

Write the output verbatim in the conversation.

The output must match the branch name from Steps 4/5. If it does not: stop and do not begin any implementation work until this line is written and matches.

**Output:** Active branch name written in conversation; matches `PDX-N-*`.

---

## Proof Requirements

- **Step 1:** `git checkout main` exit code 0 before pull.
- **Step 2:** `git pull` exit code 0 written before branch creation.
- **Step 3:** Branch existence result written before Step 4.
- **Step 5:** `"Format: PASS"` written before `git checkout -b`.
- **Step 7:** `git branch --show-current` output written; matches branch name from Step 4.

---

## State Persistence

N/A — single session. Output is the active branch confirmed in conversation.

---

## Red Flags

| Thought | Reality |
|---|---|
| "I already know we're on main — no need to checkout" | `git checkout main` is the guarantee. Memory of what branch is active is not. |
| "I'll pull after creating the branch" | Pull must run before branching. A branch from stale main is already behind remote. |
| "The branch name looks right" | "Looks right" is not a format PASS. Step 5 writes the result explicitly before the branch is created. |
| "The branch was created — I can start implementing" | Step 7 is not optional. Read `git branch --show-current` and write the output before proceeding. |
| "I'll skip the branch-exists check — it's probably fine" | `git checkout -b` on an existing branch exits non-zero. The session stalls. Step 3 costs one command and prevents it. |

---

## Integration

```
Called by: standalone at ticket start; before superpowers:test-driven-development
Calls: none
Output: Active branch named PDX-N-short-description confirmed in conversation via
        git branch --show-current; all subsequent commits land on this branch
```

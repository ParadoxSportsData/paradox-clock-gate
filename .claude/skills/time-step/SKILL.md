---
name: time-step
description: Records start/end wall-clock timestamps for SDLC steps to metrics_log.md, computing durations for end events. TRIGGER when: beginning any ticket implementation, doc generation, phase gate, or skill build ("time-step start ticket:MVP1-2"), or ending one ("time-step end ticket:MVP1-2"). Also triggered by the git commit PostToolUse hook for automatic end-of-task recording.
---

## Overview

Agent performance claims without timestamps are anecdotal. This skill records wall-clock measurements at semantic SDLC boundaries — start/end of tickets, docs, phase gates — producing a log file that aggregates into time-per-category data for the assessment story.

**Iron Law:** Every entry must land in `metrics_log.md` and be confirmed by reading back the last line. A measurement that exists only in the conversation is not a measurement.

---

## Announce at Start

"Recording [start|end] for [label]. Writing to metrics_log.md."

---

## Log File Location

```
/Users/athatcher/.claude-work/projects/-Users-athatcher-Documents-at-proj-paradox-clock-gate/memory/metrics_log.md
```

Create if it doesn't exist. Never truncate — always append.

---

## Label Format

All labels must follow: `{category}:{identifier}`

Valid categories:
- `ticket` — implementing a Jira ticket (e.g., `ticket:MVP1-2`)
- `doc` — generating a document (e.g., `doc:prd`, `doc:tech-req`, `doc:design-doc`)
- `phase` — entire implementation phase (e.g., `phase:mvp-phase-1`)
- `gate` — a specific verification gate (e.g., `gate:pre-code-review`)
- `skill` — building a skill via skill-builder (e.g., `skill:create-jira-ticket`)

---

## Steps

### Step 1 — Validate and Normalize Label (Killer)

Read the label provided.

Apply these rules in order:
1. Convert to lowercase
2. Replace spaces with hyphens
3. If no colon is present: ask which category applies (`ticket`, `doc`, `phase`, `gate`, `skill`) before continuing
4. If category is present but not in the valid list: ask for correction before continuing

Write: "Label: [original] → [normalized]" (omit if no change was needed).

**Output:** Normalized label confirmed in conversation.

---

### Step 2 — Branch: Start or End

If event is **start**: proceed to Step 5 (skip Steps 3–4).

If event is **end**: proceed to Step 3.

---

### Step 3 — Match Start Entry (Killer — end events only)

Run:
```bash
grep "| start | {normalized-label}" /Users/athatcher/.claude-work/projects/-Users-athatcher-Documents-at-proj-paradox-clock-gate/memory/metrics_log.md
```

If a matching start entry is found: extract its timestamp. Proceed to Step 4.

If no matching start entry is found:
- Write: "WARNING: no matching start entry for [{label}]. Either start was never recorded or label doesn't match."
- Ask: "Proceed with end entry anyway (duration will be blank), or fix the label?"
- Do not write the end entry until the user responds.

**Output:** Matching start timestamp extracted, or warning written and user consulted.

---

### Step 4 — Compute Duration (Killer — end events only)

Compute the elapsed time between the matched start timestamp and the current time.

Format as: `{N}m {M}s` (e.g., `28m 33s`). If under 60 seconds: `{N}s`. If over 60 minutes: `{N}h {M}m`.

Write: "Duration: {label} took {duration}"

**Output:** Duration computed and written in conversation.

---

### Step 5 — Compose and Append Log Entry (Killer)

Get the current timestamp in ISO-8601 format: `YYYY-MM-DDTHH:MM:SSZ`

Compose the entry:
- **Start:** `{timestamp} | start | {label}`
- **End:** `{timestamp} | end   | {label} | duration={duration}`

Append to `metrics_log.md` using the Write tool (append mode — do not overwrite existing content).

If `metrics_log.md` does not exist: create it with a header line first:
```
# metrics_log — paradox-clock-gate agent performance timestamps
# Format: {ISO-8601} | {start|end} | {category:identifier} | duration={Xm Ys} (end only)
# Aggregation: grep 'end.*ticket:' metrics_log.md | grep -oP 'duration=\K[^\n]+'
```

**Output:** Entry appended.

---

### Step 6 — Confirm Write (Killer)

Run:
```bash
tail -1 /Users/athatcher/.claude-work/projects/-Users-athatcher-Documents-at-proj-paradox-clock-gate/memory/metrics_log.md
```

If the output matches the entry composed in Step 5: write "✓ Written: [{entry}]"

If the output does not match: the append failed. Write the entry again using the Write tool, then re-run the tail check. Do not proceed until the last line matches.

**Output:** "✓ Written: [{entry}]" confirmed in conversation.

---

### Step 7 — Conversation Summary

Write one line:
- **Start:** `⏱ Started: {label} at {timestamp}`
- **End:** `⏱ Completed: {label} — {duration}`

**Output:** One-line summary written.

---

## Aggregation Reference

After multiple entries accumulate, use these commands to aggregate:

```bash
# All completed durations by category
grep '| end ' metrics_log.md

# Time per ticket (all tickets)
grep 'end.*ticket:' metrics_log.md | grep -oP 'duration=\K[^ \n]+'

# Time per doc generation
grep 'end.*doc:' metrics_log.md | grep -oP 'duration=\K[^ \n]+'

# All entries for a specific label
grep 'ticket:MVP1-2' metrics_log.md

# Unmatched starts (no corresponding end)
comm -23 \
  <(grep '| start |' metrics_log.md | grep -oP '\| start \| \K\S+' | sort) \
  <(grep '| end   |' metrics_log.md | grep -oP '\| end   \| \K\S+' | sort)
```

---

## Proof Requirements

- **Step 1:** Normalized label written in conversation before any file operation.
- **Step 3:** grep result written in conversation before duration computed (end events).
- **Step 4:** Duration written in conversation before file append (end events).
- **Step 5:** Entry composed and appended.
- **Step 6:** `tail -1` output matches composed entry — "✓ Written:" line present.

---

## State Persistence

Writes to `metrics_log.md` — persists across sessions. Never truncate.

---

## Red Flags

| Thought | Reality |
|---------|---------|
| "The label is close enough — I'll write it as-is" | Case differences and space-vs-hyphen differences break grep pairing silently. Step 1 normalizes before any write. |
| "I'll match the start entry from memory — I remember it" | grep the file. Memory of what was written earlier in the session is not the file. Session compaction removes it. |
| "The write succeeded — no need to check" | Write tool errors are visible; silent append-to-wrong-location is not. Step 6 tail check is the only confirmation that the entry is where it needs to be. |
| "Duration is nice to have — the timestamps are enough" | Downstream aggregation scripts must recompute duration from every pair if it's not pre-computed. Pre-compute it at write time — once, correctly. |
| "I'll add the aggregation commands later when we need them" | The aggregation reference in the skill file is authoritative. Commands written now, once, correctly — not improvised per session. |

---

## Integration

```
Called by: agent at SDLC task boundaries (start/end of tickets, docs, phases, gates);
           PostToolUse git commit hook (auto-records end-of-task timestamps)
Calls: none
Output: Entry appended to metrics_log.md; duration written to conversation on end events
```

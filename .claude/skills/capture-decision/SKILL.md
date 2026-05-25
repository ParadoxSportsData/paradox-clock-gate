---
name: capture-decision
description: Real-time context sync — invoke IMMEDIATELY when any significant decision is made, scope changes, an AI suggestion is accepted or overridden, a guardrail is added, or a flaw is identified. TRIGGER when: a design choice is confirmed, a tool or approach is accepted or rejected, scope expands or contracts, a hook or constraint is added, a known limitation is acknowledged, Aaron approves or overrides anything Claude proposed. One response turn = invoke now, not at session end.
---

## Overview

Decisions made during a build evaporate if not captured immediately. This skill is the connective tissue between the moment a decision happens and every file that needs to reflect it.

**Iron Law:** Invoke in the same response turn as the decision. Not at session end. Not after the next step. Now.

---

## Announce at Start

"Invoking capture-decision for: [one-sentence description of what was just decided]. Working through the 7-item checklist — writing to every applicable file before this response ends."

---

## Steps

### Step 1 — Timing Gate (Killer)

State in this response: "This is the turn in which [decision] was made."

If the decision was made in a prior exchange and you are now in a subsequent response, write: "Late capture: [reason for lateness]" and proceed anyway. Do not skip capture because of lateness.

**Output:** Timing statement written before proceeding.

---

### Step 2 — Intelligence Analyst: State What Was Decided (Killer)

Write explicitly in this response:

- **What was decided:** One specific sentence. Not "we discussed architecture" — "we chose `uint16` over `float64` for WinProb to keep GameState GC-free."
- **What confirms this decision:** The evidence or rationale that makes it locked.
- **What is uncertain:** Any aspect still open or conditional.

If you cannot write a specific one-sentence decision statement, the decision is not yet made. Surface the ambiguity — do not proceed.

**Output:** Decision statement, rationale, and uncertainties written before any file writes.

---

### Step 3 — Memory Files: Route and Write (Killer)

Route to the correct memory file based on the decision type:

- Design or architecture choice → append to `~/.claude-work/projects/-Users-athatcher-Documents-at-proj-paradox-clock-gate/memory/design_decisions.md`
- AI suggestion accepted or overridden → append to `~/.claude-work/projects/-Users-athatcher-Documents-at-proj-paradox-clock-gate/memory/ai_moments.md`
- Session progress or benchmark numbers → `implementation_log.md`

Write now. Do not defer.

**Output:** At least one memory file written before proceeding.

---

### Step 4 — WRITEUP.md: Update Assessment Narrative (Killer)

Open `/Users/athatcher/Documents/at-proj/paradox-clock-gate/WRITEUP.md`.

Map the decision to the section that matches:
- AI override or acceptance → §3 "Where AI Got It Wrong / Where I Pushed Back"
- Agent guardrail added or modified → §2 "How I Set Up the Agent to Work Safely"
- Known limitation or rough edge → §4 "Honest Rough Edges"
- Future improvement identified → §5 "What I'd Add With More Time"
- Broader process or workflow decision → §1 "How AI Was Used"

If a section applies: append the specific update text to it now.
If none apply: write "No WRITEUP update needed — [reason]" and proceed.

**Output:** WRITEUP section written, or explicit "no update" statement.

---

### Step 5 — Interview Prep: Add Discussable Q&A

Open `~/.claude-work/projects/-Users-athatcher-Documents-at-proj-paradox-clock-gate/memory/interview_prep.md`.

Ask: would a Prelude Origin team interviewer probe this decision in a technical follow-up?

If yes: add a Q&A entry in this format:
```
- *"[Natural interviewer question]"* → [Specific, confident answer with rationale]
```

Place under the section header that matches the topic:
- Architecture / data representation / performance → Architecture & Design
- AI suggestion accepted or overridden → AI Collaboration
- Hook, allowlist, scope guard, secrets scan → Agent Access & Guardrails
- Known limitation, design trade-off → Trade-offs & Rough Edges
- Concurrent users, cold-start, memory → Scalability & Production Architecture
- Improvement roadmap item → Continuous Improvement
- Plugin repo / AI workflow system decisions → Agent Harness & Process Documentation
- Assessment fit, paradox-platform rewrite → Broader Context

If the topic already has a Q&A entry: update it rather than duplicating.
If not discussable: write "No interview prep update — [reason]" and proceed.

**Output:** New or updated Q&A entry, or explicit "no update" statement.

---

### Step 6 — Plan File: Update if Scope or Decisions Changed

Open `~/.claude-work/plans/now-do-you-understand-rustling-sky.md`.

Update if any of these are true:
1. A locked decision in the decisions table changed.
2. Scope expanded or contracted (new phase, removed feature, timeline change).
3. The setup sequence changed (new prerequisite, step reordered).

If any trigger fires: update the specific section now.
If none fire: write "No plan update — [reason]" and proceed.

**Output:** Plan section updated, or explicit "no update" statement.

---

### Step 7 — Plugin Repo Content: Classify the Decision

Ask: does this decision (skill, hook, memory structure, settings pattern) belong in the plugin repo install manifest — i.e., would it be useful to any AI-assisted project by default?

Answer yes if:
- It solved a problem every project faces (context loss, agent scope creep, constraint enforcement, secrets hygiene)
- It is not specific to NFL data, Go's allocation model, or the Prelude assessment

If yes: append to the Plugin Repo Content table in `~/.claude-work/plans/now-do-you-understand-rustling-sky.md` under `## Plugin Repo Content`:

```
| [Component name] | [What it is] | [Why it exists] | [Universal / Go-specific / Conditional] |
```

If no: write "Not a plugin repo item — [reason]" and proceed.

**Output:** New Plugin Repo Content row, or explicit "not applicable" statement.

---

### Step 8 — Continuous Improvement: Log Flaws and Trade-offs

Open `~/.claude-work/projects/-Users-athatcher-Documents-at-proj-paradox-clock-gate/memory/continuous_improvement.md`.

Update if any of these are true:
1. A trade-off was explicitly made (gained X, gave up Y).
2. A known limitation was identified (what breaks, under what conditions).
3. A future improvement path was named (what the fix would be, why it's out of scope now).

If any trigger fires: append under the correct section (clock-gate or AI Workflow System) and category (UX gap, architectural trade-off, or enforcement gap):

```
- **[Topic]:** What the gap is. What it costs. What the fix would be. Why it was out of scope.
```

If none fire: write "No continuous improvement update — [reason]" and proceed.

**Output:** New entry written, or explicit "no update" statement.

---

### Step 9 — README: Invoke readme-sync if Interface Changed

Ask: did this decision change anything in the public interface — CLI flags, subcommands, output format, architecture overview, or install instructions?

If yes: invoke `readme-sync` skill now. Do not defer.
If no: write "No README update needed — [reason]" and proceed.

**Output:** `readme-sync` invoked and output presented, or explicit "no update" statement.

---

### Step 10 — Incident Commander: Zero-Context Agent Test (Killer)

Re-read every file written in Steps 3–9.

For each written entry, ask: could a new agent reading only this file at 2am, with zero conversation history, answer all of these?
- What was decided?
- Why was it decided?
- What was rejected and why?
- What are the known limitations?
- What comes next?

If any entry fails: rewrite it now. Do not proceed until all written content passes this test.

**Output:** Confirmation that all updates pass, or rewritten entries.

---

### Step 11 — Declare Complete

State: "capture-decision complete for: [decision statement from Step 2]. Updated: [list every file written]."

---

## Proof Requirements

- **Step 2:** A specific one-sentence decision statement exists in this response before any file writes.
- **Step 3:** At least one memory file has a new append, or an explicit "no applicable memory file — [reason]" statement exists.
- **Steps 4–9:** Either the file is updated with specific new text, or an explicit "no update needed — [reason]" statement exists. Silence does not satisfy this requirement.
- **Step 10:** Every update written in Steps 3–9 has been re-read and passes the 2am specificity test.
- **Step 11:** The completion declaration names the decision and lists all files written.

---

## State Persistence

N/A — single session, all outputs are file writes, no state.json required.

---

## Red Flags

| Thought | Reality |
|---------|---------|
| "I'll capture this at the end of the session" | End-of-session batch writes don't distinguish mid-session decisions from retrospective reconstruction. Capture now. |
| "This decision is obvious — no need to write it down" | Obvious now. Not obvious to a fresh agent next session with no context. |
| "The decision isn't fully settled yet" | Capture the partial state. Note what's uncertain. Do not wait for certainty. |
| "I already know which files need updating — I'll skip the routing step" | Routing is the step that decides which files get updated. A skipped route means wrong file or no file. |
| "I'll write a quick summary instead of specific entries" | Quick summaries fail the 2am test. Specific entries survive session handoffs. |
| "readme-sync can wait until the phase is done" | A stale README at submission is a credibility gap. Invoke it now. |
| "This is a plugin repo architectural decision, so no WRITEUP update needed" | If it affects what a Prelude reviewer reading the WRITEUP would want to know about the AI workflow system, WRITEUP §1 or §2 still applies. |

---

## Integration

```
Called by: Claude agent immediately upon any significant decision, scope change, AI moment, or flaw identification
Calls: readme-sync (Step 9, when interface changed)
Output: Updated memory files (design_decisions.md, ai_moments.md, implementation_log.md),
        WRITEUP.md sections, interview_prep.md Q&A, plan file sections,
        plugin repo content rows, continuous_improvement.md entries
```

---
name: skill-builder
description: Meta-skill for creating or improving any skill. Enforces Gawande's checklist methodology — identify killer items, apply the 8-component template, pass 3 quality tests before writing a single line of the new skill. TRIGGER when: user says "create a skill", "write a skill", "build a skill", "make a skill", "add a skill", "new skill", "create a command", "new command", "write a command", "improve a skill", "update a skill", "fix a skill", "analyze gaps in a skill", "make this skill better", "what's missing from this skill", "review this skill", or when any skill creation or improvement is requested as part of a larger workflow.
---

# Skill Builder

## Overview

Skills are READ-DO checklists for experts. They exist to prevent **errors of omission** — steps skipped under pressure, not from lack of knowledge. A good skill focuses on **killer items**: the 5–9 steps that, if skipped, cause failure. Everything else is noise.

**Iron Law:** Do not write a single line of the new or updated skill until Phases 1 and 2 are complete and their outputs are written in the conversation for the engineer to see.

**Announce at start:** "Using skill-builder to design `{skill name}`. I'll work through five phases before writing anything: (1) define the outcome, (1b) propose a domain anchor menu and let you choose the professional lens, then map selected anchors to failure modes, (2) identify killer items, (3) write the skill, (4) apply three quality tests, (5) run a compliance test to confirm agents follow the skill. Written output appears at every phase gate before proceeding."

---

## Mode: New vs. Improve

**Before Phase 1, determine the mode:**

- **New skill:** No existing SKILL.md. Proceed through all phases from scratch.
- **Improve existing skill:** Read the existing SKILL.md in full before answering Phase 1 questions. Note which killer items are currently missing, weak, or lacking verifiable output. Frame Phase 1 answers in terms of what the current skill gets wrong or omits.

State the mode explicitly: `"Mode: [New / Improving existing skill at {path}]"`

---

## Phase 1: Define the Outcome

Answer all 5 questions. Ask the engineer for any not answered by their description.

**Write the answers as a numbered list in the conversation before proceeding. Do not continue until the engineer has seen them.**

1. **What problem does this skill solve?** One sentence — not "it helps with X," but what specifically fails without it?
2. **What does perfect execution look like?** Describe the end state an engineer can inspect and verify.
3. **What are the 3 most common failure modes for this type of skill?** Failure modes, not success states.
4. **Who executes this skill, under what conditions?** Rushed AI agent? Human in an incident? Novel or repeated task?
5. **Where does it live in a larger workflow?** What calls it? What does it call? Can it run standalone?

> **PAUSE — Post the 5 answers as written text in the conversation. The engineer must be able to read and challenge them. Do not proceed until this output exists.**

---

## Phase 1b: Establish Expert Anchors

Anchors sharpen the skill beyond structure — they embed a quality lens that forces the executing agent to see what checklist steps alone cannot catch. Phase 1b runs in two parts: first choose the professional lens, then map it to failure modes.

**Rules:**
1. One anchor per distinct failure mode. Maximum 3 anchors.
2. Each anchor must address a failure mode the others do not. If two anchors address the same root cause, one is redundant — drop it.
3. Anchors are embedded as operational checkpoints inside the skill's reasoning flow. They are NOT header text or introductory flavor. They must block progress until the test passes.

---

**Step 1b.0 — Generate domain anchor menu:**

From the skill's domain (inferred from Phase 1 answers), identify 3–5 professional archetypes whose quality lens would shape what the skill optimizes for. Present a numbered menu — each entry names the expert or archetype, their core principle in one sentence, and the single focus question they would ask that others wouldn't.

Format:
```
Domain: {skill's professional domain — e.g., "engineering communication", "security review", "technical writing"}

Anchor options:
1. {Expert or archetype name} — {one-sentence principle}
   Focus question: {the single question this expert asks that a checklist alone wouldn't catch}

2. {Expert or archetype name} — {one-sentence principle}
   Focus question: ...

3. ...
```

Examples of expert/archetype anchors by domain (use as inspiration, not a fixed list):
- **Communication / writing:** Chip Heath (is every claim specific enough to act on?), journalist (would this survive a fact-check?), product manager (does a non-builder understand what changed and why?)
- **Security:** Bruce Schneier (what does an attacker see, not what does a defender intend?), compliance auditor (is every control traceable to a specific requirement?)
- **Engineering / code:** Atul Gawande (what are the killer items — the steps that cause failure if skipped?), Kent Beck (is the design reversible, or does it lock in a decision?)
- **Data / database:** query planner (will this degrade at 10× volume?), data steward (is PII handled correctly at every surface?)
- **Operations / incident response:** on-call engineer (what does someone need at 2am to diagnose this without context?)
- **Product / UX:** user advocate (what does this look like from someone encountering it for the first time?), accessibility auditor (does this work for every user, not just the happy path?)

Ask the engineer:
> "Which of these anchors should shape this skill's quality checkpoints? Select 1–3. You can also name a different expert or archetype not in the list."

Wait for the engineer's selection. Carry chosen anchors forward into Step 1b.1.

> **PAUSE — Anchor menu must be posted and selection confirmed before mapping to failure modes.**

---

**Step 1b.1 — Map selected anchors to failure modes:**

For each failure mode from Phase 1 Q3, map it to the most appropriate of the selected anchors from Step 1b.0. If a selected anchor doesn't map to a specific failure mode but adds a valuable general lens, embed it as a cross-cutting checkpoint and note where.

```
Failure mode: {from Phase 1 Q3}
Mapped anchor: {selected expert name + principle}
Checkpoint test: {the concrete, runnable question the agent asks itself}
Embedded at: {which step in the skill this checkpoint lives inside}
```

**Step 1b.2 — Overlap check:**

Read all proposed checkpoints. If two test the same underlying condition, drop the weaker one. State which was dropped and why.

**Step 1b.3 — Engineer approval:**

Present the final anchor table:

```
| Failure mode | Anchor | Checkpoint test | Embedded at |
|-------------|--------|-----------------|-------------|
```

> "These anchors will be embedded as checkpoints in the skill. Do you want to keep, replace, or add any anchor before I write the skill?"

Wait for response. Accept, modify, or proceed with confirmed anchors.

> **PAUSE — Anchor table must be posted and approved. Do not write the skill until anchors are confirmed.**

---

## Phase 2: Identify Killer Items

A killer item is any step where skipping causes: wrong output, missed gap, security hole, data loss, or unrecoverable state.

1. List every step the skill must perform.
2. For each step ask: *"If this step is skipped, what specifically fails?"*
3. Mark each **KILLER** or **ROUTINE**.
4. Count killer items. **Target: 5–9.** If more than 9: the skill is doing too much — identify the split point.

Routine items can be grouped or implied. Only killer items become explicit checkpoints in the final skill.

**Post the full table in the conversation before proceeding:**

```
| Step | KILLER / ROUTINE | What fails if skipped |
|------|------------------|-----------------------|
| ...  | KILLER           | ...                   |
```

> **PAUSE — The killer item table must be posted and visible in the conversation. Count must be ≤ 9. Do not proceed until the engineer has seen the table.**

---

## Phase 3: Write the Skill

Every skill must contain all 8 components below. Do not omit any.

**Component 1 — Frontmatter**
```yaml
---
name: {kebab-case-name}
description: {One sentence: what it does AND when to invoke it. Must include: "TRIGGER when: user says X, Y, Z" — list every real-world phrasing including improvement variants}
---
```
The TRIGGER list must catch every real-world phrasing. A vague trigger means the skill never fires.

**Description anti-pattern — never summarize workflow in the description:**
Agents read the description to decide whether to load a skill. If the description summarizes the workflow, agents follow the description *instead of reading the skill body* — the full skill becomes documentation they skip.

❌ `description: Use when creating skills — defines outcome, identifies killer items, applies 8-component template, runs 3 quality tests`
✅ `description: Meta-skill for creating or improving any skill. TRIGGER when: user says "create a skill", "improve a skill"...`

The description must state *when* to invoke, not *what the skill does*.

**Component 2 — Overview (≤ 5 lines)**
- What problem this solves (one sentence)
- Who/what executes it
- The **Iron Law** — one non-negotiable rule, stated once, clearly

**Component 3 — Announce at Start**
The exact first words the agent says when invoked. Must tell the engineer: what is about to happen, what input is needed, and that written output will appear at each phase gate.

**Component 4 — Steps (READ-DO format)**
- Each step is a direct instruction, not a description of what to think about
- Each step produces a **verifiable output** (a file, a list, a confirmed fact, a decision)
- Killer items from Phase 2 are **explicit checkpoints** — numbered, unavoidable
- Phase gates are marked: `> **PAUSE — [exact condition that must be true before continuing]**`

**Component 5 — Proof Requirements**
For every major deliverable: the objective criterion for "done."
- Fail: "ensure the output is complete"
- Pass: "the output file exists at `{path}` and contains all 8 required components"

**Component 6 — State Persistence** *(include only if skill spans multiple context turns)*
What writes to state.json, after which step, and in what format. If single-session, write: `N/A — single session, no persistence required.`

**Component 7 — Red Flags Table**
The 4–6 rationalizations that tempt agents to skip or shortcut this skill:
```
| Thought | Reality |
|---------|---------|
```

**Component 8 — Integration**
```
Called by: {skill name or "standalone"}
Calls: {sub-skills, or "none"}
Output: {exactly what this skill produces — file path, state.json key, report, etc.}
```

> **PAUSE — All 8 components must be drafted before applying the quality tests.**

---

## Phase 4: Apply the 3 Quality Tests

Run all 3. Do not finalize until all 3 pass.

**Test 1 — Rickover Test**

Read the skill top-to-bottom as if you have never seen it and have zero prior context.

Do this literally: scan every sentence and write down each instance of:
`"handle appropriately"` / `"as needed"` / `"if necessary"` / `"where relevant"` / `"appropriately"` / `"properly"`

For each instance found: replace with a specific, observable criterion before proceeding.

Passes when: the list of ambiguous phrases is empty — zero instances remain.

**Test 2 — Feynman Test**

Scan the entire skill for every instance of: *ensure, verify, confirm, check, review, validate, make sure.*

For each instance: replace with a step that produces a concrete, checkable result (a count, a file, a list, a yes/no decision with explicit criteria).

Passes when: grep for `ensure\|verify\|confirm\|check\|review\|validate\|make sure` returns 0 matches outside of this quality-test section itself.

**Test 3 — Gawande Trim**

Remove every line that fails: *"if an expert agent skipped this line, would execution fail or produce a worse outcome?"*

Cut without mercy:
- Sentences restating what was already said
- Explanations of *why* a step exists, unless the agent needs the why to make a judgment call
- Steps fully covered by an adjacent step
- Anything an expert would do automatically without being told

**Length gates:**
- Single-phase skill: target 150–300 lines. Over 400 = re-trim.
- Multi-phase skill: target 400–600 lines. Over 700 = consider splitting into sub-skills.

Passes when: the skill is as short as it can be while remaining complete.

---

## Phase 5: Compliance Test

Document quality (Phases 1–4) does not prove agent behavior. Run one compliance scenario before writing the file.

**Step 5.1 — Baseline:**
Describe a realistic invocation scenario matching this skill's trigger. In the conversation, simulate or subagent the task *without* the skill present. Write down exactly how the agent shortcuts, rationalizes, or produces wrong output.

**Step 5.2 — Compliance:**
Provide the same scenario *with* the written skill draft. The skill passes if agent output matches the perfect execution state from Phase 1 Q2.

**Step 5.3 — Close loopholes:**
If the agent found a new rationalization in Step 5.2, add it to the Red Flags table and repeat from Step 5.2. Done when Step 5.2 produces compliant output with no new rationalizations.

> **PAUSE — Baseline failure (Step 5.1) and compliance confirmation (Step 5.2) must be documented in the conversation before writing the file.**

---

## Write the File

After all 3 quality tests pass and Phase 5 compliance is documented, write to:

```
.claude/skills/{name}/SKILL.md
```

(Path is relative to the project root — `/Users/athatcher/Documents/at-proj/paradox-clock-gate/.claude/skills/{name}/SKILL.md`)

**Proof of completion:**
1. File exists at `.claude/skills/{name}/SKILL.md`.
2. File contains all 8 required components (frontmatter, overview, announce, steps, proof requirements, state persistence, red flags, integration).
3. The Rickover Test passes: zero ambiguous phrases remain in the written file.
4. The Feynman Test passes: zero unresolved `ensure/verify/confirm/check/review/validate` instances in the written file.
5. Phase 5 compliance test is documented in the conversation: baseline failure (Step 5.1) and compliance confirmation (Step 5.2) both visible.

---

## Red Flags

| Thought | Reality |
|---------|---------|
| "I know what a good skill looks like — I'll just write it" | Skipping intake produces skills that look complete but miss killer items |
| "The steps are obvious — no need to list them" | Obvious to whom? Rickover Test: a zero-context agent must execute reliably |
| "This is getting long — I'll document everything to be safe" | Length defeats execution. Killer items only. Cut the rest. |
| "I'll add a note to verify this later" | Verification is not a note. It is a testable criterion embedded in the step. |
| "Phase gates slow things down" | The failures gates prevent are exactly what happens when gates are skipped |
| "I answered the questions in my head — good enough" | Phase 1 and 2 output must be written text the engineer can read and challenge. Unwritten answers cannot be reviewed. |
| "The skill is basically right, just needs a tweak" | Improvement requests require the same Phase 1–2 intake as new skills. Tweaking without intake produces the same errors of omission. |
| "The quality tests passed — no need to test with an agent" | Rickover/Feynman/Gawande test the document, not agent behavior. A clean document can still fail to change what agents do. Phase 5 is not optional. |
| "I'll put the workflow summary in the description so agents know what's coming" | Agents follow the description instead of reading the skill body. Description = when to invoke only, never what the skill does. |

## Integration

**Called by:** Any skill creation or improvement request (auto-triggered)
**Calls:** none
**Output:** `.claude/skills/{name}/SKILL.md` written to project, passing all 3 quality tests, compliance test documented in conversation

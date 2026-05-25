---
name: artifact-context
description: Produces a structured metadata block for embedding in any project artifact — assessment context, GitHub org, contact, current phase, and cross-references to all related documents. TRIGGER when: any artifact-generating skill (create-prd, create-tech-requirements, create-design-doc, create-jira-ticket, task-breakdown) is about to generate output and needs canonical project metadata.
---

## Overview

Artifacts generated without consistent metadata become islands — inconsistent headers, missing cross-references, stale phase markers. This skill reads the authoritative source files and produces a single metadata block all artifact skills embed.

**Iron Law:** Every field in the output block must name its source file. No field is written from memory.

---

## Announce at Start

"Invoking artifact-context. Reading project_context.md and plan file now — producing metadata block for embedding."

---

## Steps

### Step 1 — Read project_context.md (Killer)

Read this file in full using the Read tool:
`~/.claude-work/projects/-Users-athatcher-Documents-at-proj-paradox-clock-gate/memory/project_context.md`

Extract:
- Assessment role and company
- Contact email
- GitHub org and repo name
- Submission format
- Related project (paradox-platform)

Do not use session memory for any of these values. If the file does not exist: stop and write "ERROR: project_context.md not found — artifact-context cannot run without it."

**Output:** Extracted values written before assembly.

---

### Step 2 — Read Current Phase from Plan File (Killer)

Read the plan file using the Read tool:
`~/.claude-work/plans/now-do-you-understand-rustling-sky.md`

Find the section headed "## Pre-Code Work Sequence" or "## Key Decisions Log" — look for explicit DONE / IN PROGRESS markers to determine the current implementation state.

Write: "Current phase sourced from plan file section: [section name]. Value: [Phase N — status]."

Do not use session memory for the phase value. The plan file is the authoritative source.

**Output:** Current phase with source citation written before assembly.

---

### Step 3 — Assemble Metadata Block (Killer)

Produce the following block. Append `[source: filename]` to every field:

```
## Project Metadata

Assessment: Software Engineer / Forward Deployed Engineer — Prelude, Origin team [source: project_context.md]
Contact: claudia@preludesecurity.com [source: project_context.md]
GitHub: ParadoxSportsData/paradox-clock-gate [source: project_context.md]
Submission: Prelude assessment portal — GitHub link [source: project_context.md]
Current phase: [value from Step 2] [source: plan file]

## Document Cross-References

Confluence container page: https://paradox-platform.atlassian.net/wiki/spaces/PDX/pages/21397505/paradox-clock-gate (ID: 21397505) [source: created 2026-05-24]
PRD: https://paradox-platform.atlassian.net/wiki/spaces/PDX/pages/21364737/paradox-clock-gate+Product+Requirements+Document [source: created 2026-05-24]
Technical Requirements: https://paradox-platform.atlassian.net/wiki/spaces/PDX/pages/21495810/paradox-clock-gate+Technical+Requirements [source: created 2026-05-24]
Architecture Design Doc: https://paradox-platform.atlassian.net/wiki/spaces/PDX/pages/21528578/paradox-clock-gate+Architecture+Design+Document [source: created 2026-05-24]
Task Tracker: https://paradox-platform.atlassian.net/wiki/spaces/PDX/pages/21594113/paradox-clock-gate+Task+Tracker+amp+Milestone+Doc [source: created 2026-05-24]
GitHub repo: https://github.com/ParadoxSportsData/paradox-clock-gate [source: project_context.md]

## Confluence

Space: PDX (paradox-platform.atlassian.net)
Parent page ID for all new docs: 21397505
All paradox-clock-gate documents nest under the container page — pass parent_id: 21397505 to confluence-page.

## Related

Source data: ../paradox-platform/data/raw/ (270 NFL game JSON files, 2011 season) [source: project_context.md]
Do not modify: ../paradox-platform [source: project_context.md]
```

Any field whose URL has not yet been provided in this session: write "Not yet created" — do not omit the field.

**Output:** Complete metadata block written with source labels on every field.

---

### Step 4 — Librarian Check (Killer)

Read the Document Cross-References section of the block from Step 3.

Every slot must have either a URL or "Not yet created." A blank or absent slot is a gap.

Write: "[N] document(s) not yet created: [list]. These will need URLs added when the documents are generated."

If all slots have URLs: write "All document cross-references populated."

**Output:** Librarian check statement written.

---

### Step 5 — Output for Calling Skill (Killer)

Write: "artifact-context complete. Metadata block above is ready for embedding. Calling skill: copy the block into the artifact's header or preamble section."

Do not summarize or abbreviate the block. The calling skill uses the full block from Step 3.

**Output:** Completion statement written; full block available for calling skill.

---

## Proof Requirements

- **Step 1:** Read tool call targeting `project_context.md` visible in conversation — no memory substitution.
- **Step 2:** Read tool call targeting the plan file visible in conversation; current phase stated with section citation.
- **Step 3:** Every field in the metadata block has a `[source: filename]` label — zero unlabeled fields.
- **Step 4:** Every Document Cross-References slot has a URL or "Not yet created" — no blank or absent slots.
- **Step 5:** Full metadata block available in conversation for the calling skill; no summarization.

---

## State Persistence

N/A — single session. Output is the metadata block written in the conversation.

---

## Red Flags

| Thought | Reality |
|---------|---------|
| "I know the GitHub org — no need to read the file" | Config manager: every field names its source. Memory is not a source. Read the file. |
| "The phase is Phase 2 — I remember from this session" | Auditor: the plan file is the authoritative source for phase status. Session memory drifts. Read it. |
| "There's no PRD yet so I'll just skip that slot" | Librarian: a missing slot is a gap. Write "Not yet created" so a reviewer knows the document is expected but absent. |
| "I'll add source labels after I assemble the block" | Labels must be written during assembly. Post-hoc labeling means the value was already written from memory. |
| "The calling skill already has this metadata from earlier in the session" | Each invocation is a live read. Metadata recalled from a prior invocation or session context is stale. |

---

## Integration

```
Called by: create-prd, create-tech-requirements, create-design-doc, create-jira-ticket, task-breakdown
Calls: none
Output: Structured metadata block written in conversation, ready for calling skill to embed in artifact
```

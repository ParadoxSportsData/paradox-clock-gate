---
name: create-jira-ticket
description: Creates a single Jira ticket in the personal workspace with all 10 template sections populated and blocker relationships set. TRIGGER when: user says "create a Jira ticket", "create the ticket for [phase/component]", "add this to Jira", "open a ticket for", or task-breakdown needs to create implementation tickets.
---

## Overview

Jira tickets created without subskills produce free-form bodies that a follow-on agent cannot execute — vague algorithm steps, "tests pass" as acceptance criteria, blank cross-reference links. This skill enforces subskill invocation and validates output before any MCP call.

**Iron Law:** `artifact-context` and `jira-ticket-template` must both be invoked and their outputs validated before `jira_create_issue` is called. There is no shortcut.

---

## Announce at Start

"Creating Jira ticket for: [title]. Invoking artifact-context for metadata, jira-ticket-template for the 10-section body, then validating both before submitting. Output visible in conversation at each step."

---

## Inputs Required

Collect from user or calling skill before proceeding:
- **Title intent:** what the ticket is implementing (e.g., "internal/ingestion/parser.go — ParseFile()")
- **Phase:** implementation phase (e.g., MVP Phase 1, Setup)
- **Component:** ingestion | matrix | gate | presenter | cli | setup | doc
- **Depends on:** list of ticket IDs this ticket is blocked by (can be empty)
- **Blocks:** list of ticket IDs this ticket blocks (can be empty)
- **Algorithm intent:** enough detail to write each numbered algorithm step naming an exact Go construct
- **Acceptance condition:** the verifiable end state (fed into AC and DoD sections)

If any required input is missing: ask before proceeding. Do not invent inputs.

---

## Steps

### Step 1 — Invoke `artifact-context` (Killer)

Invoke the `artifact-context` skill. Pass: current phase, component, GitHub org (`ParadoxSportsData/paradox-clock-gate`).

Write the returned metadata block in the conversation.

**Output:** Artifact-context metadata block written in conversation.

---

### Step 2 — Validate Artifact-Context Output (Killer)

Read the metadata block from Step 1. For each field in Document Cross-References:
- If the field has a URL: write "✓ [field name]: [url]"
- If the field says "Not yet created": write "⚠ [field name]: Not yet created — will appear as placeholder in ticket"
- If the field is blank: write "✗ [field name]: BLANK — rerun artifact-context before proceeding"

If any field is blank (not "Not yet created", but genuinely empty): stop and rerun artifact-context.

Write: "Cross-reference scan: [N] populated, [M] placeholders, [K] blank."

**Output:** Cross-reference scan result written. Zero blank fields before continuing.

---

### Step 3 — Invoke `jira-ticket-template` (Killer)

Invoke the `jira-ticket-template` skill. Pass all inputs from the Inputs Required section plus the artifact-context metadata block from Step 1.

The returned ticket body must contain all 10 sections. Write the complete ticket body in the conversation.

**Output:** Complete 10-section ticket body written in conversation.

---

### Step 4 — Validate Algorithm and Acceptance Criteria (Killer)

**Algorithm scan (Section 6):** Read every numbered step in the Algorithm section. For each step: does it name an exact Go identifier (type name, function, method, constant, sort key), loop bound, or data structure — or does it use a vague verb ("process", "handle", "iterate", "deal with") without naming a specific construct? Count vague steps.

Write: "Algorithm scan: [N] steps total, [K] vague. [List each vague step verbatim if K > 0]"

If K > 0: rewrite each flagged step to name the exact construct before continuing.

**AC scan (Section 8):** Read every row in the Acceptance Criteria checklist. For each row: is there an exact shell command producing a binary pass/fail result? Count rows without an exact command.

Write: "AC scan: [N] criteria, [K] without exact command. [List each failing row if K > 0]"

If K > 0: add the exact command and expected output to each flagged row before continuing.

Both counts must be 0 before proceeding to Step 5.

**Output:** Algorithm scan and AC scan results written; both counts = 0 before continuing.

---

### Step 5 — Discover Jira Project Key (Killer)

Call `mcp__plugin_greenfi-engineering_atlassian-custom__jira_search` to list projects in the personal workspace.

Read the response. Identify the project key for the personal workspace belonging to aaron.thatcher11@gmail.com.

Write: "Jira project key: [KEY] — sourced from live MCP query, not recalled."

Do not use a project key recalled from a prior session or conversation.

**Output:** Project key stated with "sourced from live MCP query" written before any creation call.

---

### Step 6 — State Parameters and Create Ticket (Killer)

Write every parameter before calling:

```
Project key:  [from Step 5]
Issue type:   Task
Summary:      [Section 1 TITLE from ticket body]
Description:  ADF JSON object (ticket body converted to ADF — NOT a string, NOT Markdown)
Labels:       [from Section 10 LABELS]
```

Convert the ticket body sections to ADF format. The description must be a JSON object beginning with `{"type":"doc","version":1,...}`. Reference the ADF node types in the `confluence-page` skill if needed.

Call `mcp__plugin_greenfi-engineering_atlassian-custom__jira_create_issue` with the above parameters.

Extract the ticket ID (e.g., `CG-12`) from the response.

Write: "Ticket created: [ID]"

**Output:** Parameter block written before MCP call; ticket created; ticket ID extracted and written.

---

### Step 7 — Set Blocker Relationships (Killer)

Read Section 10 BLOCKS from the ticket body.

If BLOCKS lists ticket IDs: for each ID, call `mcp__plugin_greenfi-engineering_atlassian-custom__jira_create_issue_link` with:
- `inwardIssue`: this ticket's ID (from Step 6)
- `outwardIssue`: the blocked ticket's ID
- `type`: "blocks"

Write after each call: "Link created: [this ticket ID] blocks [blocked ticket ID]"

If BLOCKS is "none": write "No blocker relationships to create."

Count BLOCKS entries vs. link calls made — they must match.

Write: "Blocker links: [N] entries in Section 10, [N] link calls made."

**Output:** All blocker links created; count confirmed matching.

---

### Step 8 — Output

Write:

```
Ticket: [ID] — [Section 1 TITLE]
URL: [full Jira ticket URL]
Blocks: [list or "none"]
Blocked by: [DEPENDS ON entries from Section 4, or "none"]
```

**Output:** Final ticket summary written in conversation.

---

## Proof Requirements

- **Step 1:** `artifact-context` skill invocation visible in conversation before Step 3.
- **Step 2:** Cross-reference scan result written; zero blank fields before Step 3.
- **Step 3:** `jira-ticket-template` skill invocation visible; full 10-section body written.
- **Step 4:** Algorithm scan count = 0 vague steps AND AC scan count = 0 rows without command — both written before Step 5.
- **Step 5:** "sourced from live MCP query" written before `jira_create_issue` call.
- **Step 6:** Parameter block written before MCP call; ticket ID extracted and written.
- **Step 7:** Link call count = BLOCKS entry count, written in conversation.

---

## State Persistence

N/A — single session. Output is the ticket URL written in conversation.

---

## Red Flags

| Thought | Reality |
|---------|---------|
| "I know the artifact-context metadata — no need to invoke the skill" | PRD and Tech Req URLs may not exist yet or may have changed this session. Invoke the skill. |
| "The jira-ticket-template output looks complete — skip the validation scans" | "Looks complete" is not the same as zero vague algorithm steps and zero AC rows without commands. Run the scans. |
| "I know the project key from last session" | Project key must come from a live `jira_search` query this session. Recalled keys can be wrong after workspace changes. |
| "The description looks fine as Markdown" | Jira stores Markdown as literal text. Description must be ADF JSON object beginning with `{"type":"doc",...}`. |
| "BLOCKS says 'none' — no link step needed" | Write "No blocker relationships to create." explicitly. Do not silently skip Step 7. |
| "The blocker relationship is implied by the ticket title" | Implied relationships are invisible to agents reading the ticket. Links must be created via `jira_create_issue_link`. |

---

## Integration

```
Called by: task-breakdown (for each ticket in the phase breakdown); standalone for single tickets
Calls: artifact-context, jira-ticket-template
Output: Jira ticket URL written in conversation; ticket created with 10-section body, ADF description, blocker links set
MCP tools: mcp__plugin_greenfi-engineering_atlassian-custom__jira_create_issue
           mcp__plugin_greenfi-engineering_atlassian-custom__jira_create_issue_link
           mcp__plugin_greenfi-engineering_atlassian-custom__jira_search
```

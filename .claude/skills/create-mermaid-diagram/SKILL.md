---
name: create-mermaid-diagram
description: Produces valid Mermaid diagram syntax for one of the 4 known clock-gate diagram types, built from actual source files. TRIGGER when: create-design-doc or create-tech-requirements needs a diagram, user says "generate the diagram", "draw the architecture", "create the [data flow / StateMatrix / package / serve mode] diagram", or a diagram needs updating after a code change.
---

## Overview

Architecture diagrams drawn from memory describe intent, not reality. This skill reads the actual source files, builds node lists from real identifiers, and produces syntax that renders correctly and accurately represents the current codebase.

**Iron Law:** Every node must represent a real artifact that exists in the current codebase, or be explicitly labeled "(future)."

---

## Announce at Start

"Generating [diagram name] diagram. Reading source files now — every node will be traced to a real code artifact before the syntax is written."

---

## The 4 Known Diagram Types

### Type 1 — Data Flow (system pipeline)
**Relationship:** Data moves through distinct transformation stages (A produces B).
**Mermaid type:** `flowchart LR`
**Source files to read:** `cmd/clock-gate/main.go`, `internal/ingestion/parser.go`, `internal/matrix/compiler.go`, `internal/gate/gate.go`, `internal/presenter/snapshot.go`

Template:
```
flowchart LR
    JSON["JSON file\n~90KB"] --> Parser["ingestion.ParseFile()"]
    Parser --> Plays["[]RawPlay"]
    Plays --> Compiler["matrix.Compile()"]
    Compiler --> Matrix["StateMatrix\n~252KB"]
    Matrix --> Gate["gate.Validate()"]
    Gate --> Presenter["presenter.RenderText()"]
    Presenter --> Output["CLI output"]
```

### Type 2 — StateMatrix Internals
**Relationship:** Algorithm visualization — how irregular plays map onto a dense array.
**Mermaid type:** `flowchart TB`
**Source files to read:** `internal/matrix/types.go`, `internal/matrix/compiler.go`

Template:
```
flowchart TB
    subgraph Input["Input: N plays (irregular)"]
        P0["Play @0s\nplay_id=1"]
        P150["Play @150s\nplay_id=47"]
        P900["Play @900s\nplay_id=89"]
    end
    subgraph Array["[9001]GameState — indexed by elapsed second"]
        S0["States[0]"]
        S1["States[1]–States[149]\nforward-filled"]
        S150["States[150]"]
        S151["States[151]–States[899]\nforward-filled"]
        S900["States[900]"]
    end
    P0 --> S0
    P150 --> S150
    P900 --> S900
    S0 -->|"Compile(): forward-fill"| S1
    S150 -->|"Compile(): forward-fill"| S151
    Query["Query: States[T]"] -->|"O(1) — one array dereference"| Array
```

### Type 3 — Package Structure
**Relationship:** Ownership — which packages import which (A depends on B).
**Mermaid type:** `graph TD`
**Source files to read:** Import blocks in all `*.go` files under `cmd/` and `internal/`.

Template:
```
graph TD
    main["cmd/clock-gate/main.go"]
    ingestion["internal/ingestion"]
    matrix["internal/matrix"]
    gate["internal/gate"]
    presenter["internal/presenter"]
    main --> ingestion
    main --> matrix
    main --> gate
    main --> presenter
    ingestion -.-> matrix
```
*(solid arrow = direct import; dashed = type dependency only)*

### Type 4 — Serve Mode (future state)
**Relationship:** Concurrent access — many users sharing one pre-compiled StateMatrix per game.
**Mermaid type:** `flowchart LR`
**Source files to read:** Check if `cmd/clock-gate/serve.go` and `internal/server/` exist. If not: all serve-mode nodes are future state.

Template:
```
flowchart LR
    Users["N concurrent users\nwatching game X"] --> HTTP["HTTP GET\n/game/:id/state?tick=T"]
    HTTP --> Server["clock-gate serve\n(future)"]
    Server --> Cache["StateMatrix cache\none per game"]
    Cache -->|"States[T] — O(1)\nno new allocations"| Response["JSON response"]
    Response --> Users
    style Server stroke-dasharray: 5 5
    style Cache stroke-dasharray: 5 5
```

---

## Steps

### Step 1 — Identify Diagram Type and Justify (Killer)

Write:
- "Diagram requested: [name]"
- "Relationship being shown: [data moving through stages / algorithm / ownership/imports / concurrent access]"
- "Correct Mermaid type: [flowchart LR / flowchart TB / graph TD] — because [one-sentence reason]"

Do not write any Mermaid node or edge syntax until this statement is written.

**Output:** Three-line type justification written.

---

### Step 2 — Read Source Files (Killer)

Read every source file listed for this diagram type in the Known Diagram Types section above, using the Read tool.

If a file does not exist: write "FILE NOT FOUND: [path] — nodes for this component will be marked (future)."

Do not proceed to Step 3 using session memory of what these files contain.

**Output:** All relevant source files read with explicit Read tool calls; any missing files noted.

---

### Step 3 — Build Node and Edge List from Real Identifiers (Killer)

From the source files read in Step 2, extract actual identifiers:
- Function names (e.g., `ParseFile`, `Compile`, `Validate`)
- Package names (e.g., `ingestion`, `matrix`)
- Type names (e.g., `StateMatrix`, `GameState`)
- File paths (e.g., `cmd/clock-gate/main.go`)

Write an explicit node list:
```
Node: [label] — represents: [exact function name / file path / type name from source]
Edge: [source] --> [target] — represents: [function call in source / Go import statement]
```

Any label that is not a plain Go identifier must be quoted in the Mermaid syntax: `["label with spaces"]`.

**Output:** Node and edge list with real-artifact citations.

---

### Step 4 — Write Mermaid Syntax

Using the template for this diagram type and the node list from Step 3, write the complete Mermaid syntax.

**Syntax rules:**
- Labels with `(`, `)`, `/`, spaces, or `\n`: wrap in double quotes inside brackets — `["ParseFile()"]`
- Newlines in labels: use `\n` inside quoted string — `["line1\nline2"]`
- Arrow types: `-->` plain, `-->|"label"|` labeled, `-.->` dashed, `==>` thick
- No bare `->` — always use `-->`
- Subgraph IDs must be alphanumeric: `subgraph ArrayID["Display Name"]`
- Style references node ID, not label: `style NodeID stroke-dasharray: 5 5`
- Use `flowchart` when data moves through stages. Use `graph` when showing import/ownership hierarchy.

Write as a fenced code block:
````
```mermaid
[syntax here]
```
````

**Output:** Complete fenced Mermaid code block.

---

### Step 5 — Cartographer Accuracy Scan (Killer)

For each node in the diagram:
- If the node represents a file or directory: run `ls [path]` — exit code 0 means ✓ exists; non-zero means mark (future).
- If the node represents a function or type: run `grep -r "[identifier]" [package dir]` — match found means ✓ exists; no match means mark (future).

For any node not confirmed:
1. Update its label to include "(future)"
2. Add `style NodeID stroke-dasharray: 5 5` to the diagram

Write scan results:
```
Node [label]: ✓ exists at [path] / marked (future)
```

**Output:** Scan result for every node; diagram updated for future-state nodes.

---

### Step 6 — Syntax Scan

Read the Mermaid syntax from Step 4 (with Step 5 updates applied).

For each node label:
- Does it contain `(`, `)`, `/`, `\`, or spaces outside of `"..."` quotes? If yes: wrap in quotes.

For each arrow:
- Is it a bare `->` (single dash)? If yes: replace with `-->`.

Write "Syntax scan: PASS — no issues found" or list each fix made.

**Output:** Syntax scan result written.

---

### Step 7 — Output

Write:

```
## [Diagram Name] Diagram

[fenced mermaid code block — final version from Steps 4/5/6]

Embedding:
- README.md: paste the fenced block directly — renders on GitHub automatically
- Confluence: paste via confluence-page skill as a Mermaid macro node in ADF
```

**Output:** Final diagram block with embedding instructions.

---

## Proof Requirements

- **Step 1:** Three-line type justification written before any Mermaid node syntax.
- **Step 2:** Read tool call visible for each source file; missing files explicitly noted.
- **Step 3:** Node list written with real-artifact citation for each node and edge.
- **Step 5:** Every node has a scan result — "✓ exists at [path]" or "marked (future)."
- **Step 6:** "Syntax scan: PASS" or list of specific fixes written before output.

---

## State Persistence

N/A — single session. Output is the fenced Mermaid code block written in the conversation.

---

## Red Flags

| Thought | Reality |
|---------|---------|
| "I know the package structure from this session — no need to read the files" | Cartographer anchor: memory produces aspirational diagrams. Read the files. |
| "flowchart works for everything" | `flowchart` shows data moving through stages. `graph TD` shows ownership/imports. Wrong type produces the wrong mental model for the reader. |
| "The serve mode components are obvious — I'll include them" | Every node must exist in the codebase or be marked (future). Serve mode files don't exist yet. Mark them explicitly with dashed style. |
| "The label looks fine without quotes" | Labels with parentheses or spaces break rendering silently — no error, just blank or partial diagram. Step 6 catches these. |
| "I'll use the template directly — it matches this project" | Templates are starting points. Real function names, import paths, and package names may differ. Read the source first; update the template from Step 3 output. |

---

## Integration

```
Called by: create-design-doc, create-tech-requirements; standalone for diagram updates
Calls: none
Output: Fenced Mermaid code block written in conversation with embedding instructions
```

---
name: confluence-page
description: Creates a Confluence page in the paradox-platform workspace (space PDX) nested under the paradox-clock-gate container page, and returns the URL to the calling skill. TRIGGER when: any artifact-generating skill (create-prd, create-tech-requirements, create-design-doc, task-breakdown) needs to create a Confluence page.
---

## Overview

Each artifact skill that creates Confluence pages independently re-derives space discovery, content format, and URL extraction — and each makes a different mistake. This skill centralizes that logic.

**Iron Law:** Content passed to `confluence_create_page` must be Confluence storage format (HTML-like XHTML). Never pass ADF JSON. Never pass Markdown. The MCP tool's `content` parameter is a string, not an object.

**Organization rule:** All paradox-clock-gate documents nest under the container page (ID: 21397505) at `https://paradox-platform.atlassian.net/wiki/spaces/PDX/pages/21397505/paradox-clock-gate`. Pass `parent_id: "21397505"` on every create call unless creating the container page itself.

---

## Announce at Start

"Invoking confluence-page. Space: PDX, parent: paradox-clock-gate container (21397505). Will format body as Confluence storage format, create the page, and return the URL."

---

## Confluence Storage Format Reference

Common storage format elements:

```html
<!-- Paragraph -->
<p>Your text here.</p>

<!-- Heading -->
<h2>Section Title</h2>
<h3>Subsection</h3>

<!-- Bullet list -->
<ul>
  <li>Item one</li>
  <li>Item two</li>
</ul>

<!-- Numbered list -->
<ol>
  <li>First</li>
  <li>Second</li>
</ol>

<!-- Inline code -->
<code>someIdentifier</code>

<!-- Code block -->
<ac:structured-macro ac:name="code">
  <ac:parameter ac:name="language">go</ac:parameter>
  <ac:plain-text-body><![CDATA[code here]]></ac:plain-text-body>
</ac:structured-macro>

<!-- Strong / emphasis -->
<strong>bold</strong>
<em>italic</em>

<!-- Table -->
<table>
  <tr><th>Column A</th><th>Column B</th></tr>
  <tr><td>Value 1</td><td>Value 2</td></tr>
</table>

<!-- Link -->
<a href="https://example.com">Link text</a>
```

---

## Steps

### Step 1 — Accept Inputs

The calling skill provides:
- `title`: the page title string
- `content`: the page body prose or structured content to convert in Step 3
- `parent_page_id` (optional override): if not provided, use the default container page ID 21397505

Write: "Creating Confluence page titled '[title]'. Parent: [parent_page_id or 21397505 default]."

**Output:** One-line confirmation with parent ID stated.

---

### Step 2 — Confirm Space and Parent (Killer)

Space is PDX — confirmed via live query at session start (pages at `/rest/api/space/PDX` returned results).

Parent page ID: use `parent_page_id` from calling skill if provided; otherwise use `21397505` (paradox-clock-gate container page).

Write: "Space: PDX (confirmed). Parent page ID: [id] — [source: provided by calling skill / default container]."

Do not create pages at space root unless explicitly instructed.

**Output:** Space and parent stated before any creation call.

---

### Step 3 — Assemble Confluence Storage Format Body (Killer)

Convert the calling skill's content to Confluence storage format HTML using the reference above.

Rules:
- Each section heading → `<h2>` or `<h3>`
- Each paragraph of prose → `<p>...</p>`
- Bullet lists → `<ul><li>...</li></ul>`
- Numbered lists → `<ol><li>...</li></ol>`
- Inline identifiers (Go types, function names, flags) → `<code>...</code>`
- Code blocks → `<ac:structured-macro ac:name="code">` with language parameter
- Special characters in text: `&amp;`, `&lt;`, `&gt;`, `&quot;`, `&apos;`

Write the complete storage format HTML as a code block in the conversation.

**Output:** Full storage format HTML written in conversation as a code block.

---

### Step 4 — Storage Format Compliance Check (Killer)

Read the HTML written in Step 3. It passes when ALL of the following are true:

1. It is a string, not a JSON object or array
2. It contains at least one block-level element (`<h2>`, `<h3>`, `<p>`, `<ul>`, `<ol>`, `<table>`)
3. All tags are properly closed (every `<tag>` has a matching `</tag>` or is self-closing)
4. No raw Markdown syntax present (`##`, `**`, `- `, `` ` ``)
5. Special characters in text content are HTML-escaped (`<` as `&lt;`, etc.)

Write "Storage format compliance: PASS" if all five conditions are met.

If any condition fails: write "Storage format compliance: FAILED — [specific condition]." Fix and re-run. Do not proceed to Step 5 until PASS is written.

**Output:** "Storage format compliance: PASS" or "Storage format compliance: FAILED — [reason]."

---

### Step 5 — State All MCP Parameters (Killer)

Write every parameter before calling the MCP tool:

```
Space key:  PDX
Title:      [from calling skill]
Parent ID:  [from Step 2]
Content:    Confluence storage format string from Step 3
```

**Output:** Parameter block written in conversation.

---

### Step 6 — Create Page

Call `mcp__paradox-confluence__confluence_create_page` with:
- `space_key`: `"PDX"`
- `title`: title from calling skill
- `parent_id`: parent page ID from Step 2
- `content`: the storage format HTML string from Step 3 (pass as a string, not as an object)

**Output:** MCP tool call executed; response received.

---

### Step 7 — Extract and Return URL (Killer)

Extract the page URL from the MCP response.

Write on its own line:
`PAGE URL: [full url]`

Then write: "confluence-page complete. Page '[title]' created under paradox-clock-gate container. Calling skill: use the PAGE URL above for artifact-context cross-references."

**Output:** A line beginning exactly "PAGE URL:" written in the conversation.

---

## Proof Requirements

- **Step 2:** Space PDX and parent page ID stated before any creation call.
- **Step 3:** Full storage format HTML written as code block before Step 4 executes.
- **Step 4:** "Storage format compliance: PASS" written before Step 5 — no creation call made until compliance passes.
- **Step 5:** Parameter block written with all four fields before `confluence_create_page` is called.
- **Step 7:** A line beginning exactly "PAGE URL:" written in conversation after successful creation.

---

## State Persistence

N/A — single session. Output is the PAGE URL written in conversation for the calling skill.

---

## Red Flags

| Thought | Reality |
|---------|---------|
| "I'll pass an ADF JSON object as the content parameter" | The MCP tool's `content` field is a string. ADF objects silently fail or render as literal JSON text. Use storage format HTML. |
| "I'll pass Markdown — Confluence will render it" | Confluence stores Markdown as literal text. Storage format compliance check must pass before any creation call. |
| "No parent was specified — I'll create at space root" | All paradox-clock-gate documents go under container page 21397505. Space root is only for the container page itself. |
| "The MCP call succeeded — the URL is in the response somewhere" | "Somewhere in the response" is not usable. Step 7 requires a line beginning exactly "PAGE URL:" in the conversation. |
| "I'll skip the parameter list — the call is straightforward" | Unverified parameters silently create the page in the wrong space or with wrong nesting. State all parameters before calling. |
| "The HTML looks right — I'll skip the compliance check" | Unclosed tags and raw Markdown produce pages that render as garbage. Step 4 is the only gate. |

---

## Integration

```
Called by: create-prd, create-tech-requirements, create-design-doc, task-breakdown
Calls: mcp__paradox-confluence__confluence_create_page (page creation)
Space: PDX (paradox-platform.atlassian.net)
Container page: ID 21397505 — https://paradox-platform.atlassian.net/wiki/spaces/PDX/pages/21397505/paradox-clock-gate
Output: "PAGE URL: [url]" written in conversation for the calling skill to capture
```

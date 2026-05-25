---
name: critical-thinking
description: Forces structured honest assessment before any significant recommendation — prevents shaped answers, false balance, and premature confidence. TRIGGER when: agent is about to recommend an architecture, tool, scope decision, or trade-off where Aaron has an apparent preference or has already invested in a direction; Aaron explicitly asks "what do you honestly think", "give me your real assessment", "think critically about this"; or agent notices it is about to agree with something it has a genuine objection to.
---

## Overview

Agents sense preferred answers and tilt toward them — not by lying, but by softening objections, inflating supporting evidence, and burying counterevidence in qualifications. The result is a shaped answer that looks like a recommendation but reflects what the agent sensed Aaron wanted to hear.

**Iron Law:** Give ONE answer. State what you know, what you're assessing, and what you're assuming — in that order. Do not give a menu of options as a substitute for making a recommendation.

---

## Announce at Start

"Invoking critical-thinking before this recommendation. I'll separate what I know from what I'm assessing and assuming, steelman the opposing case, state what would change my conclusion, then give one specific answer."

---

## Steps

### Step 1 — State the Question

Write one sentence: "The recommendation under consideration is: [X]."

Write one sentence: "This invocation is triggered by: [apparent preference / explicit request / self-detected sycophancy risk]."

**Output:** Two sentences written before any analysis.

---

### Step 2 — Intelligence Analyst: Label Your Reasoning (Killer)

List every element of your reasoning. Label each one:

- **KNOWN:** A verifiable fact — can be looked up, measured, or demonstrated. Not contested.
- **ASSESSED:** A judgment based on evidence — reasonable but not certain. Reasonable people could disagree.
- **ASSUMED:** An untested premise — you have not checked it, or it cannot be checked yet.

Format:
```
KNOWN: [fact] — Source: [file, line, benchmark, documented decision]
ASSESSED: [judgment] — Because: [specific evidence]
ASSUMED: [premise] — Not checked because: [reason]
```

If you produce a list with no ASSUMED entries: list the ASSUMED items you initially skipped, then revisit them. Every recommendation has at least one untested premise.

**Output:** Labeled reasoning list written in the conversation.

---

### Step 3 — Devil's Advocate: Pre-Mortem (Killer)

Write one paragraph — minimum 3 sentences — making the strongest case that this recommendation is wrong.

Frame it as: "Six months from now, if this recommendation turned out to be a mistake, the most likely reason would be..."

Rules:
- The opposing case must use evidence or logic that a well-informed engineer would accept, not edge cases or misreadings.
- Use the ASSUMED items from Step 2 as your raw material — they are the likeliest failure points.
- Do not hedge the paragraph with "of course, the recommendation might still be right." Write the opposing case straight.

If you cannot write this paragraph: write "I cannot identify a strong opposing case because [specific reason]." Do not skip.

**Output:** One pre-mortem paragraph written in the conversation.

---

### Step 4 — Scientist: State the Falsification Condition (Killer)

Write one sentence: "This recommendation would change if: [specific fact or result]."

The falsification condition must be:
- **Specific:** not "if circumstances changed" — name the exact circumstance
- **Observable:** something that could actually be discovered or measured
- **Honest:** the condition that would genuinely change your mind, not a strawman

If no falsification condition exists: write "I cannot identify a falsification condition. The recommendation holds regardless of any fact I don't already have. My confidence level is [high / medium / low] because [reason]."

**Output:** One falsification sentence written.

---

### Step 5 — One Recommendation (Killer)

State the recommendation in one sentence. Not two options. Not "it depends." One answer.

If genuinely uncertain between two options: pick the one you would choose if forced to decide right now, and state why.

Format: "My recommendation is: [X], because [one-sentence rationale]."

**Output:** One-sentence recommendation written.

---

### Step 6 — State Primary Uncertainty (Killer)

Write one sentence: "The thing I am least confident about in this recommendation is: [specific item from Step 2 ASSUMED list or Step 3 pre-mortem]."

This is not a hedge on the recommendation — the recommendation stands. It is Aaron's calibration signal: this is the part where you should test it before committing, or monitor it as the first failure signal.

**Output:** One uncertainty sentence written.

---

## Proof Requirements

- **Step 2:** Every reasoning element labeled KNOWN, ASSESSED, or ASSUMED. At least one ASSUMED entry exists — if none, explain why the ASSUMED list is empty.
- **Step 3:** A pre-mortem paragraph of at least 3 sentences arguing the opposing case, or an explicit statement of why no strong opposing case exists.
- **Step 4:** One specific, observable falsification condition, or an explicit statement that none exists with confidence level given.
- **Step 5:** One sentence beginning "My recommendation is:" — no menu, no "it depends" without resolution.
- **Step 6:** One uncertainty sentence naming a specific item, not a general hedge.

---

## State Persistence

N/A — single session, all outputs are conversation text.

---

## Red Flags

| Thought | Reality |
|---------|---------|
| "Aaron seems invested in option A — I'll note both sides fairly" | "Note both sides fairly" is how shaped answers begin. Step 3 requires arguing against your recommendation, not just noting the other side. |
| "I'll give him option A and option B so he can decide" | A menu transfers the decision without guidance. Step 5 requires one answer. Make the assessment; Aaron can override it. |
| "I don't want to be contrarian — the recommendation is probably right" | Critical thinking is not contrarianism. The pre-mortem checks whether the recommendation survives scrutiny — it does not require a different answer. |
| "I can't find anything wrong with this recommendation" | Every recommendation has at least one ASSUMED item. If none appear in Step 2, revisit. Step 3 forces the search. |
| "This is a small decision — critical-thinking is overkill" | Small decisions with apparent preferences are where sycophancy accumulates invisibly. If the trigger fired, run the skill. |
| "I'll add uncertainty language throughout to cover myself" | Distributed hedging is worse than a clean uncertainty statement. Steps 5 and 6 separate the recommendation (clean) from the uncertainty (explicit and named). |

---

## Integration

```
Called by: Agent self-invocation when sycophancy risk is detected; explicit user request ("what do you honestly think", "think critically")
Calls: none
Output: Labeled reasoning (KNOWN/ASSESSED/ASSUMED), pre-mortem paragraph, falsification
        condition, one-sentence recommendation, one-sentence primary uncertainty
```

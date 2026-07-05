---
name: adr-writer
description: Write and review Architecture Decision Records for repository-level architecture, dependency, infrastructure, API, data-modeling, testing, observability, and operational decisions. Use when asked to write, review, update, supersede, or check an ADR, or when a change introduces a meaningful architectural trade-off.
---

# ADR Writer Skill

## Purpose

Create and review Architecture Decision Records that explain:

- what decision was made;
- why the decision was needed;
- what alternatives were considered;
- what trade-offs were accepted;
- how the decision affects implementation;
- how future contributors and agents can verify compliance.

The ADR must be useful inside this repository, not a generic architecture essay.

## When to use

Use this skill when:

- the user asks to write, review, update, supersede, or check an ADR;
- a change introduces or removes a major dependency;
- a change affects storage, queues, APIs, auth, logging, metrics, configuration, migrations, deployment, testing strategy, or other cross-cutting patterns;
- the repository is choosing between multiple implementation approaches;
- future maintainers may reasonably ask: "Why did we do it this way?"

Do not create an ADR for routine implementation details, small bug fixes, formatting-only changes, or decisions already covered by an accepted ADR.

## Repository ADR discovery

Before writing or reviewing an ADR, inspect existing conventions.

Look for ADRs in:

- `adrs/`
- `docs/decisions/`
- `docs/adr/`
- `adr/`
- `decisions/`

Read:

- the ADR index, if present;
- the last 3-5 accepted ADRs;
- ADRs related to the current area;
- relevant code, configs, migrations, CI files, manifests, and documentation.

Do not contradict accepted ADRs silently. If a new decision conflicts with an accepted ADR, either mark the older ADR as superseded or explicitly recommend a superseding ADR.

## Repository ADR location

This repository stores ADRs in:

```text
adrs/
```

Use filenames like:

```text
NNNN-short-title-with-dashes.md
```

Example:

```text
adrs/0006-use-pretty-print-false-for-nilaway-snapshots.md
```

## Repository ADR template

Use this lightweight template unless the repository adopts a stricter one:

```markdown
# ADR NNNN: Short Decision Title

## Status

Proposed

## Context

Describe the situation that forced this decision.

Include:

- what changed or broke;
- why the decision matters now;
- relevant repository context;
- constraints and forces;
- related ADRs, issues, PRs, or code paths.

## Decision

Chosen option: **Option X**.

Explain why this option best fits the drivers and constraints.

## Alternatives Considered

- Option A
- Option B
- Option C

## Consequences

Positive:

- ...

Negative:

- ...

Risks:

- ...

Mitigations:

- ...

## Implementation Notes

Affected areas:

- `path/to/file_or_directory`

Implementation rules:

- ...

Patterns to follow:

- ...

Patterns to avoid:

- ...

## Verification

This ADR is implemented when:

- [ ] ...
- [ ] ...
- [ ] ...
```

Shorter existing ADRs may use only `Status`, `Context`, `Decision`, and `Consequences`, but new ADRs should include alternatives and verification when the decision has implementation impact.

## ADR review checklist

When reviewing an ADR, check:

- Is the decision specific enough to implement?
- Is the context repository-specific?
- Is the trigger clear?
- Are alternatives honestly represented?
- Is the chosen option explicit?
- Are rejected options explained?
- Are consequences concrete?
- Are risks and mitigations named?
- Are affected files, packages, configs, workflows, or docs listed?
- Are verification criteria checkable?
- Does the ADR conflict with any accepted ADR?
- Should it supersede or update an older ADR?
- Would a coding agent be able to implement the decision without guessing?

Output the review as:

```markdown
## ADR review

### Passes

- ...

### Gaps

- ...

### Suggested changes

- ...

### Recommendation

Ship it / revise before accepting / needs more context
```

## Status lifecycle

Allowed statuses in this repository:

- `Proposed`
- `Accepted`
- `Rejected`
- `Deprecated`
- `Superseded`

Existing ADRs use a `## Status` section rather than YAML front matter. Preserve that convention unless the repository intentionally changes it.

Never rewrite historical rationale silently. If reality changed, append an update with a date or create a superseding ADR.

## ADR index

The ADR index lives at:

```text
adrs/README.md
```

Update it when creating, renaming, or superseding ADRs.

## Code links

When an ADR governs specific code, link both ways when useful.

In ADRs, mention affected paths.

In code, add lightweight comments only at important entry points, for example:

```go
// ADR-0004: NilAway snapshots must run with deterministic no-color output.
```

Do not spam ADR references across trivial files.

## Style

Prefer:

- direct language;
- short paragraphs;
- concrete trade-offs;
- repository-specific names and paths;
- explicit assumptions;
- checkable verification criteria.

Avoid:

- vague "best practice" claims;
- generic architecture filler;
- pretending there were no trade-offs;
- long essays;
- unverifiable claims like "more scalable", "cleaner", or "better" without explanation.

# ADR Review

Date: 2026-07-05

## Summary

The repository already uses ADRs. They live in `adrs/`, use four-digit numeric filenames, and currently use a lightweight Markdown structure with `## Status`, `## Context`, `## Decision`, and `## Consequences`.

This review was performed after adding the repo-local `$adr-writer` skill.

## Existing ADR Convention

- Location: `adrs/`
- Numbering: `NNNN-short-title-with-dashes.md`
- Heading format: `# ADR NNNN: Title`
- Status format: `## Status` section with `Accepted`
- Index: `adrs/README.md`

## Findings

- The five existing ADRs are consistently numbered from `0001` to `0005`.
- ADR statuses are consistent: all current ADRs are `Accepted`.
- ADR 0005 correctly captures the GitHub Actions container safe-directory decision and references the workflow behavior it governs.
- ADR 0004 captures the NilAway ANSI/cache behavior that affects snapshot stability.
- Earlier ADRs are intentionally lightweight, but they do not always include alternatives considered or explicit verification criteria.
- The ADR index was missing before this review and has now been added.

## Recommended Fixes

- For future ADRs with implementation impact, include alternatives considered and verification criteria.
- When changing NilAway snapshot behavior, check ADR 0004 first.
- When changing GitHub Actions container behavior or the final `git diff` snapshot check, check ADR 0005 first.
- Keep `AGENTS.md` short and put detailed ADR-writing rules in `.agents/skills/adr-writer/SKILL.md`.

## Follow-up ADR Candidates

- Explicitly set NilAway `pretty-print: "false"` in golangci-lint configs for deterministic snapshots.
- Add a new tool family beyond NilAway, such as `goleak`, if it introduces new snapshot or CI conventions.

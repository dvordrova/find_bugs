# Implemented Examples Documentation Cleanup

Commit: `9fb3c1f Clarify implemented bug examples`

## Context

After implementing the planned 11 examples, `BUGS.md` still called the final list `First Implementation Order`.

## Options

1. Leave the heading unchanged.
   Score: 5/10. Technically harmless, but stale wording after the implementation work.

2. Rename the section to `Implemented Examples` and link each README.
   Score: 10/10. Makes the catalog easier to navigate and accurately reflects the current repository.

3. Remove the section entirely and rely on root README.
   Score: 6/10. Avoids duplication, but `BUGS.md` is the planning/catalog file and should show what has landed.

4. Split implemented and future candidates into separate documents.
   Score: 5/10. Better later, but too much process for this cleanup.

## Chosen

Option 2.

## Why

The repository now has all 11 planned examples. A linked implemented list is more useful than a historical implementation order.

## Verification

- `git diff --check`
- `make test` had passed immediately before this documentation cleanup.

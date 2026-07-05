# Tooling Contract Documentation Decision

Commit: `740cbcb Clarify local and CI tooling contract`

## Options

1. Leave the contract implicit in Makefiles.
   - Score: 2/5
   - The structure is there, but a new reader has to infer the local/CI workflow.

2. Add a root README "Tooling Contract" and set `make` to the full repository check.
   - Score: 5/5
   - It makes the intended user path explicit: local one-command check, per-example detector command, and CI snapshot command.

3. Add only CONTRIBUTING guidance.
   - Score: 3/5
   - Good for contributors, but users and testers often start from README and would miss the workflow.

4. Add a separate long process document.
   - Score: 2/5
   - Too much ceremony for a small repository; the important contract fits in README and CONTRIBUTING.

## Selected

Option 2, with the CONTRIBUTING checklist from option 3.

The root Makefile now makes `make` equivalent to `make test`. README explains the stable targets each example should expose: `run`, `lint`, `test`, and `ci-test`. CONTRIBUTING adds an acceptance checklist so future examples keep teaching how to avoid the bug with tools, not just how to recognize it during review.

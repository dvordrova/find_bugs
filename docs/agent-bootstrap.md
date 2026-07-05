# Agent Bootstrap Prompt

Use this when you want another coding agent to reproduce the repository style and tooling contract from `find_bugs` in a new Go repository.

## Copy-Paste Prompt

```text
Use the `find_bugs` repository as the reference style.

Build a catalog of small, runnable Go bug examples and the tools that catch them. The result must be useful for local developers and CI, not just a set of code snippets.

Repository contract:

- Every example lives in `<tool_or_problem>/<category>/`.
- Every example is self-contained and usually contains:
  - `go.mod`
  - `Makefile`
  - `main.go`
  - `README.md`
  - `README.ru.md`
- Example code should feel production-like, but stay small enough to understand quickly.
- Do not comment every line. Put the explanation in README files.
- Prefer real tooling behavior over made-up output.
- Pin tool versions through Go tool dependencies or explicit module versions.
- Do not use `go get -u` for routine setup.

Makefile contract for every example:

- Put user-facing targets first:
  - `make run` when the program behavior matters
  - `make lint` or the closest detector command
  - `make test`
- Put maintainer/CI targets later:
  - `make ci-test`
  - `make tool-update`
  - generated helper binary targets
  - `make clean`
- `make lint` may fail when the example intentionally demonstrates a bug.
- `make ci-test` must be stable in CI: it should regenerate committed snapshot files such as `lint.logs`, `race.logs`, or `lint-fixed.logs`.
- If module downloads can pollute snapshot logs, run `go mod download` before the snapshot-producing command.

Root repository contract:

- `make test-update` walks all example directories and runs their `ci-test` target.
- `make test` runs `make test-update`, then `git diff --exit-code`.
- The default `make` target should be `test`.
- CI should run one command: `make test`.
- If examples accept a custom golangci-lint config, root `make test LINT_CONFIG=/absolute/path/to/.golangci.yaml` should forward it to relevant examples.

Documentation contract:

- Root `README.md` explains purpose, quick start, tool contract, example list, tools, and why snapshots exist.
- Root `README.ru.md` mirrors the important user-facing parts in Russian.
- `BUGS.md` lists planned and implemented bug patterns.
- `CONTRIBUTING.md` explains how to add one example.
- ADRs are used only for repository-level decisions, not for routine examples.
- Each example README explains:
  - what production bug is demonstrated;
  - how to run it locally;
  - what tool catches it;
  - how to read the expected output;
  - one realistic fix or mitigation.

Acceptance check:

After your changes, a new developer or tester should be able to run:

```sh
make
```

and understand failures through committed snapshot diffs. For one example, they should be able to run:

```sh
cd <tool_or_problem>/<category>
make run
make lint
make test
```

and understand what bug was caught without reading the whole repository.

Avoid:

- examples that require global tools to be installed manually;
- hidden setup steps outside Makefiles;
- unpinned tool upgrades;
- debug-only files committed to the repository;
- snapshots that contain cold-cache download noise;
- code that is so artificial that the bug no longer resembles production.
```

## How To Use It

Give the prompt above to an agent together with a link or local path to this repository.

If the target repository already exists, ask the agent to first inspect its current build system and adapt the contract instead of replacing everything. The goal is to preserve the same operational shape:

- one command for local and CI checks;
- small runnable examples;
- committed tool output snapshots;
- pinned tool versions;
- docs that teach the bug, not only the command.

## Good Result

A good clone of this approach should make a reader think:

- I can run the whole catalog with one command.
- I can run one example without reading global setup docs.
- I can see which tool catches the bug.
- I can compare my own lint config against the catalog.
- I can review tool behavior changes through `git diff`.

## Useful Follow-Up Prompt

```text
Review this repository against `docs/agent-bootstrap.md`.

Tell me which parts already match the contract, which parts are missing, and make the smallest changes needed so a developer can run the full check with one command and inspect per-example tool output snapshots.
```

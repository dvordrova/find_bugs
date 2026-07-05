# Contributing

This repository is a catalog of small, runnable Go bug examples and the tools that catch them.

Keep contributions focused: one example should teach one bug pattern or one tool behavior.

## Add A New Example

1. Create a self-contained example directory:

   ```text
   <tool_or_problem>/<category>/
   ```

2. Include the usual files:

   ```text
   go.mod
   Makefile
   main.go
   README.md
   README.ru.md
   ```

3. Keep user-facing Make targets obvious:

   ```sh
   make run
   make lint
   make test
   ```

   Use the tool-specific target name only when `lint` does not fit.

4. Put maintainer and CI targets after user-facing targets.

   `ci-test` should regenerate committed snapshot logs such as `lint.logs` or `lint-fixed.logs`.

5. Pin tool versions in the example Makefile.

   Do not use `go get -u` for routine updates.

6. Regenerate snapshots:

   ```sh
   make test-update
   ```

7. Review changed `*.logs` files.

   Snapshot diffs should be understandable to a reader who is learning the tool.

8. Run the repository check:

   ```sh
   make test
   ```

## Acceptance Checklist

Before an example is done, a new reader should be able to say:

- I know which production bug this demonstrates.
- I know which local command catches it.
- I know what CI command protects the snapshot.
- I know whether the signal comes from `go tool`, `go vet -vettool`, `go test -race`, or goleak.
- I know the main fix or mitigation.

To compare a custom golangci-lint config against the existing examples, run:

```sh
make test LINT_CONFIG=/absolute/path/to/.golangci.yaml
```

The root Makefile forwards that config into golangci-lint examples and keeps non-golangci examples on their normal tools. Use an absolute path because each example runs in its own directory.

## Documentation

Each example README should explain:

- what bug or behavior is demonstrated;
- how to run it;
- what output to expect;
- how to read the tool report;
- one realistic fix or mitigation.

Avoid commenting every line of code. Prefer clear code and focused explanation in the README.

## ADRs

Use ADRs only for repository-level decisions: tooling strategy, CI snapshot model, dependency approach, cross-cutting conventions, or decisions that future contributors may reasonably question.

Do not add ADRs for routine bug examples.

ADR guidance lives in [AGENTS.md](AGENTS.md) and [.agents/skills/adr-writer/SKILL.md](.agents/skills/adr-writer/SKILL.md).

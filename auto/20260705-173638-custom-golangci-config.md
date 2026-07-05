# Custom golangci-lint Config Snapshot Decision

Commit: `99026da Support custom golangci config snapshots`

## Options

1. Let users run `golangci-lint` manually inside each example.
   - Score: 2/5
   - It works, but it loses the repository snapshot workflow and makes comparison tedious.

2. Add a root `make test LINT_CONFIG=/path/to/.golangci.yaml` mode.
   - Score: 5/5
   - It keeps the existing workflow: regenerate logs, then use `git diff` to inspect behavior.

3. Add a separate target such as `make test-custom-config`.
   - Score: 4/5
   - Clear, but it duplicates the root test loop and gives contributors another command to remember.

4. Replace every example config with the supplied config.
   - Score: 1/5
   - Too broad: NilAway fixed-config examples and non-golangci examples should keep their own semantics.

## Selected

Option 2.

The root Makefile now accepts either `LINT_CONFIG=/absolute/path` or the shorter `config=/absolute/path`, converts it to an absolute path, and forwards it into example Makefiles. golangci-lint examples write the custom run output into `lint.logs` without failing early, so the final `git diff` shows exactly which expected reports changed.

Non-golangci examples still run their normal `race` or `goleak` checks. The NilAway false-positive example still verifies `.golangci.fixed.yaml`, because that fixed config is part of the example's contract rather than a user experiment.

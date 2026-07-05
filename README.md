# find_bugs

Small, self-contained Go bug examples and the tools that can catch them.

Each problem is meant to be easy to run, easy to inspect, and close enough to production code to be useful. The examples intentionally keep the bug visible without turning the code into a puzzle.

The first example uses Go's tool dependency mechanism for the custom linter build.

## Structure

```text
README.md
README.ru.md
BUGS.md
<problem_name>/<category>/
  go.mod
  Makefile
  main.go
  README.md
  README.ru.md
```

## First Example

- [nilaway/cross_package_nil](nilaway/cross_package_nil/README.md): a repository function returns `nil, nil`; the caller trusts the nil error and dereferences the nil result. NilAway can report the nil flow through a custom `golangci-lint` build before the program panics.
- [nilaway/dependency_contract_false_positive](nilaway/dependency_contract_false_positive/README.md): a dependency module exports a pointer initialized in `init`; runtime is safe, but NilAway reports the global pointer as nilable.

## Tools

- `golangci-lint`: common driver for many Go linters. NilAway currently needs to be added as a custom module plugin.
- `nilaway`: Uber's static analyzer for potential nil panics, used here through a custom `golangci-lint` build.
- `go test -race`: runtime race detector for some shared-memory concurrency bugs.
- Uber leak detector (`go.uber.org/goleak`): test-time detector for leaked goroutines.

## Repository Checks

- `make test-update`: regenerates pinned tool files and lint snapshot logs for every example.
- `make test`: runs `make test-update`, then fails if tracked files changed.

The generated `*.logs` files are committed snapshots. When a Go, golangci-lint, or NilAway version changes, `make test` shows the changed reports through `git diff`.

The full catalog of planned bugs is in [BUGS.md](BUGS.md).

Contribution notes are in [CONTRIBUTING.md](CONTRIBUTING.md).

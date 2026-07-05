# find_bugs

Small, self-contained Go bug examples and the tools that can catch them.

Each problem is meant to be easy to run, easy to inspect, and close enough to production code to be useful. The examples intentionally keep the bug visible without turning the code into a puzzle.

## Purpose

This is executable documentation for Go bug patterns and tooling behavior.

The goal is not to benchmark tools or collect clever broken snippets. The goal is to preserve small, realistic examples with the exact tool reports they produce, including true positives, false positives, and the configuration needed to handle them responsibly.

## Quick Start

Run the full repository check:

```sh
make test
```

Run one example:

```sh
cd nilaway/cross_package_nil
make run
make lint
```

Some `make lint` targets are expected to fail because they demonstrate the bug. CI uses `ci-test` targets and committed snapshot logs to verify that the expected reports stay stable.

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

## Examples

- [nilaway/cross_package_nil](nilaway/cross_package_nil/README.md): a repository function returns `nil, nil`; the caller trusts the nil error and dereferences the nil result. NilAway can report the nil flow through a custom `golangci-lint` build before the program panics.
- [nilaway/dependency_contract_false_positive](nilaway/dependency_contract_false_positive/README.md): a dependency module exports a pointer initialized in `init`; runtime is safe, but NilAway reports the global pointer as nilable.
- [goleak/channel_timeout_leak](goleak/channel_timeout_leak/README.md): a request times out while a background worker later sends to an unbuffered channel. The normal test passes, but `go.uber.org/goleak` reports the leaked goroutine.
- [goleak/context_not_cancelled](goleak/context_not_cancelled/README.md): a background cache warmer accepts a context but does not use it inside the worker. The normal test passes, but `go.uber.org/goleak` reports the goroutine left in `select`.
- [race/shared_map](race/shared_map/README.md): a metrics collector stores mutable counters in a map and reads them while another goroutine writes. `go test -race` reports the conflicting accesses.
- [race/config_pointer](race/config_pointer/README.md): a config cache refreshes a shared `*Config` while request handlers read it. `go test -race` reports the unsynchronized pointer access.
- [race/shutdown_flag](race/shutdown_flag/README.md): a worker reads a plain shutdown boolean while another goroutine writes it. `go test -race` reports the unsynchronized flag access.
- [govet/copylocks](govet/copylocks/README.md): a method copies a struct that contains `sync.Mutex`. `govet` through `golangci-lint` reports the copied lock value.
- [govet/nocopy_marker](govet/nocopy_marker/README.md): a type opts into copy detection with a private `noCopy` marker. `govet` through `golangci-lint` reports accidental value copies.
- [govet/lostcancel](govet/lostcancel/README.md): a timeout context is created but its cancel function is discarded. `govet` through `golangci-lint` reports the context leak.
- [govet/waitgroup_add_inside_goroutine](govet/waitgroup_add_inside_goroutine/README.md): `WaitGroup.Add` is called inside the goroutine it should track. `govet` through `golangci-lint` reports the lifecycle ordering bug.

## Tools

- `golangci-lint`: common driver for many Go linters. NilAway currently needs to be added as a custom module plugin.
- `nilaway`: Uber's static analyzer for potential nil panics, used here through a custom `golangci-lint` build.
- `go test -race`: runtime race detector for some shared-memory concurrency bugs.
- Uber leak detector (`go.uber.org/goleak`): test-time detector for leaked goroutines.

## Repository Checks

- `make test-update`: regenerates pinned tool files and lint snapshot logs for every example.
- `make test`: runs `make test-update`, then fails if tracked files changed.

## Why Snapshots

The generated `*.logs` files are committed snapshots. When a Go, golangci-lint, or NilAway version changes, `make test` shows the changed reports through `git diff`.

This keeps the repository honest about tool behavior. If a diagnostic changes, the diff shows what changed and forces a human to decide whether the new output is expected.

The full catalog of planned bugs is in [BUGS.md](BUGS.md).

Contribution notes are in [CONTRIBUTING.md](CONTRIBUTING.md).

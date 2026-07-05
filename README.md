# find_bugs

Small, self-contained Go bug examples and the tools that can catch them.

Each problem is meant to be easy to run, easy to inspect, and close enough to production code to be useful. The examples intentionally keep the bug visible without turning the code into a puzzle.

## Purpose

This is executable documentation for Go bug patterns and tooling behavior.

The goal is not to benchmark tools or collect clever broken snippets. The goal is to preserve small, realistic examples with the exact tool reports they produce, including true positives, false positives, and the configuration needed to handle them responsibly.

## Background Paper

The concurrency part of this catalog is guided by ["Understanding Real-World Concurrency Bugs in Go"](https://songlh.github.io/paper/go-study.pdf) by Tu, Liu, Song, and Zhang. The paper studies 171 bugs from production Go projects and uses a useful taxonomy: `blocking` vs `non-blocking` behavior, crossed with `shared memory` vs `message passing` causes.

Use that paper as the map and the examples in this repository as runnable checkpoints.

## Quick Start

Run the full repository check:

```sh
make
```

Same explicit target:

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

## Tooling Contract

Every example is shaped so a reader can turn the lesson into a local or CI guard:

- `make run` shows the program behavior.
- `make lint` runs the detector that should catch the bug.
- `make test` runs the ordinary test path for that example.
- `make ci-test` is the repository check that regenerates committed logs and asserts that the expected detector signal is still present.

Tool versions live in the example module, usually through Go tool dependencies and `go tool`. Generated helper binaries such as NilAway's custom `golangci-lint` build or the scannererr vettool are ignored by git.

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
- [concurrency/select_priority_assumption](concurrency/select_priority_assumption/README.md): a dispatcher assumes the first ready `select` case has priority. A repeated schedule check exposes that Go `select` can choose another ready case.
- [race/shared_map](race/shared_map/README.md): a metrics collector stores mutable counters in a map and reads them while another goroutine writes. `go test -race` reports the conflicting accesses.
- [race/config_pointer](race/config_pointer/README.md): a config cache refreshes a shared `*Config` while request handlers read it. `go test -race` reports the unsynchronized pointer access.
- [race/shutdown_flag](race/shutdown_flag/README.md): a worker reads a plain shutdown boolean while another goroutine writes it. `go test -race` reports the unsynchronized flag access.
- [govet/copylocks](govet/copylocks/README.md): a method copies a struct that contains `sync.Mutex`. `govet` through `golangci-lint` reports the copied lock value.
- [govet/nocopy_marker](govet/nocopy_marker/README.md): a type opts into copy detection with a private `noCopy` marker. `govet` through `golangci-lint` reports accidental value copies.
- [govet/lostcancel](govet/lostcancel/README.md): a timeout context is created but its cancel function is discarded. `govet` through `golangci-lint` reports the context leak.
- [govet/waitgroup_add_inside_goroutine](govet/waitgroup_add_inside_goroutine/README.md): `WaitGroup.Add` is called inside the goroutine it should track. `govet` through `golangci-lint` reports the lifecycle ordering bug.
- [govet/scannererr_vettool](govet/scannererr_vettool/README.md): a line importer uses `bufio.Scanner` with a small token limit and forgets `scanner.Err`. A local `go vet -vettool` wrapper runs the `scannererr` analyzer from `golang.org/x/tools`.
- [golangci/sql_rows_not_closed](golangci/sql_rows_not_closed/README.md): a repository method scans database rows and checks iteration errors, but forgets `rows.Close`. `sqlclosecheck` through `golangci-lint` reports the resource leak.
- [teamrules/ddd_repository_boundary](teamrules/ddd_repository_boundary/README.md): service code calls `*sql.DB` directly. A type-aware `ruleguard` rule keeps database calls inside repository packages.
- [teamrules/force_sqlc_query_layer](teamrules/force_sqlc_query_layer/README.md): repository code writes raw SQL through `*sql.DB` even though a generated query layer exists. A type-aware `ruleguard` rule keeps database calls inside sqlc packages.
- [teamrules/transaction_boundary](teamrules/transaction_boundary/README.md): service code starts, rolls back, and commits a transaction directly. A type-aware `ruleguard` rule keeps transaction lifecycle in transaction manager packages.
- [teamrules/no_wall_clock_in_domain](teamrules/no_wall_clock_in_domain/README.md): domain code calls `time.Now` directly. A narrow `ruleguard` rule keeps wall-clock reads in adapters or composition roots.
- [teamrules/no_panic_in_service_path](teamrules/no_panic_in_service_path/README.md): service code panics for an ordinary business failure. A narrow `ruleguard` rule keeps service paths returning errors instead.
- [synctest/context_afterfunc_negative_assertion](synctest/context_afterfunc_negative_assertion/README.md): a cancellation hook writes an audit record before the context is canceled. `testing/synctest` makes the "nothing happened yet" assertion deterministic.
- [synctest/context_timeout_without_wall_clock](synctest/context_timeout_without_wall_clock/README.md): a lease timeout is accidentally doubled. `testing/synctest` advances fake time to the deadline without wall-clock sleeps.

## Tools

- `golangci-lint`: common driver for many Go linters. NilAway currently needs to be added as a custom module plugin.
- `sqlclosecheck` and `rowserrcheck`: focused golangci-lint linters for SQL rows resource lifetime and iteration errors.
- `scannererr`: Go analysis pass from `golang.org/x/tools`, run here with a small local `go vet -vettool` binary until it is available through standard `go vet`.
- `ruleguard`: custom team rules for architecture boundaries and project-specific conventions.
- `testing/synctest`: standard-library support for deterministic tests of concurrent code, fake time, and negative assertions without wall-clock sleeps.
- `nilaway`: Uber's static analyzer for potential nil panics, used here through a custom `golangci-lint` build.
- `go test -race`: runtime race detector for some shared-memory concurrency bugs.
- Uber leak detector (`go.uber.org/goleak`): test-time detector for leaked goroutines.

## Repository Checks

- `make test-update`: regenerates pinned tool files and lint snapshot logs for every example.
- `make` or `make test`: runs `make test-update`, then fails if tracked files changed.

To try another golangci-lint config against the catalog:

```sh
make test LINT_CONFIG=/Users/me/project/.golangci.yaml
```

`config=/Users/me/project/.golangci.yaml` works as a shorter alias. In this mode golangci-lint examples write whatever the custom config reports into their `lint.logs`; examples based on `go test -race` or `goleak` keep using their own tools. The final `git diff` shows whether the custom config still catches the expected problems.

## Why Snapshots

The generated `*.logs` files are committed snapshots. When a Go, golangci-lint, or NilAway version changes, `make test` shows the changed reports through `git diff`.

This keeps the repository honest about tool behavior. If a diagnostic changes, the diff shows what changed and forces a human to decide whether the new output is expected.

The full catalog of planned bugs is in [BUGS.md](BUGS.md).

Contribution notes are in [CONTRIBUTING.md](CONTRIBUTING.md).

The working backlog for future examples and tooling is in [docs/backlog.md](docs/backlog.md).

To reproduce this repository style in another project, use the agent handoff prompt in [docs/agent-bootstrap.md](docs/agent-bootstrap.md).

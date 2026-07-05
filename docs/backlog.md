# Project Backlog

This file is the working checklist for future agents and maintainers.

`BUGS.md` is the catalog of bug patterns. This file answers a different question: what is already implemented, what is only planned, and what should be done next.

## Current Inventory

Implemented examples:

- `nilaway/cross_package_nil`
- `nilaway/dependency_contract_false_positive`
- `goleak/channel_timeout_leak`
- `goleak/context_not_cancelled`
- `race/shared_map`
- `race/config_pointer`
- `race/shutdown_flag`
- `govet/copylocks`
- `govet/nocopy_marker`
- `govet/lostcancel`
- `govet/waitgroup_add_inside_goroutine`
- `govet/scannererr_vettool`
- `golangci/sql_rows_not_closed`
- `synctest/context_afterfunc_negative_assertion`
- `teamrules/ddd_repository_boundary`
- `teamrules/no_wall_clock_in_domain`
- `teamrules/no_panic_in_service_path`
- `synctest/context_timeout_without_wall_clock`

Current shape by area:

| Area | Status |
| --- | --- |
| NilAway | Started: 2 examples, including one false-positive/config example. |
| Goleak | Started: timeout send leak and missing context propagation. |
| Race detector | Good initial set: map, pointer config, shutdown flag. |
| govet / golangci-lint | Good initial set: copylocks, noCopy, lostcancel, WaitGroup, scannererr vettool, SQL rows close. |
| synctest | Started: negative assertion and fake-time timeout examples. |
| teamrules | Started: DDD repository boundary, no wall clock in domain, and no panic in service path. This is not complete. |
| metadata/provenance | Not implemented yet. |
| GoBench import/curation | Not implemented yet. |
| GFuzz-style schedule/order examples | Not implemented yet, except backlog entries. |

## Important Gaps

### Team Rules

`teamrules` currently has only `ddd_repository_boundary`. That was the first example, not the intended endpoint.

Planned team-rule examples from `BUGS.md`:

- [x] DDD repository boundary: database calls belong in repository packages.
- [ ] Force sqlc query layer: application code should use generated query methods instead of raw SQL or direct `database/sql`.
- [ ] Transaction boundary: `BeginTx`, `Commit`, and `Rollback` belong in a unit-of-work or transaction manager package.
- [ ] No infrastructure imports in domain packages: ban `database/sql`, HTTP clients, queue clients, and loggers from domain code.
- [x] No wall clock in domain logic: ban `time.Now()` in domain packages.
- [x] No panic in service paths: ban `panic` in service packages.
- [ ] Context first argument: I/O-facing functions should accept `context.Context` as the first argument.

Good next team-rule candidates:

- `teamrules/transaction_boundary`: useful, but requires more careful package layout.
- `teamrules/force_sqlc_query_layer`: useful, but probably needs a tiny generated-like package and import/call boundary story.

### Synctest

Implemented:

- [x] `synctest/context_afterfunc_negative_assertion`
- [x] `synctest/context_timeout_without_wall_clock`

Planned:

- [ ] `synctest/unbuffered_send_after_timeout`
- [ ] document synctest limitations with mutexes and external I/O in an example README or a small docs note.

### GFuzz-Inspired Schedule Examples

Planned:

- [ ] `concurrency/select_priority_assumption`
- [ ] `concurrency/message_order_assumption`
- [ ] small stress target pattern, probably `make stress`, for examples where repeated scheduling is the detector.

Do not implement a full GFuzz clone here. Keep the repository focused on small runnable examples.

### GoBench-Inspired Metadata

GoBench is useful for its provenance and bug/fix discipline. The next repository-level improvement should be a machine-readable catalog.

Planned `catalog.yaml` fields:

- `id`
- `title`
- `category`
- `behavior`: `blocking`, `non-blocking`, `resource`, `nil`, `api-misuse`, `architecture`, etc.
- `cause`: `shared-memory`, `message-passing`, `api-contract`, `lifetime`, `architecture-boundary`, etc.
- `detector`
- `expected_signal`
- `false_positive_possible`
- `source`: paper, issue, PR, project, or manual
- `fixed_variant`: yes/no
- `difficulty`: beginner/intermediate/advanced

Follow-up tooling:

- [ ] validate `catalog.yaml`
- [ ] generate or check the implemented examples list in `BUGS.md`
- [ ] optionally produce JSON for external consumers

### Resource And API Leak Examples

Planned:

- [ ] HTTP response body leak with `bodyclose`
- [ ] ticker/timer leak
- [ ] file handle leak
- [ ] SQL rows iteration error ignored with `rowserrcheck`
- [ ] bad printf shape with `govet`

### Nil Safety

Planned:

- [ ] nil interface payload
- [ ] optional dependency not initialized
- [ ] nil map/slice write or missing initialization

## P0: Do Next

Pick one small, deterministic example and finish it end to end.

Recommended order:

1. `concurrency/select_priority_assumption`
2. `teamrules/transaction_boundary`
3. `teamrules/force_sqlc_query_layer`
4. `synctest/unbuffered_send_after_timeout`

For each example:

- create the normal self-contained directory;
- keep `make run`, `make lint`, `make test`, `make ci-test`;
- commit snapshot logs;
- update `README.md`, `README.ru.md`, and `BUGS.md`;
- run root `make`.

## P1: Add Metadata

Add `catalog.yaml` after a few more examples, before the catalog gets too large to track manually.

Acceptance criteria:

- every implemented example has one metadata entry;
- a script validates that every metadata path exists;
- `make test` or a new root target runs the validation;
- docs explain whether metadata is the source of truth or only a checked index.

## P2: Curate From GoBench

Pick 5-10 GoKer kernels that can become production-like `find_bugs` examples without copying code blindly.

For each curated example:

- preserve source/provenance;
- explain the original project and PR/issue;
- keep a small buggy version;
- include fixed variant or mitigation in README;
- use the repository Makefile/snapshot contract;
- check license/source before copying any code.

Good candidate classes:

- double lock / missing unlock;
- channel and context deadlock;
- WaitGroup misuse;
- data race on shared config;
- order violation.

## P3: Tooling Helpers

Possible helpers:

- stress runner for `go test -count=N -shuffle=on -race`;
- tiny schedule/order permutation helper for channel examples;
- custom `go/analysis` pass for one local pattern;
- ruleguard shared patterns for team rules;
- catalog JSON export.

Do not add global tool installation requirements. Everything should be reachable through `make` and Go module/tool dependencies.

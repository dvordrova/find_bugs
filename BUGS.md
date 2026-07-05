# Bug Catalog

This catalog is based on practical Go failure modes and on the taxonomy from ["Understanding Real-World Concurrency Bugs in Go"](https://songlh.github.io/paper/go-study.pdf) by Tu, Liu, Song, and Zhang. The paper studies 171 real-world Go concurrency bugs and groups them along two useful axes: behavior (`blocking` vs `non-blocking`) and cause (`shared memory` vs `message passing`).

## Nil Safety

| Bug | Example | Detection |
| --- | --- | --- |
| Nil result with nil error | A repository returns `(*T)(nil), nil`; caller checks only `err` and dereferences `T`. | NilAway through a custom `golangci-lint` build. |
| SDK init global false positive | A dependency exports a pointer initialized in `init`; runtime initialization order makes it non-nil before `main`, but the global remains nilable to the analyzer. | NilAway reports a conservative cross-module flow; handle with an explicit boundary check or a stronger SDK API that avoids exported mutable pointer globals. |
| Nil interface payload | A typed nil pointer is stored in an interface and later used as if the interface were non-nil. | NilAway may catch direct flows; tests should cover interface boundary behavior. |
| Missing map/slice initialization | Code writes to a nil map or assumes an optional slice/map exists. | `go test`; some cases by `staticcheck`; NilAway for pointer nilness rather than all container misuse. |
| Optional dependency not initialized | Struct has a pointer dependency that is only set in some constructors. | NilAway when the nil dependency flows to a dereference. |

## Concurrency: Blocking Bugs

| Bug | Example | Detection |
| --- | --- | --- |
| Unbuffered channel send after caller timeout | Worker sends to a channel after parent returns on timeout. | Usually not caught by `go test -race`; can be exposed by `testing/synctest` or checked with goleak; implemented in [synctest/unbuffered_send_after_timeout](synctest/unbuffered_send_after_timeout/README.md). |
| Receive from never-closed channel | Consumer ranges over a channel that no producer closes. | Tests with timeout; goleak for leaked goroutines. |
| Nil channel blocks forever | A select case uses a nil channel because initialization was skipped. | Tests with timeout; static review; sometimes linters catch impossible paths. |
| WaitGroup counter mismatch | `Add` and `Done` do not match on all paths. | Tests with timeout; goleak if goroutines remain blocked. |
| Mutex double lock or lock-order deadlock | One goroutine locks the same mutex twice, or two goroutines take locks in opposite order. | Tests with timeout; some static analyzers for copylocks, but lock-order bugs usually need tests or specialized tools. |
| Context cancellation not propagated | Child goroutine waits after parent request is gone. | `go.uber.org/goleak` in tests. |

## Concurrency: Non-Blocking Bugs

| Bug | Example | Detection |
| --- | --- | --- |
| Data race on shared map/state | Multiple goroutines read/write shared state without synchronization. | `go test -race`; `golangci-lint` can catch adjacent issues such as `copylocks`. |
| Loop variable captured by goroutine | Anonymous function uses changing loop variable. | Modern Go fixed the common range case, but indexed loops and pointer captures still need care; tests and linters can help. |
| WaitGroup `Add` inside goroutine | Parent may call `Wait` before the child increments the counter. | Linters can catch some WaitGroup patterns; tests with stress runs help. |
| Channel close/send race | One goroutine closes a channel while another may still send. | `go test -race` may report related races; tests should exercise shutdown ordering. |
| Select nondeterminism assumption | Code assumes priority when multiple select cases are ready. | Tests with repeated runs; design review. |

## Race Detector Examples

| Bug | Example | Detection |
| --- | --- | --- |
| Shared map race | A request metrics collector writes to a map from multiple goroutines without a mutex. | `go test -race`; snapshot the race report in `race.logs`. |
| Read/write race on cached config | One goroutine refreshes a pointer to config while request handlers read it directly. | `go test -race`; fix with `atomic.Pointer`, `sync.RWMutex`, or immutable handoff. |
| Race on shutdown flag | Worker goroutines read a plain boolean while another goroutine writes it during shutdown. | `go test -race`; fix with context cancellation, channel close, mutex, or atomic state. |

## Deterministic Concurrency Testing

| Bug | Example | Detection |
| --- | --- | --- |
| Early `context.AfterFunc` callback | Code registers cleanup or audit work for cancellation, but starts it during registration. A weak ordinary test only checks that the work eventually happened after cancel. | `testing/synctest` with `synctest.Wait` makes the negative assertion "nothing happened before cancel" deterministic. |
| Context timeout tested with wall-clock sleeps | Timeout logic is tested with real sleeps, making tests slow or flaky under load. | `testing/synctest` fake time can move to the deadline without waiting on real time; implemented in [synctest/context_timeout_without_wall_clock](synctest/context_timeout_without_wall_clock/README.md). |
| Unbuffered send after timeout | Caller returns on timeout, then the worker completes and blocks forever sending a late result to an unbuffered channel. | `testing/synctest` can advance fake time until the late send and report blocked goroutines; implemented in [synctest/unbuffered_send_after_timeout](synctest/unbuffered_send_after_timeout/README.md). |
| Select priority assumption | Code assumes a `select` prefers one ready case over another. | Repeated tests or a small schedule harness can expose the non-deterministic order assumption; implemented in [concurrency/select_priority_assumption](concurrency/select_priority_assumption/README.md). |
| Message ordering assumption | A channel protocol only works when messages arrive in one order. | A helper can run table-driven order permutations before a larger GFuzz-style tool is worth building; implemented in [concurrency/message_order_assumption](concurrency/message_order_assumption/README.md). |

## Go Vet Through golangci-lint Examples

| Bug | Example | Detection |
| --- | --- | --- |
| Copied lock value | A service struct embeds `sync.Mutex` and a method copies the struct by value before using it. | `govet` through `golangci-lint`, analyzer `copylocks`. |
| Explicit no-copy marker | A type owns background state and embeds a `noCopy` marker with `Lock`/`Unlock` pointer methods so accidental value copies are reported. | `govet` through `golangci-lint`, analyzer `copylocks`; this documents an intentional non-copyable API contract. |
| Lost context cancel | Code calls `context.WithTimeout` and forgets to call the returned cancel function on all paths. | `govet` through `golangci-lint`, analyzer `lostcancel`. |
| WaitGroup misuse | `WaitGroup.Add` is called from inside the goroutine it is supposed to track. | `govet` through `golangci-lint` when the active Go version includes the WaitGroup analyzer; otherwise use a dedicated linter or keep as a documented limitation. |
| Scanner error ignored | Code loops over `scanner.Scan()` and never checks `scanner.Err()`, so EOF and scanner failure are treated the same. | `scannererr` from `golang.org/x/tools` through a small local `go vet -vettool` wrapper; standard `go vet` does not include it yet in Go 1.26. |
| Bad printf shape | Logging or formatting uses a dynamic format with mismatched arguments. | `govet` through `golangci-lint`, analyzer `printf`. |

## Input And Stream Handling

| Bug | Example | Detection |
| --- | --- | --- |
| `bufio.Scanner` error swallowed | Scanner uses a small maximum token size, `Scan` returns `false` on a long token, and the caller forgets to check `scanner.Err()`. | `scannererr` through a local `go vet -vettool`; this also remains a good test case because standard `go vet`, `govet`, `staticcheck`, and `errcheck` do not catch it in Go 1.26. |
| SQL rows iteration error ignored | Code loops over `rows.Next()` and returns partial data without checking `rows.Err()`. | `rowserrcheck` through `golangci-lint`. |

## Resource Leaks

| Bug | Example | Detection |
| --- | --- | --- |
| Goroutine leak on early return | Function starts a worker and returns without canceling or draining it. | `go.uber.org/goleak` in tests. |
| Ticker/timer leak | `time.NewTicker` is not stopped. | Tests and review; goleak can show goroutines caused by leaked background work. |
| HTTP response body leak | Client does not close `resp.Body`. | `bodyclose` via `golangci-lint`. |
| SQL rows leak | Code reads from `*sql.Rows` and even checks `rows.Err()`, but forgets `rows.Close()`. | `sqlclosecheck` through `golangci-lint`. |
| File handle leak | Error path returns before closing a file. | `go vet`, tests, and linters such as `errcheck` for ignored close errors. |

## Team Rules And Architecture Guards

| Rule | Example | Detection |
| --- | --- | --- |
| DDD repository boundary | Service/application packages use `*sql.DB` or `*sql.Tx` directly instead of going through repository packages. | `ruleguard` with a type-aware rule that allows `database/sql` calls only under packages ending in `/repository`. |
| Force sqlc query layer | Application code writes raw SQL strings or calls `database/sql` directly instead of using generated sqlc query methods. | `ruleguard` for direct `*sql.DB`/`*sql.Tx` query calls outside generated packages; implemented in [teamrules/force_sqlc_query_layer](teamrules/force_sqlc_query_layer/README.md). |
| Transaction boundary | Code starts or commits transactions outside a unit-of-work/transaction manager package. | `ruleguard` for `BeginTx`, `Commit`, and `Rollback` calls outside allowed packages; implemented in [teamrules/transaction_boundary](teamrules/transaction_boundary/README.md). |
| No infrastructure imports in domain | Domain packages import `database/sql`, HTTP clients, loggers, or queue clients. | `depguard` for broad import boundaries; implemented in [teamrules/no_infrastructure_imports_in_domain](teamrules/no_infrastructure_imports_in_domain/README.md). |
| No wall clock in domain logic | Domain code calls `time.Now` directly instead of accepting a clock. | `ruleguard` for `time.Now()` in domain packages; implemented in [teamrules/no_wall_clock_in_domain](teamrules/no_wall_clock_in_domain/README.md). |
| No panic in service paths | Service/application packages use `panic` for ordinary error handling. | `ruleguard` for `panic($*_)` in service packages; implemented in [teamrules/no_panic_in_service_path](teamrules/no_panic_in_service_path/README.md). |
| Context first argument | I/O-facing functions accept `context.Context`, but not as the first argument. | `revive`/custom analyzer; `ruleguard` can cover common local signatures. |

## Implemented Examples

1. [nilaway/cross_package_nil](nilaway/cross_package_nil/README.md)
2. [nilaway/dependency_contract_false_positive](nilaway/dependency_contract_false_positive/README.md)
3. [goleak/channel_timeout_leak](goleak/channel_timeout_leak/README.md)
4. [goleak/context_not_cancelled](goleak/context_not_cancelled/README.md)
5. [race/shared_map](race/shared_map/README.md)
6. [race/config_pointer](race/config_pointer/README.md)
7. [race/shutdown_flag](race/shutdown_flag/README.md)
8. [govet/copylocks](govet/copylocks/README.md)
9. [govet/nocopy_marker](govet/nocopy_marker/README.md)
10. [govet/lostcancel](govet/lostcancel/README.md)
11. [govet/waitgroup_add_inside_goroutine](govet/waitgroup_add_inside_goroutine/README.md)
12. [golangci/sql_rows_not_closed](golangci/sql_rows_not_closed/README.md)
13. [govet/scannererr_vettool](govet/scannererr_vettool/README.md)
14. [teamrules/ddd_repository_boundary](teamrules/ddd_repository_boundary/README.md)
15. [synctest/context_afterfunc_negative_assertion](synctest/context_afterfunc_negative_assertion/README.md)
16. [teamrules/no_wall_clock_in_domain](teamrules/no_wall_clock_in_domain/README.md)
17. [teamrules/no_panic_in_service_path](teamrules/no_panic_in_service_path/README.md)
18. [synctest/context_timeout_without_wall_clock](synctest/context_timeout_without_wall_clock/README.md)
19. [concurrency/select_priority_assumption](concurrency/select_priority_assumption/README.md)
20. [teamrules/transaction_boundary](teamrules/transaction_boundary/README.md)
21. [teamrules/force_sqlc_query_layer](teamrules/force_sqlc_query_layer/README.md)
22. [synctest/unbuffered_send_after_timeout](synctest/unbuffered_send_after_timeout/README.md)
23. [concurrency/message_order_assumption](concurrency/message_order_assumption/README.md)
24. [teamrules/no_infrastructure_imports_in_domain](teamrules/no_infrastructure_imports_in_domain/README.md)

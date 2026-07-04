# Bug Catalog

This catalog is based on practical Go failure modes and on the taxonomy from "Understanding Real-World Concurrency Bugs in Go" by Tu, Liu, Song, and Zhang. The paper studies 171 real-world Go concurrency bugs and groups them along two useful axes: behavior (`blocking` vs `non-blocking`) and cause (`shared memory` vs `message passing`).

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
| Unbuffered channel send after caller timeout | Worker sends to a channel after parent returns on timeout. | Usually not caught by `go test -race`; can be exposed by tests with timeouts and checked with goleak. |
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

## Resource Leaks

| Bug | Example | Detection |
| --- | --- | --- |
| Goroutine leak on early return | Function starts a worker and returns without canceling or draining it. | `go.uber.org/goleak` in tests. |
| Ticker/timer leak | `time.NewTicker` is not stopped. | Tests and review; goleak can show goroutines caused by leaked background work. |
| HTTP response body leak | Client does not close `resp.Body`. | `bodyclose` via `golangci-lint`. |
| File handle leak | Error path returns before closing a file. | `go vet`, tests, and linters such as `errcheck` for ignored close errors. |

## First Implementation Order

1. `nilaway/cross_package_nil`
2. `nilaway/dependency_contract_false_positive`
3. `goleak/channel_timeout_leak`
4. `goleak/context_not_cancelled`
5. `race/shared_map`
6. `race/waitgroup_add_inside_goroutine`
